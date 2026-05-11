package validator

import "github.com/go-playground/validator/v10"

var V *validator.Validate

func Init() {
	v := validator.New()

	v.RegisterValidation("dns_type", DNSType)
	v.RegisterValidation("ipv4", IPv4)

	V = v
}
