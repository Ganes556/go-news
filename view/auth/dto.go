package view_auth

import (
	"github.com/gofiber/fiber/v2"
	req_dto_user "github.com/news/internal/dto/request/user"
)

type ParamAuth struct {
	C         *fiber.Ctx
	Method    string
	Action    string
	CsrfToken string
	IsInvalid bool
	OldData   req_dto_user.Login
}
