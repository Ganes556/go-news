package handler_news

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/news/helper"
	req_dto_news "github.com/news/internal/dto/request/news"
	uc_news "github.com/news/internal/usecase/news"
	"github.com/news/pkg"
	"github.com/sujit-baniya/flash"
)

type NewsHandler interface {
	PostNews(c *fiber.Ctx) error
	// DelNews(c *fiber.Ctx) error
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
	if err := h.validator.Validate(req); err != nil && len(err.Errs) > 0 {
		errs, _ := json.Marshal(err.Errs)
		mp := fiber.Map{
			"error": true,
			"code": 400,
			"messages": string(errs),
		}
		return flash.WithError(c, mp).Redirect("/user?page=create-news")
	}
	if err != nil {
		helper.LogsError(err)
		mp := fiber.Map{
			"error": true,
			"code": 500,
			"message": "Something wrong!",
		}
		return flash.WithError(c, mp).Redirect("/user?page=create-news")
	}
	req.Cover = file
	ctx := c.UserContext()
	sess, _ := h.session.Get(c)
	userid := sess.Get("id")
	err = h.uc.Create(uc_news.ParamCreate{
		Ctx:    ctx,
		Req:    *req,
		UserID: userid.(uint),
	})
	
	if err != nil {
		mp := fiber.Map{
			"error": true,
			"code": 400,
			"message": "bad request!",
		}
		return flash.WithError(c, mp).Redirect("/user?page=create-news")
	}

	mp := fiber.Map{
		"success": true,
		"message": "successfully add article",
	}
	return flash.WithSuccess(c, mp).Redirect("/user?page=create-news")
}

// func (h *newsHandler) DelNews(c *fiber.Ctx) error {

// }