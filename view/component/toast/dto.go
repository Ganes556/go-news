package view_toast

import dto_error "github.com/news/internal/dto/error"

type ParamToastErr struct {
	ErrMessages []dto_error.ErrResponse
	Timer  int
}

type ParamToastSucc struct {
	Message string
	Timer  int
}