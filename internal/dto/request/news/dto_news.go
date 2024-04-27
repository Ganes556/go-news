package req_dto_news

import "mime/multipart"

type DeleteNews struct {
	ID string `params:"id"`
}
type CreateNews struct {
	Title    string                `form:"title" validate:"required"`
	Category string                `form:"category" validate:"required"`
	Contents string                `form:"contents" validate:"required"`
	Cover    *multipart.FileHeader `form:"cover" validate:"required"`
}
