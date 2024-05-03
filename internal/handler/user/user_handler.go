package handler_user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/news/helper"
	req_dto_user "github.com/news/internal/dto/request/user"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	helper_handler "github.com/news/internal/handler"
	uc_category "github.com/news/internal/usecase/category"
	uc_news "github.com/news/internal/usecase/news"
	uc_user "github.com/news/internal/usecase/user"
	"github.com/news/pkg"
	view_admin_content_category "github.com/news/view/admin/content/category"
	view_admin_content_dashboard "github.com/news/view/admin/content/dashboard"
	view_admin_content_news "github.com/news/view/admin/content/news"
	view_admin_layout "github.com/news/view/admin/layout"
	view_auth "github.com/news/view/auth"
	view_navbar "github.com/news/view/component/navbar"
	view_layout "github.com/news/view/layout"
	"github.com/sujit-baniya/flash"
)

type HandlerUser interface {
	GetLogin(c *fiber.Ctx) error
	PostLogin(c *fiber.Ctx) error
	GetLogout(c *fiber.Ctx) error
	GetDashboard(c *fiber.Ctx) error
	GetCategory(c *fiber.Ctx) error
	GetNews(c *fiber.Ctx) error
}

type handlerUser struct {
	uc         uc_user.UcUser
	ucNews     uc_news.UcNews
	ucCategory uc_category.UcCategory
	validator  pkg.Validator
	session    *session.Store
}

func NewHandlerUser(uc uc_user.UcUser, ucNews uc_news.UcNews, ucCategory uc_category.UcCategory, validator pkg.Validator, session *session.Store) HandlerUser {
	return &handlerUser{uc, ucNews, ucCategory, validator, session}
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
		C: c,
	}))
}

func (h *handlerUser) PostLogin(c *fiber.Ctx) error {
	req := new(req_dto_user.Login)
	c.BodyParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		errs, _ := json.Marshal(err.Errs)
		mp := fiber.Map{
			"error":    true,
			"messages": string(errs),
		}
		return flash.WithError(c, mp).Redirect("/user/login")
	}

	ctx := c.UserContext()

	user, err := h.uc.Login(ctx, *req)
	if err != nil {

		var errRes dto_response.Response = dto_response.Response{
			Message: fiber.ErrInternalServerError.Message,
			Code:    fiber.ErrInternalServerError.Code,
		}

		if customErr, ok := err.(*dto_response.Response); ok {
			errRes = *customErr
		}

		msg, _ := json.Marshal(errRes)

		mp := fiber.Map{
			"error":    true,
			"messages": string(msg),
		}

		return flash.WithError(c, mp).Redirect("/user/login")
	}

	sess, _ := h.session.Get(c)
	sess.Set("id", user.ID)
	sess.Set("name", user.Name)
	sess.Set("username", user.Username)
	sess.SetExpiry(time.Hour * 24 * 7)
	sess.Save()

	return c.Redirect("/user")
}

func (h *handlerUser) GetLogout(c *fiber.Ctx) error {
	sess, _ := h.session.Get(c)
	if id := sess.Get("id"); id != nil {
		sess.Destroy()
	}
	return c.Redirect("/user/login")
}

func (h *handlerUser) GetNews(c *fiber.Ctx) error {
	q := c.Queries()
	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	ctx := c.UserContext()
	csrfToken := c.Locals("csrfToken").(string)

	news, err := h.ucNews.GetNews(uc_news.ParamGetNews{
		Ctx: ctx,
	})

	var component templ.Component

	var categories []entity.Category
	var errRes []dto_response.Response

	if err != nil {
		errRes = append(errRes, dto_response.Response{
			Message: fiber.ErrInternalServerError.Message,
			Code:    fiber.ErrInternalServerError.Code,
		})
		c.Set("HX-Refresh", "true")
	}

	if q["page"] != "" {
		categories, err = h.ucCategory.GetAll(ctx)
		if err != nil {
			helper.LogsError(err)
			errRes = append(errRes, dto_response.Response{
				Message: fiber.ErrInternalServerError.Message,
				Code:    fiber.ErrInternalServerError.Code,
			})
			c.Set("HX-Refresh", "true")
		}
	}

	switch q["page"] {
	case "create":
		component = view_admin_content_news.ModifiedNews(view_admin_content_news.DtoModifiedNews{
			C:          c,
			CsrfToken:  csrfToken,
			Method:     "POST",
			Url:        "/user/news",
			Categories: categories,
		})
	case "update":
		if q["id"] == "" {
			errRes = append(errRes, dto_response.Response{
				Message: "id news wajib diisi",
				Code:    400,
			})
		} else {
			id := q["id"]
			news, err := h.ucNews.GetNewsById(ctx, id)
			if err != nil {
				if errRe, ok := err.(*dto_response.Response); ok {
					errRes = append(errRes, *errRe)
				} else {
					errRes = append(errRes, dto_response.Response{
						Message: fiber.ErrInternalServerError.Message,
						Code:    fiber.ErrInternalServerError.Code,
					})
				}
				c.Set("HX-Refresh", "true")
			}
			component = view_admin_content_news.ModifiedNews(view_admin_content_news.DtoModifiedNews{
				C:          c,
				CsrfToken:  csrfToken,
				Method:     "PUT",
				Url:        "/user/news",
				OldNews:    news,
				Categories: categories,
			})
		}
	default:
		component = view_admin_content_news.GetNews(news, csrfToken)
	}

	if len(errRes) > 0 {
		header := c.GetReqHeaders()
		if header["Hx-Request"] != nil && header["Hx-Request"][0] == "true" {
			return flash.WithError(c, fiber.Map{
				"error":    true,
				"messages": helper.JSONStringify(errRes),
			}).SendString("")
		}
		return flash.WithError(c, fiber.Map{
			"error":    true,
			"messages": helper.JSONStringify(errRes),
		}).Redirect(c.Path())
	}

	if q["partial"] == "1" {
		fmt.Println("kena")
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: component,
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}))
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "News",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: component,
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
		C: c,
	}))
}

func (h *handlerUser) GetCategory(c *fiber.Ctx) error {
	q := c.Queries()
	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	ctx := c.UserContext()
	csrfToken := c.Locals("csrfToken").(string)

	categories, err := h.ucCategory.GetAll(ctx)
	if err != nil {
		header := c.GetReqHeaders()
		errRes := dto_response.Response{
			Message: fiber.ErrInternalServerError.Message,
			Code:    500,
		}
		if header["Hx-Request"] != nil && header["Hx-Request"][0] == "true" {
			return flash.WithError(c, fiber.Map{
				"error":    true,
				"messages": helper.JSONStringify(errRes),
			}).SendString("")
		}
		return flash.WithError(c, fiber.Map{
			"error": true,
			"messages": helper.JSONStringify(errRes),
		}).Redirect(c.Path())
	}

	if q["partial"] == "1" {
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_category.Category(categories, csrfToken),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}))
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "News Category",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_category.Category(categories, csrfToken),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
		C: c,
	}))
}

func (h *handlerUser) GetDashboard(c *fiber.Ctx) error {
	q := c.Queries()
	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	if q["partial"] == "1" {
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_dashboard.Dashboard(username, name),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}))
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "Dashboard",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_dashboard.Dashboard(username, name),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
		C: c,
	}))
}
