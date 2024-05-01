package handler_error

import (
	"github.com/gofiber/fiber/v2"
	helper_handler "github.com/news/internal/handler"
	view_error "github.com/news/view/error"
	view_layout "github.com/news/view/layout"
)

type ErrorHandler interface {
	NotFound(c *fiber.Ctx) error
}

type errorHandler struct {}

func NewErrorHandler() ErrorHandler {
	return &errorHandler{}
}

func (h *errorHandler) NotFound(c *fiber.Ctx) error {
	return helper_handler.Render(c, view_layout.Layout(
		view_layout.ParamLayout{
			C: c,
			Title: fiber.ErrNotFound.Message,
			Contents: view_error.Error(fiber.ErrNotFound.Message, 404),
		},
	))
}
