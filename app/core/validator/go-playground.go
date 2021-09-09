package validator

import (
	"fmt"
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func (v *Validator) Validate(i interface{}) []string {
	err := v.validator.Struct(i)
	if err == nil {
		return []string{}
	}

	return v.TranslateErrors(err)
}

func NewValidator() (*Validator, error) {
	en := en.New()
	uni := ut.New(en, en)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, fmt.Errorf("could not find en translator")
	}

	v := &Validator{validator: validator.New(), trans: trans}
	en_translations.RegisterDefaultTranslations(v.validator, trans)

	return v, nil
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
