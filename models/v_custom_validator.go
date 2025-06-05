package models

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()

	_ = validate.RegisterValidation("username_valid", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		return re.MatchString(fl.Field().String())
	})
}

func ValidateAccount(account Account) []string {
	if err := validate.Struct(account); err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' failed on rule '%s'", e.Field(), e.Tag()))
		}
		return errors
	}
	return nil
}
