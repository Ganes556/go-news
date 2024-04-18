package view_admin_layout

import (
	"github.com/a-h/templ"
	dto_error "github.com/news/internal/dto/error"
)

type ParamAdminLayout struct {
	ErrRes []dto_error.ErrResponse
	SlideBar templ.Component
	Content  templ.Component
}