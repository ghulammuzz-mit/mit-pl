package dto

import (
	"mit/platform/internal/validator"
)

type CreateDNSRecordRequest struct {
	Name    string      `json:"name" validate:"required"`
	TTL     int         `json:"ttl" validate:"required,gt=0"`
	Type    string      `json:"type" validate:"required,dns_type"`
	Comment string      `json:"comment,omitempty"`
	Content interface{} `json:"content" validate:"required"`
	Proxied *bool       `json:"proxied,omitempty"`
}

func (r *CreateDNSRecordRequest) Validate() error {
	return validator.V.Struct(r)
}
