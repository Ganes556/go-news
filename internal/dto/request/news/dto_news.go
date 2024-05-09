package req_dto_news

import "mime/multipart"

type DeleteNews struct {
	ID int `params:"id" validate:"required,gt=0"`
}
type CreateNews struct {
	Title      string                `form:"title" validate:"required"`
	CategoryID uint                  `form:"category_id" validate:"required"`
	Contents   string                `form:"contents" validate:"required"`
	Cover      *multipart.FileHeader `form:"cover" validate:"required"`
}

type ViewNews struct {
	IdEdit    uint   `query:"id" validate:"omitempty"`
	Partial   string `query:"partial" validate:"omitempty,oneof=1 0"`
	LastIndex int    `query:"last_index" validate:"omitempty,min=1"`
	Page      string `query:"page" validate:"omitempty,oneof=create update"`
	Next      uint   `query:"next" validate:"omitempty"`
	Limit     uint   `query:"next" validate:"omitempty"`
}
