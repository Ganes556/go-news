package view_admin_content_news

import (
	"github.com/gofiber/fiber/v2"
	"github.com/news/internal/entity"
)

type DtoModifiedNews struct {
	C          *fiber.Ctx
	OldNews    entity.News
	CsrfToken  string
	Method     string
	Url        string
	Categories []entity.Categories
}
