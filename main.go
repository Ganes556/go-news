package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jasonlvhit/gocron"
	_ "github.com/joho/godotenv/autoload"
	"github.com/news/helper"
	conf_internal "github.com/news/internal/conf"
	"github.com/news/internal/cron_func"
	handler_categories "github.com/news/internal/handler/categories"
	handler_error "github.com/news/internal/handler/error"
	handler_news "github.com/news/internal/handler/news"
	handler_user "github.com/news/internal/handler/user"
	"github.com/news/internal/middleware"
	uc_categories "github.com/news/internal/usecase/categories"
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
	DB, fStorage := conf_internal.NewGorm(db_param)
	app := fiber.New()
	app.Use(recover.New())

	// helper
	encryptor := helper.NewEncryptor()
	// pkg
	validator := pkg.NewValidator()
	gcp := pkg.NewGcloud(nil, os.Getenv("GC_BUCKET"))

	// cron
	cronFunc := cron_func.NewCronFunc(DB)
	s := gocron.NewScheduler()
	err := s.Every(1).Day().Do(cronFunc.DBResetIpread)
	if err != nil {
		panic(err)
	}
	go func() {
		<-s.Start()
	}()

	// middleware
	session := session.New(session.Config{
		CookieHTTPOnly: true,
		Storage:        fStorage,
	})

	app.Static("/", "./static")

	app.Use(csrf.New(csrf.Config{
		CookieHTTPOnly:    true,
		ContextKey:        "csrfToken",
		CookieSessionOnly: true,
		SingleUseToken:    true,
		Expiration:        1 * time.Hour,
		Extractor: func(c *fiber.Ctx) (string, error) {
			tokenFromQuery := c.Query("csrfToken")
			if tokenFromQuery != "" {
				return tokenFromQuery, nil
			}
			// If not found in the query parameters, attempt to extract from form data
			tokenFromForm := c.FormValue("csrfToken")

			if tokenFromForm != "" {
				return tokenFromForm, nil
			}

			// If not found in either query parameters or form data, return an error
			return "", errors.New("csrf token not found")
		},
		KeyGenerator: utils.UUIDv4,
		Session:      session,
	}))

	if os.Getenv("SSL") == "1" {
		app.Use(redirectToHTTPS)
	}

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/news",
		},
		StatusCode: 301,
	}))

	timeoutMid := middleware.NewTimeoutMiddleware()
	authMid := middleware.NewAuthMiddleware(session)
	commonMid := middleware.NewCommonMid(session)

	// usecase
	userUc := uc_user.NewUcUser(DB, encryptor)
	newsUc := uc_news.NewNewsUc(gcp, DB)
	categoriesUc := uc_categories.NewCategoriesUc(DB, gcp)
	// handler
	userHandler := handler_user.NewHandlerUser(userUc, newsUc, validator, session)
	newsHandler := handler_news.NewNewsHandler(newsUc, categoriesUc, validator, session)
	categoriesHandler := handler_categories.NewHandlerCategories(categoriesUc, validator, session)

	userGroup := app.Group("/user", timeoutMid.Timeout(nil), authMid.Authorized)
	{

		userGroup.Get("/profile", userHandler.Profile)
		userGroup.Post("/profile", userHandler.PostProfile)
		userGroup.Get("/login", userHandler.ViewLogin)
		userGroup.Post("/login", userHandler.Login)
		userGroup.Get("/logout", userHandler.Logout)
		userGroup.Get("", userHandler.ViewDashboard)
		userGroup.Get("/news", newsHandler.ViewNewsAdmin)
		userGroup.Post("/news", newsHandler.PostNews)
		userGroup.Put("/news/:id", newsHandler.PutNews)
		userGroup.Post("/news/categories", categoriesHandler.PostCategories)
		userGroup.Get("/news/categories", categoriesHandler.ViewCategoriesAdmin)
		userGroup.Put("/news/categories/:id", categoriesHandler.PutCategories)
		userGroup.Delete("/news/categories/:id", categoriesHandler.DelCategories)
		userGroup.Delete("/news/:id", newsHandler.DelNews)
	}

	newsGroup := app.Group("/news", timeoutMid.Timeout(nil))
	{
		newsGroup.Get("/", newsHandler.ViewNewsHomeUser)
		newsGroup.Get("/:slug", newsHandler.ViewNewsContentUser)
	}

	errHandler := handler_error.NewErrorHandler()
	app.Use(commonMid.IsAdmin, errHandler.NotFound)

	if os.Getenv("SSL") == "1" {
		cer, err := tls.LoadX509KeyPair(os.Getenv("SSL_CERT"), os.Getenv("SSL_KEY"))
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{
			Certificates: []tls.Certificate{cer},
		}

		go func() {
			if err := app.Listen(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))); err != nil {
				panic(err)
			}
		}()

		ln, err := tls.Listen("tcp", ":443", config)
		if err != nil {
			panic(err)
		}

		log.Fatal(app.Listener(ln))
	} else {
		app.Listen(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
	}
}

func redirectToHTTPS(c *fiber.Ctx) error {
	if c.Protocol() == "http" {
		// Redirect to the HTTPS version of the requested URL
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
	}
	return c.Next()
}
