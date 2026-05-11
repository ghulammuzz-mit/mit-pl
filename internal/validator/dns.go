package validator

import "github.com/go-playground/validator/v10"

func DNSType(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case "A", "AAAA", "CNAME", "TXT", "MX":
		return true
	default:
		return false
	}
}
