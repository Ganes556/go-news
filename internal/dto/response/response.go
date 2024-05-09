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

func (e *Response) Err404(field string) error {
	e.Code = 404
	e.Message = field + " not found"
	return e
}

func (e *Response) Err409(field string) error {
	e.Code = 409
	e.Message = field + " already exist"
	return e
}
