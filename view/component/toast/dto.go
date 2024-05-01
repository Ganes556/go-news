package view_toast

import (
	"encoding/json"

	dto_response "github.com/news/internal/dto/response"
)

type ParamToast struct {
	Messages string
	Mode     string
	Timer    int
}

func (p *ParamToast) getResponse() []dto_response.Response {
	var r dto_response.Response
	err := json.Unmarshal([]byte(p.Messages), &r)
	if err == nil {
		return []dto_response.Response{r}
	}

	var res []dto_response.Response
	err = json.Unmarshal([]byte(p.Messages), &res)
	if err != nil {
		res = append(res, dto_response.Response{
			Code:    500,
			Message: "Something Wrong!",
		})
	}

	return res
}
