package handler_categories

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	req_dto_categories "github.com/news/internal/dto/request/categories"
	dto_response "github.com/news/internal/dto/response"
	helper_handler "github.com/news/internal/handler"
	uc_categories "github.com/news/internal/usecase/categories"
	"github.com/news/pkg"
	view_admin_content_categories "github.com/news/view/admin/content/categories"
	view_admin_layout "github.com/news/view/admin/layout"
	view_navbar "github.com/news/view/component/navbar"
	view_layout "github.com/news/view/layout"
)

type HandlerCategories interface {
	ViewCategoriesAdmin(c *fiber.Ctx) error
	PostCategories(c *fiber.Ctx) error
	PutCategories(c *fiber.Ctx) error
	DelCategories(c *fiber.Ctx) error
}

type handlerCategories struct {
	uc        uc_categories.UcCategories
	validator pkg.Validator
	session   *session.Store
}

func NewHandlerCategories(uc uc_categories.UcCategories, validator pkg.Validator, session *session.Store) HandlerCategories {
	return &handlerCategories{uc, validator, session}
}

func (h *handlerCategories) ViewCategoriesAdmin(c *fiber.Ctx) error {
	fmt.Println("path 1 -> ", c.Path())
	req := new(req_dto_categories.ViewCategories)

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	ctx := c.UserContext()
	csrfToken := c.Locals("csrfToken").(string)

	categories, err := h.uc.GetAll(ctx)
	if err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	if req.Partial == "1" {
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_categories.GetCategories(categories, csrfToken),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}))
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title: "News Category",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: view_admin_content_categories.GetCategories(categories, csrfToken),
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
		C: c,
	}))
}

func (h *handlerCategories) PostCategories(c *fiber.Ctx) error {
	fmt.Println("path 2 -> ", c.Path())
	req := new(req_dto_categories.CreateCategory)
	c.BodyParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}
	ctx := c.UserContext()

	if err := h.uc.Create(ctx, *req); err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	return helper_handler.ReturnOkFlash(c, "", dto_response.Response{
		Message: "success add a new category",
		Code:    fiber.StatusOK,
	})
}

func (h *handlerCategories) PutCategories(c *fiber.Ctx) error {
	req := new(req_dto_categories.UpdateCategory)
	c.ParamsParser(req)
	c.BodyParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}
	ctx := c.UserContext()

	if err := h.uc.Update(ctx, *req); err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	return helper_handler.ReturnOkFlash(c, "", dto_response.Response{
		Message: "success edit category",
		Code:    fiber.StatusOK,
	})
}

func (h *handlerCategories) DelCategories(c *fiber.Ctx) error {
	req := new(req_dto_categories.DeleteNews)
	c.ParamsParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}
	ctx := c.UserContext()

	if err := h.uc.Delete(ctx, req.ID); err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	return helper_handler.ReturnOkFlash(c, "", dto_response.Response{
		Message: "success delete category",
		Code:    fiber.StatusOK,
	})
}
