package handler_news

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/news/helper"
	req_dto_news "github.com/news/internal/dto/request/news"
	dto_response "github.com/news/internal/dto/response"
	"github.com/news/internal/entity"
	helper_handler "github.com/news/internal/handler"
	uc_categories "github.com/news/internal/usecase/categories"
	uc_news "github.com/news/internal/usecase/news"
	"github.com/news/pkg"
	view_admin_content_news "github.com/news/view/admin/content/news"
	view_admin_layout "github.com/news/view/admin/layout"
	view_navbar "github.com/news/view/component/navbar"
	view_layout "github.com/news/view/layout"
	view_news "github.com/news/view/news"
	"github.com/sujit-baniya/flash"
)

type HandlerNews interface {
	PostNews(c *fiber.Ctx) error
	PutNews(c *fiber.Ctx) error
	ViewNewsUser(c *fiber.Ctx) error
	DelNews(c *fiber.Ctx) error
	ViewNewsAdmin(c *fiber.Ctx) error
}

type handlerNews struct {
	uc           uc_news.UcNews
	ucCategories uc_categories.UcCategories
	validator    pkg.Validator
	session      *session.Store
}

func NewNewsHandler(uc uc_news.UcNews, ucCategories uc_categories.UcCategories, validator pkg.Validator, session *session.Store) HandlerNews {
	return &handlerNews{uc, ucCategories, validator, session}
}

func (h *handlerNews) ViewNewsAdmin(c *fiber.Ctx) error {

	req := new(req_dto_news.ViewNews)
	c.QueryParser(req)
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}
	sess, _ := h.session.Get(c)
	username := sess.Get("username").(string)
	name := sess.Get("name").(string)

	ctx := c.UserContext()
	csrfToken := c.Locals("csrfToken").(string)

	news, err := h.uc.GetNews(uc_news.ParamGetNews{
		Ctx:   ctx,
		Next:  req.Next,
		Limit: req.Limit,
	})

	var component templ.Component

	if err != nil {
		helper.LogsError(err)
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	var categories []entity.Categories
	if req.Page != "" {
		categories, err = h.ucCategories.GetAll(ctx)
		if err != nil {
			helper.LogsError(err)
			return helper_handler.ReturnErrFlash(c, "", nil)
		}
	}

	switch req.Page {
	case "create":
		component = view_admin_content_news.ModifiedNews(view_admin_content_news.DtoModifiedNews{
			C:          c,
			CsrfToken:  csrfToken,
			Method:     "POST",
			Url:        "/user/news",
			Categories: categories,
		})
	case "update":
		if req.IdEdit == 0 {
			if err != nil {
				helper.LogsError(err)
				return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{
					{
						Message: "id news wajib diisi",
						Code:    400,
					},
				})
			}
		} else {
			news, err := h.uc.GetNewsById(ctx, req.IdEdit)
			if err != nil {
				helper.LogsError(err)
				if errRe, ok := err.(*dto_response.Response); ok {
					return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRe})
				} else {
					return helper_handler.ReturnErrFlash(c, "", nil)
				}
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

	if req.Partial == "1" && c.GetReqHeaders()["Hx-Request"] != nil {
		if req.Next != 0 {
			return helper_handler.Render(c, view_admin_content_news.TrNews(news, csrfToken, req.LastIndex))
		}
		return helper_handler.Render(c, view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			C:       c,
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

func (h *handlerNews) PostNews(c *fiber.Ctx) error {
	req := new(req_dto_news.CreateNews)
	c.BodyParser(req)
	file, err := c.FormFile("cover")

	if err == nil {
		req.Cover = file
	}

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	if err != nil {
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	ctx := c.UserContext()
	sess, _ := h.session.Get(c)
	userid := sess.Get("id")

	err = h.uc.Create(uc_news.ParamCreate{
		Ctx:    ctx,
		Req:    *req,
		UserID: userid.(uint),
	})

	if err != nil {
		helper.LogsError(err)
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	return flash.WithSuccess(c, fiber.Map{
		"success": true,
		"messages": helper.JSONStringify(dto_response.Response{
			Message: "successfully add article",
			Code:    200,
		}),
	}).Redirect("/user/news?page=create")
}

func (h *handlerNews) PutNews(c *fiber.Ctx) error {
	req := new(req_dto_news.UpdateNews)
	c.ParamsParser(req)
	c.BodyParser(req)
	file, err := c.FormFile("cover")

	if err == nil {
		req.Cover = file
	}

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	ctx := c.UserContext()
	sess, _ := h.session.Get(c)
	userid := sess.Get("id")

	err = h.uc.Update(uc_news.ParamUpdate{
		Ctx:    ctx,
		Req:    *req,
		UserID: userid.(uint),
	})

	if err != nil {
		helper.LogsError(err)
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}

	return helper_handler.ReturnOkFlash(c, "/user/news", dto_response.Response{
		Message: "successfully edit article",
		Code:    200,
	})
}

func (h *handlerNews) ViewNewsUser(c *fiber.Ctx) error {
	req := req_dto_news.ViewNewsUser{}
	c.QueryParser(&req)
	if err := h.validator.Validate(&req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	header := c.GetReqHeaders()
	ctx := c.UserContext()

	if req.Category != "" && header["Hx-Request"] != nil && header["Hx-Request"][0] == "true" {		
		news, err := h.uc.GetNewsByCategory(uc_news.ParamGetNewsByCategory{
			Ctx:      ctx,
			Category: req.Category,
		})

		if err != nil {
			if errRes, ok := err.(*dto_response.Response); ok {
				return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
			}
			return helper_handler.ReturnErrFlash(c, "", nil)
		}
		return helper_handler.Render(c, view_news.DataNews(news))
	}

	categories, err := h.ucCategories.GetAll(ctx)
	if err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}
	if len(categories) > 0 {
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title:    categories[0].Name,
			Contents: view_news.News(categories, categories[0].Name),
			C:        c,
		}))
	}
	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title:    "News",
		Contents: view_news.News(categories, ""),
		C:        c,
	}))
}

func (h *handlerNews) DelNews(c *fiber.Ctx) error {

	req := new(req_dto_news.DeleteNews)

	c.ParamsParser(req)

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	ctx := c.UserContext()

	if err := h.uc.Delete(uc_news.ParamDelete{
		Ctx: ctx,
		Req: *req,
	}); err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "", nil)
	}
	return helper_handler.ReturnOkFlash(c, "", dto_response.Response{
		Code:    200,
		Message: "successfully delete news",
	})
}
