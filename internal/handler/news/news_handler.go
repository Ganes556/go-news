package handler_news

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	req_dto_news "github.com/news/internal/dto/request/news"
	uc_news "github.com/news/internal/usecase/news"
	"github.com/news/pkg"
)

type NewsHandler interface {
	PostNews(c *fiber.Ctx) error
}

type newsHandler struct {
	uc        uc_news.NewsUc
	validator pkg.Validator
	session   *session.Store
}

func NewNewsHandler(uc uc_news.NewsUc, validator pkg.Validator, session *session.Store) NewsHandler {
	return &newsHandler{uc, validator, session}
}

func (h *newsHandler) PostNews(c *fiber.Ctx) error {
	req := new(req_dto_news.CreateNews)
	c.BodyParser(req)
	file, err := c.FormFile("cover")
	if err != nil {
		return c.SendStatus(500)
	}
	req.Cover = file
	ctx := c.UserContext()
	sess, _ := h.session.Get(c)
	userid := sess.Get("id")
	err = h.uc.Create(uc_news.ParamCreate{
		Ctx: ctx,
		Req: *req,
		UserID: uint(userid.(float64)),
	})

	if err != nil {
		return c.SendStatus(500)
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully add news",
	})
}
