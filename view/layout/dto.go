package view_layout

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type ParamLayout struct {
	C *fiber.Ctx
	Title    string
	Contents templ.Component
}
