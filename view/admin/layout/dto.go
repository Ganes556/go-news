package view_admin_layout

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type ParamAdminLayout struct {
	C *fiber.Ctx
	SlideBar templ.Component
	Content  templ.Component
}