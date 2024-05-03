package dto_response

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Response) Error() string {
	return e.Message
}

func (e *Response) ErrAuth400() error {
	e.Code = 400
	e.Message = "invalid username or password"
	return e
}

func (e *Response) ErrAuth403() error {
	e.Code = 401
	e.Message = "forbidden"
	return e
}

func (e *Response) ErrAuth500() error {
	e.Code = 500
	e.Message = "internal server error"
	return e
}

func (e *Response) ErrNews404() error {
	e.Code = 404
	e.Message = "news not found"
	return e
}

func (e *Response) ErrCategory409() error {
	e.Code = 409
	e.Message = "category already exist"
	return e
}
