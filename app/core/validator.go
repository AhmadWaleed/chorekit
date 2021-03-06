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
	trans     ut.Translator
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewValidator() (*Validator, error) {
	en := en.New()
	uni := ut.New(en, en)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, fmt.Errorf("could not find en translator")
	}

	v := &Validator{validator: validator.New(), trans: trans}
	v.RegisterCustomTrans()

	return v, nil
}

func (v *Validator) RegisterCustomTrans() error {
	v.validator.RegisterTranslation("required", v.trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})
	return nil
}

func (v *Validator) TranslateErrors(err error) []string {
	var errs []string
	if _, ok := err.(*validator.InvalidValidationError); ok {
		errs = append(errs, fmt.Sprintf("%s: %v", http.StatusText(http.StatusBadRequest), err))
	}

	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Translate(v.trans))
	}

	return errs
}
