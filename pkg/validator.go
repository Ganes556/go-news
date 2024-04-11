package pkg

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	dto_error "github.com/news/internal/dto/error"
)

type ValidatorRes struct {
	Errs []dto_error.ErrResponse
}
type Validator interface {
	Validate(data any) *ValidatorRes
}

type xValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func translator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	return trans
}

func NewValidator() Validator {
	validate := validator.New()
	trans := translator()
	en_trans.RegisterDefaultTranslations(validate, trans)
	return &xValidator{validate, trans}
}

func convToReadAble(s string) string {
	var buf bytes.Buffer
	var prevIsUpper bool
	for _, r := range s {
		if unicode.IsNumber(r) || !unicode.IsLetter(r) {
			continue
		}
		if unicode.IsUpper(r) {
			if prevIsUpper {
				buf.WriteRune(unicode.ToLower(r))
			} else {
				if buf.Len() > 0 {
					buf.WriteRune(' ')
				}
				buf.WriteRune(unicode.ToLower(r))
				prevIsUpper = true
			}
		} else {
			buf.WriteRune(r)
			prevIsUpper = false
		}
	}

	return buf.String()
}

func (x *xValidator) Validate(data any) *ValidatorRes {

	errs := x.validator.Struct(data)

	if errs != nil {
		errosMsg := make([]dto_error.ErrResponse, len(errs.(validator.ValidationErrors)))
		for i, err := range errs.(validator.ValidationErrors) {
			errosMsg[i].Message = strings.ReplaceAll(err.Translate(x.trans), err.Field(), convToReadAble(err.Field()))
			errosMsg[i].Code = 400
		}
		return &ValidatorRes{
			Errs: errosMsg,
		}
	}

	return nil
}
