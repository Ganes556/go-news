package req_dto_news

import "mime/multipart"

type DeleteNews struct {
	ID int `params:"id" validate:"required,gt=0"`
}
type CreateNews struct {
	Title      string                `form:"title" validate:"required"`
	CategoriID uint                  `form:"category_id" validate:"required"`
	Contents   string                `form:"contents" validate:"required"`
	Cover      *multipart.FileHeader `form:"cover" validate:"required"`
}
