package rest

import (
	"bytes"
	validator "github.com/go-playground/validator/v10"
	"strings"
)

type Form interface {
	IsValid() []error
}

type validationErrors struct {
	errors []error
}

func (ve *validationErrors) Error() string {
	buff := bytes.NewBufferString("")

	var e error

	for i := 0; i < len(ve.errors); i++ {
		e = ve.errors[i]
		buff.WriteString(e.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

type FormValidator struct {
	goValidator *validator.Validate
}

func CreateValidator() *FormValidator {
	return &FormValidator{goValidator: validator.New()}
}

func (v *FormValidator) validate(f Form) error {
	err := v.goValidator.Struct(f)

	if err != nil {
		return err
	}

	errs := f.IsValid()
	if errs != nil && len(errs) > 0 {
		return &validationErrors{errors: errs}
	}

	return nil
}
