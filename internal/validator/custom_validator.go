package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	indonesianPhoneRegex = regexp.MustCompile(`^(08|\+628)[0-9]{8,11}$`)
	skuRegex             = regexp.MustCompile(`^[A-Z0-9\-]{3,30}$`)
)

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone_id", validateIndonesianPhone)
		_ = v.RegisterValidation("valid_uuid", validateUUID)
		_ = v.RegisterValidation("sku_format", validateSKU)
	}
}

func validateIndonesianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return indonesianPhoneRegex.MatchString(phone)
}

func validateUUID(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	_, err := uuid.Parse(val)
	return err == nil
}

func validateSKU(fl validator.FieldLevel) bool {
	sku := fl.Field().String()
	return skuRegex.MatchString(sku)
}
