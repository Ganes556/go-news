package view_auth

import (
	dto_error "github.com/news/internal/dto/error"
	req_dto_user "github.com/news/internal/dto/request/user"
)

type ParamAuth struct {
	Method      string
	Action      string
	CsrfToken   string
	ErrMessages []dto_error.ErrResponse
	IsInvalid   bool
	OldData     req_dto_user.Login
}
