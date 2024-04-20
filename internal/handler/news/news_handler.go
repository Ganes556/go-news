package handler_news

import (
	"github.com/gofiber/fiber/v2"
	req_dto_news "github.com/news/internal/dto/request/news"
)

type NewsHandler interface {
	PostNews(c *fiber.Ctx) error
}

type newsHandler struct {
}


func NewNewsHandler() NewsHandler {
	return &newsHandler{}
}

func (h *newsHandler) PostNews(c *fiber.Ctx) error {
	req := new(req_dto_news.CreateNews)
	c.BodyParser(req)
	file, err := c.FormFile("cover")
	if err != nil {
		return c.SendStatus(500)
	}
	req.Cover = file
	return c.Status(fiber.StatusOK).JSON(req)
}