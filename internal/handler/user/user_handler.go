package handler_user

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	req_dto_user "github.com/news/internal/dto/request/user"
	dto_response "github.com/news/internal/dto/response"
	helper_handler "github.com/news/internal/handler"
	uc_news "github.com/news/internal/usecase/news"
	uc_user "github.com/news/internal/usecase/user"
	"github.com/news/pkg"
	view_admin_content_dashboard "github.com/news/view/admin/content/dashboard"
	view_admin_content_profile "github.com/news/view/admin/content/profile"
	view_admin_layout "github.com/news/view/admin/layout"
	view_auth "github.com/news/view/auth"
	view_navbar "github.com/news/view/component/navbar"
	view_layout "github.com/news/view/layout"
	"github.com/sujit-baniya/flash"
)

type HandlerUser interface {
	ViewLogin(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Profile(c *fiber.Ctx) error
	PostProfile(c *fiber.Ctx) error
	ViewDashboard(c *fiber.Ctx) error
}

type handlerUser struct {
	uc        uc_user.UcUser
	ucNews    uc_news.UcNews
	validator pkg.Validator
	session   *session.Store
}

func NewHandlerUser(uc uc_user.UcUser, ucNews uc_news.UcNews, validator pkg.Validator, session *session.Store) HandlerUser {
	return &handlerUser{uc, ucNews, validator, session}
}

func (h *handlerUser) ViewLogin(c *fiber.Ctx) error {
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

func (h *handlerUser) Login(c *fiber.Ctx) error {
	req := new(req_dto_user.Login)
	c.BodyParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "/user/login", err.Errs)
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

func (h *handlerUser) Profile(c *fiber.Ctx) error {
	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)
	csrfToken := c.Locals("csrfToken").(string)

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "Profile",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_profile.Profile(name, username, csrfToken),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
		C: c,
	}))
}

func (h *handlerUser) PostProfile(c *fiber.Ctx) error {
	req := new(req_dto_user.EditProfile)
	c.BodyParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}
	sess, _ := h.session.Get(c)
	id := sess.Get("id").(uint)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)
	ctx := c.UserContext()
	err := h.uc.EditProfile(ctx, id, *req)
	if err != nil {
		if errRe, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRe})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	if req.Name != name {
		sess.Set("name", req.Name)
	}
	if req.Username != username {
		sess.Set("username", req.Username)
	}
	if req.Name != name || req.Username != username {
		sess.Save()
	}

	return helper_handler.ReturnOkFlash(c, "/user/profile", dto_response.Response{
		Message: "successfully edit profile",
		Code: 200,
	})
}

func (h *handlerUser) Logout(c *fiber.Ctx) error {
	sess, _ := h.session.Get(c)
	if id := sess.Get("id"); id != nil {
		sess.Destroy()
	}
	return c.Redirect("/user/login")
}

func (h *handlerUser) ViewDashboard(c *fiber.Ctx) error {
	req := new(req_dto_user.ViewDashboard)
	c.QueryParser(req)

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	ctx := c.UserContext()

	totalPost, totalViews := h.ucNews.GetTotalPostAndViews(ctx)

	if req.Partial == "1" {
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_dashboard.Dashboard(username, name, totalPost, totalViews),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}))
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "Dashboard",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_dashboard.Dashboard(username, name, totalPost, totalViews),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
		C: c,
	}))
}
