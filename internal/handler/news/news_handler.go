package handler_news

import (
	"fmt"

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
	view_error "github.com/news/view/error"
	view_layout "github.com/news/view/layout"
	view_news "github.com/news/view/news"
	"github.com/sujit-baniya/flash"
)

type HandlerNews interface {
	PostNews(c *fiber.Ctx) error
	PutNews(c *fiber.Ctx) error
	DelNews(c *fiber.Ctx) error
	ViewNewsHomeUser(c *fiber.Ctx) error
	ViewNewsContentUser(c *fiber.Ctx) error
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

	var news []entity.News
	var err error

	news, err = h.uc.GetNews(uc_news.ParamGetNews{
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
	categories, err = h.ucCategories.GetAll(ctx)
	if err != nil {
		helper.LogsError(err)
		return helper_handler.ReturnErrFlash(c, "", nil)
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
		component = view_admin_content_news.GetNews(news, categories, csrfToken)
	}

	if req.Partial == "1" && c.GetReqHeaders()["Hx-Request"] != nil {
		fmt.Println("request->",req.Title, req.Category)
		if req.Title != nil || req.Category != nil {	
			if *req.Title == "" && *req.Category == "" {
				return helper_handler.Render(c, view_admin_content_news.TrNews(news, csrfToken, req.LastIndex))
			}
			
			news, err = h.uc.GetNewsByFilter(uc_news.ParamGetNewsByFilter{
				Ctx:   ctx,
				Title: *req.Title,
				Category: *req.Category,
				Next:  req.Next,
				Limit: 10,
			})
			if err != nil {
				return helper_handler.ReturnErrFlash(c, "", []dto_response.Response{{
					Message: fiber.ErrInternalServerError.Message,
					Code: fiber.ErrInternalServerError.Code,
				}})
			}
			return helper_handler.Render(c, view_admin_content_news.TrNews(news, csrfToken, req.LastIndex))
		}
		
		// data infinite scroll
		if req.Next != 0 {
			return helper_handler.Render(c, view_admin_content_news.TrNews(news, csrfToken, req.LastIndex))
		}
	}

	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		C:     c,
		Title: "News",
		Contents: view_admin_layout.AdminLayout(view_admin_layout.ParamAdminLayout{
			Content: component,
			SlideBar: view_navbar.Slidebar(view_navbar.ParamNavbar{
				Username: username,
				Name:     name,
			}),
		}),
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
		return helper_handler.ReturnErrFlash(c, "/user/news?page=create", err.Errs)
	}

	if err != nil {
		return helper_handler.ReturnErrFlash(c, "/user/news?page=create", nil)
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
			return helper_handler.ReturnErrFlash(c, "/user/news?page=create", []dto_response.Response{*errRes})
		}
		return helper_handler.ReturnErrFlash(c, "/user/news?page=create", nil)
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

func (h *handlerNews) ViewNewsHomeUser(c *fiber.Ctx) error {
	req := req_dto_news.ViewNewsUser{}
	c.QueryParser(&req)
	c.BodyParser(&req)
	if err := h.validator.Validate(&req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "", err.Errs)
	}

	header := c.GetReqHeaders()
	ctx := c.UserContext()

	if header["Hx-Request"] != nil && header["Hx-Request"][0] == "true" {
		if req.Category != "" {
			news, err := h.uc.GetNewsByFilter(uc_news.ParamGetNewsByFilter{
				Ctx:      ctx,
				Category: req.Category,
				Next:     req.Next,
			})

			if err != nil {
				c.Set("HX-Retarget", "#error-get-content-news")
				c.Set("HX-Reswap", "innerHTML")
				if errRes, ok := err.(*dto_response.Response); ok {
					return c.SendString(fmt.Sprintf("%d: %s", errRes.Code, errRes.Message))
				}
				return c.SendString(fmt.Sprintf("%d: %s", fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message))
			}
			if req.Next != 0 {
				return helper_handler.Render(c, view_news.DataListNews(news, true))
			}
			return helper_handler.Render(c, view_news.DataListNews(news, false))
		}

		if req.Search != "" {
			news, err := h.uc.GetNewsByFilter(uc_news.ParamGetNewsByFilter{
				Ctx:   ctx,
				Title: req.Search,
				Next:  req.Next,
			})

			if err != nil {
				c.Set("HX-Retarget", "#error-search-news")
				c.Set("HX-Reswap", "innerHTML")
				if errRes, ok := err.(*dto_response.Response); ok {
					return c.SendString(fmt.Sprintf("%d: %s", errRes.Code, errRes.Message))
				}
				return c.SendString(fmt.Sprintf("%d: %s", fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message))
			}
			return helper_handler.Render(c, view_news.DataSearchTitle(news, req.Search))
		}

		if req.MostViewed == "1" {
			news, err := h.uc.GetNewsMostViewed(ctx)
			if err != nil {
				c.Set("HX-Retarget", "#error-most-viewed-news")
				c.Set("HX-Reswap", "innerHTML")
				if errRes, ok := err.(*dto_response.Response); ok {
					return c.SendString(fmt.Sprintf("%d: %s", errRes.Code, errRes.Message))
				}
				return c.SendString(fmt.Sprintf("%d: %s", fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message))
			}
			return helper_handler.Render(c, view_news.DataNewsMostViewed(news))
		}

		return nil
	}

	categories, err := h.ucCategories.GetAll(ctx)

	if err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
				Title:    categories[0].Name,
				Contents: view_error.Error(errRes.Message, errRes.Code),
				C:        c,
			}))
		}
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title:    categories[0].Name,
			Contents: view_error.Error(fiber.ErrInternalServerError.Message, fiber.ErrInternalServerError.Code),
			C:        c,
		}))
	}

	if req.Category != "" {
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title:    req.Category,
			Contents: view_news.NewsHome(categories, req.Category),
			C:        c,
		}))
	}
	if len(categories) > 0 {
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title:    categories[0].Name,
			Contents: view_news.NewsHome(categories, categories[0].Name),
			C:        c,
		}))
	}
	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title:    "News",
		Contents: view_news.NewsHome(categories, ""),
		C:        c,
	}))
}

