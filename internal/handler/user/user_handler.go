package handler_user

import (
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	dto_error "github.com/news/internal/dto/error"
	req_dto_user "github.com/news/internal/dto/request/user"
	"github.com/news/internal/entity"
	helper_handler "github.com/news/internal/handler"
	uc_user "github.com/news/internal/usecase/user"
	"github.com/news/pkg"
	view_admin_content_dashboard "github.com/news/view/admin/content/dashboard"
	view_admin_content_news "github.com/news/view/admin/content/news"
	view_admin_layout "github.com/news/view/admin/layout"
	view_auth "github.com/news/view/auth"
	view_navbar "github.com/news/view/component/navbar"
	view_layout "github.com/news/view/layout"
)

type HandlerUser interface {
	GetLogin(c *fiber.Ctx) error
	PostLogin(c *fiber.Ctx) error
	PostNews(c *fiber.Ctx) error
	GetDashboard(c *fiber.Ctx) error
}

type handlerUser struct {
	uc        uc_user.UcUser
	validator pkg.Validator
	session   *session.Store
}

func NewHandlerUser(uc uc_user.UcUser, validator pkg.Validator, session *session.Store) HandlerUser {
	return &handlerUser{uc, validator, session}
}

func (h *handlerUser) GetLogin(c *fiber.Ctx) error {
	csrfToken := c.Locals("csrfToken").(string)
	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "Login",
		Contents: view_auth.Login(view_auth.ParamAuth{
			Method:    "POST",
			Action:    "/user/login",
			CsrfToken: csrfToken,
		}),
	}))
}

func (h *handlerUser) PostLogin(c *fiber.Ctx) error {
	req := new(req_dto_user.Login)
	c.BodyParser(req)
	csrfToken := c.Locals("csrfToken").(string)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title: "Login",
			Contents: view_auth.Login(view_auth.ParamAuth{
				CsrfToken:   csrfToken,
				ErrMessages: err.Errs,
				IsInvalid:   true,
				OldData:     *req,
			}),
		}))
	}
	ctx := c.UserContext()

	user, err := h.uc.Login(ctx, *req)
	if err != nil {
		var errRes dto_error.ErrResponse = dto_error.ErrResponse{
			Message: fiber.ErrInternalServerError.Message,
			Code:    fiber.ErrInternalServerError.Code,
		}
		if customErr, ok := err.(*dto_error.ErrResponse); ok {
			errRes = *customErr
		}
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title: "Login",
			Contents: view_auth.Login(view_auth.ParamAuth{
				Method:      "POST",
				Action:      "/user/login",
				CsrfToken:   csrfToken,
				ErrMessages: []dto_error.ErrResponse{errRes},
				IsInvalid:   true,
				OldData:     *req,
			}),
		}))

	}

	sess, _ := h.session.Get(c)
	sess.Set("name", user.Name)
	sess.Set("username", user.Username)
	sess.SetExpiry(time.Hour * 24 * 7)
	sess.Save()

	return c.Status(fiber.StatusSeeOther).Redirect("/user")
}

func (h *handlerUser) PostNews(c *fiber.Ctx) error {
	req := new(req_dto_user.CreateNews)
	c.BodyParser(req)
	return c.Status(fiber.StatusOK).JSON(req)
}

func (h *handlerUser) GetDashboard(c *fiber.Ctx) error {
	q := c.Queries()
	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	ctx := c.UserContext()
	var title string = "Dashboard"
	var contentComponent templ.Component
	var errRes = []dto_error.ErrResponse{}
	if q["page"] != "" {
		title = "News"
	}
	switch q["page"] {
	case "news":
		news, err := h.uc.GetNews(ctx)
		if err != nil {
			errRes = append(errRes, dto_error.ErrResponse{
				Message: fiber.ErrInternalServerError.Message,
				Code:    fiber.ErrInternalServerError.Code,
			})
		}
		contentComponent = view_admin_content_news.GetNews(news)
	case "create-news":
		contentComponent = view_admin_content_news.ModifiedNews(entity.News{}, "POST", "/news")
	case "edit-news":
		if q["id"] == "" {
			errRes = append(errRes, dto_error.ErrResponse{
				Message: "id news wajib diisi",
				Code:    400,
			})
		} else {
			id := q["id"]
			news, err := h.uc.GetNewsById(ctx, id)
			if err != nil {
				if errRe, ok := err.(*dto_error.ErrResponse); ok {
					errRes = append(errRes, *errRe)
				}
			}
			contentComponent = view_admin_content_news.ModifiedNews(news, "PUT", "/news")
		}
	default:
		contentComponent = view_admin_content_dashboard.Dashboard(username, name)
	}
	
	if q["partial"] == "1" {
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: contentComponent,
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
			ErrRes: errRes,
		}))
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: title,
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: contentComponent,
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
			ErrRes: errRes,
		}),
	}))
}
