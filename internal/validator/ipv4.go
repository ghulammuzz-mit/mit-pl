package validator

import (
	"net"

	"github.com/go-playground/validator/v10"
)

func IPv4(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	ip := net.ParseIP(value)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}
