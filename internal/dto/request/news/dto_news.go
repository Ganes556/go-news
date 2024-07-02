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

type UpdateNews struct {
	ID         uint                  `params:"id" validate:"required"`
	Title      string                `form:"title" validate:"omitempty"`
	CategoryID uint                  `form:"category_id" validate:"omitempty"`
	Contents   string                `form:"contents" validate:"omitempty"`
	Cover      *multipart.FileHeader `form:"cover" validate:"omitempty"`
}
type ViewNews struct {
	IdEdit    uint    `query:"id" validate:"omitempty"`
	Partial   string  `query:"partial" validate:"omitempty,oneof=1 0"`
	LastIndex int     `query:"last_index" validate:"omitempty,min=1"`
	Page      string  `query:"page" validate:"omitempty,oneof=create update"`
	Title     *string `query:"title" validate:"omitempty"`
	Category  *string `query:"category" validate:"omitempty"`
	Next      uint    `query:"next" validate:"omitempty"`
	Limit     uint    `query:"next" validate:"omitempty"`
}

type ViewNewsUser struct {
	Category   string `query:"category" validate:"omitempty"`
	MostViewed string `query:"most_viewed" validate:"omitempty,oneof=1 0"`
	Search     string `form:"search" xml:"search" json:"search" validate:"omitempty"`
	Next       uint   `query:"next" validate:"omitempty"`
	Limit      uint   `query:"next" validate:"omitempty"`
}

type ViewNewsContentUser struct {
	Title string `params:"title" validate:"required"`
}
