package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/news/helper"
	conf_internal "github.com/news/internal/conf"
	handler_error "github.com/news/internal/handler/error"
	handler_news "github.com/news/internal/handler/news"
	handler_user "github.com/news/internal/handler/user"
	"github.com/news/internal/middleware"
	uc_news "github.com/news/internal/usecase/news"
	uc_user "github.com/news/internal/usecase/user"
	"github.com/news/pkg"
)

func main() {
	db_param := conf_internal.ParamNewGorm{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
	}
	DB := conf_internal.NewGorm(db_param)
	app := fiber.New()
	app.Use(recover.New())

	// helper
	encryptor := helper.NewEncryptor()
	// pkg
	validator := pkg.NewValidator()
	gcp := pkg.NewGcloud(nil, os.Getenv("GC_BUCKET"))

	// middleware
	session := session.New(session.Config{
		Storage: mysql.New(mysql.Config{
			Host: db_param.Host,
			Port: func() int {
				port, _ := strconv.Atoi(db_param.Port) // Handle error appropriately
				return port
			}(),
			Database:   db_param.Database,
			Username:   db_param.Username,
			Password:   db_param.Password,
			Table:      "fiber_store",
			Reset:      false,
			GCInterval: 10 * time.Second,
		}),
	})

	app.Static("/", "./public")

	app.Use(csrf.New(csrf.Config{
		KeyLookup:         "form:csrfToken",
		ContextKey:        "csrfToken",
		CookieSessionOnly: true,
		SingleUseToken:    true,
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Session:           session,
	}))

	timeoutMid := middleware.NewTimeoutMiddleware()
	authMid := middleware.NewAuthMiddleware(session)
	commonMid := middleware.NewCommonMid(session)

	userUc := uc_user.NewUcUser(DB, encryptor)
	newsUc := uc_news.NewNewsUc(gcp, DB)

	userGroup := app.Group("/user", timeoutMid.Timeout(nil), authMid.Authorized)
	userHandler := handler_user.NewHandlerUser(userUc, newsUc, validator, session)
	newsHandler := handler_news.NewNewsHandler(newsUc, validator, session)
	{
		userGroup.Get("/login", userHandler.GetLogin)
		userGroup.Post("/login", userHandler.PostLogin)
		userGroup.Get("/logout", userHandler.GetLogout)
		userGroup.Get("", userHandler.GetDashboard)
		userGroup.Post("/news", newsHandler.PostNews)
	}

	newsGroup := app.Group("/news", timeoutMid.Timeout(nil))
	{
		newsGroup.Get("/news", newsHandler.GetNewsUser)
	}


	errHandler := handler_error.NewErrorHandler()
	app.Use(commonMid.IsAdmin, errHandler.NotFound)

	app.Listen(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))

}
