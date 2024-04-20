package req_dto_news

import "mime/multipart"

type CreateNews struct {
	Title    string                `form:"title"`
	Category string                `form:"category"`
	Contents string                `form:"contents"`
	Cover    *multipart.FileHeader `form:"cover"`
}
