package pkg

import (
	"strconv"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

func ValidateRequest(input interface{}) (err error) {
	
	//if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
	//	err = fmt.Errorf("Content type not supported")
	//	break
	//}

	// validate body
	err = validate.Struct(input)
	return err
}

func ValidatorInit() {
	validate = validator.New()
	_ = validate.RegisterValidation("positive", func(fl validator.FieldLevel) bool {
		val, _ := strconv.Atoi(fl.Field().String())
		return val > 0
	})
}