func (h *handlerNews) ViewNewsContentUser(c *fiber.Ctx) error {
	req := new(req_dto_news.ViewNewsContentUser)

	c.ParamsParser(req)

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return helper_handler.ReturnErrFlash(c, "/", err.Errs)
	}

	ctx := c.UserContext()

	news, err := h.uc.GetNewsBySlug(ctx, req.Slug)
	if err != nil {
		var contentErr templ.Component = view_error.Error(fiber.ErrInternalServerError.Message, fiber.ErrInternalServerError.Code)
		if errRes, ok := err.(*dto_response.Response); ok {
			contentErr = view_error.Error(errRes.Message, errRes.Code)
		}
		return helper_handler.Render(c, view_layout.Layout(
			view_layout.ParamLayout{
				C:        c,
				Title:    "Error",
				Contents: contentErr,
			},
		))
	}

	// compute views
	if err := h.uc.AddViewingNews(uc_news.ParamAddViewingNews{
		Ctx:    ctx,
		Ip:     c.IP(),
		IdNews: news.ID,
	}); err != nil {
		var contentErr templ.Component = view_error.Error(fiber.ErrInternalServerError.Message, fiber.ErrInternalServerError.Code)
		if errRes, ok := err.(*dto_response.Response); ok {
			contentErr = view_error.Error(errRes.Message, errRes.Code)
		}
		return helper_handler.Render(c, view_layout.Layout(
			view_layout.ParamLayout{
				C:        c,
				Title:    "Error",
				Contents: contentErr,
			},
		))
	}

	categories, err := h.ucCategories.GetAll(ctx)

	if err != nil {
		if errRes, ok := err.(*dto_response.Response); ok {
			return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
				Title:    categories[0].Name,
				Contents: view_error.Error(errRes.Message, errRes.Code),
				C:        c,
			}))
		}
		return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
			Title:    categories[0].Name,
			Contents: view_error.Error(fiber.ErrInternalServerError.Message, fiber.ErrInternalServerError.Code),
			C:        c,
		}))
	}
	// if len(categories) > 0 {
	// 	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
	// 		Title:    categories[0].Name,
	// 		Contents: view_news.NewsContent(news, categories, categories[0].Name),
	// 		C:        c,
	// 	}))
	// }
	fmt.Println("kena")
	return helper_handler.Render(c, view_layout.Layout(view_layout.ParamLayout{
		Title:    news.Title,
		Contents: view_news.NewsContent(news, categories, categories[0].Name),
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
