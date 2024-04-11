package dto_error

type ErrResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *ErrResponse) Error() string {
	return e.Message
}

func (e *ErrResponse) ErrAuth400() error {
	e.Code = 400
	e.Message = "username atau password tidak valid"
	return e
}

func (e *ErrResponse) ErrAuth403() error {
	e.Code = 401
	e.Message = "forbidden"
	return e
}

func (e *ErrResponse) ErrAuth500() error {
	e.Code = 500
	e.Message = "internal server error"
	return e
}
