package req_dto_user

type CreateNews struct {
	Title    string `json:"title"`
	Contents string `json:"contents"`
}
