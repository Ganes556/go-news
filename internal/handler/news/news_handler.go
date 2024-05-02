package handler_news

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/news/helper"
	req_dto_news "github.com/news/internal/dto/request/news"
	dto_response "github.com/news/internal/dto/response"
	helper_handler "github.com/news/internal/handler"
	uc_news "github.com/news/internal/usecase/news"
	"github.com/news/pkg"
	view_layout "github.com/news/view/layout"
	view_user_news "github.com/news/view/user"
	"github.com/sujit-baniya/flash"
)

type NewsHandler interface {
	PostNews(c *fiber.Ctx) error
	GetNewsUser(c *fiber.Ctx) error
	DelNews(c *fiber.Ctx) error
}

type newsHandler struct {
	uc        uc_news.UcNews
	validator pkg.Validator
	session   *session.Store
}

func NewNewsHandler(uc uc_news.UcNews, validator pkg.Validator, session *session.Store) NewsHandler {
	return &newsHandler{uc, validator, session}
}

func (h *newsHandler) PostNews(c *fiber.Ctx) error {
	req := new(req_dto_news.CreateNews)
	c.BodyParser(req)
	file, err := c.FormFile("cover")
	if err == nil {
		req.Cover = file
	}
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		mp := fiber.Map{
			"error":    true,
			"messages": helper.JSONStringify(err.Errs),
		}
		return flash.WithError(c, mp).Redirect("/user?page=create-news")
	}
	if err != nil {
		helper.LogsError(err)
		return flash.WithError(c, fiber.Map{
			"error": true,
			"messages": helper.JSONStringify(dto_response.Response{
				Message: fiber.ErrInternalServerError.Message,
				Code:    500,
			}),
		}).Redirect("/user?page=create-news")
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
		if errRes,ok := err.(*dto_response.Response); ok {
			return flash.WithError(c, fiber.Map{
				"error": true,
				"messages": helper.JSONStringify(errRes),
			}).Redirect("/user?page=create-news")
		}
		return flash.WithError(c, fiber.Map{
			"error": true,
			"messages": helper.JSONStringify(dto_response.Response{
				Message: fiber.ErrInternalServerError.Message,
				Code:    500,
			}),
		}).Redirect("/user?page=create-news")
	}

	return flash.WithSuccess(c, fiber.Map{
		"success": true,
		"messages": helper.JSONStringify(dto_response.Response{
			Message: "successfully add article",
			Code:    200,
		}),
	}).Redirect("/user?page=create-news")
}

func (h *newsHandler) GetNewsUser(c *fiber.Ctx) error{
	ctx := c.UserContext()
	news, err := h.uc.GetNews(uc_news.ParamGetNews{
		Ctx: ctx,
	})
	if err != nil {

	}
	
	return helper_handler.Render(c,view_layout.Layout(view_layout.ParamLayout{
		Title: "News",
		Contents: view_user_news.News(news),
		C: c,
	}))
}

func (h *newsHandler) DelNews(c *fiber.Ctx) error {

	req := new(req_dto_news.DeleteNews)

	c.ParamsParser(req)

	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		return flash.WithError(c, fiber.Map{
			"error": true,
			"messages": helper.JSONStringify(err.Errs),
		}).Redirect("/user?page=news")
	}

	ctx := c.UserContext()
	
	if err := h.uc.Delete(uc_news.ParamDelete{
		Ctx: ctx,
		Req: *req,
	}); err != nil {
		return flash.WithError(c, fiber.Map{
			"error": true,
			"messages": helper.JSONStringify(dto_response.Response{
				Code: 500,
				Message: fiber.ErrInternalServerError.Message,
			}),
		}).Redirect("/user?page=news")
	}
	return flash.WithSuccess(c, fiber.Map{
		"success": true,
		"messages": helper.JSONStringify(dto_response.Response{
			Code: 200,
			Message: fmt.Sprintf("successfully delete news with id %d", req.ID),
		}),
	}).Redirect("/user?page=news")
}
