package core

import (
	"fmt"
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)

}

func NewTranslator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	return trans
}

func TransValidationErrors(err error) []string {
	var errors []string
	if _, ok := err.(*validator.InvalidValidationError); ok {
		errors = append(errors, fmt.Sprintf("%s: %v", http.StatusText(http.StatusBadRequest), err))
	}

	trans := NewTranslator()
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, err.Translate(trans))
	}

	return errors
}
