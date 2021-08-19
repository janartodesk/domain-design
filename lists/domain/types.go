package domain

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EmailAddress string

func NewEmailAddress(addr string) EmailAddress {
	return EmailAddress(strings.TrimSpace(addr))
}

func (a EmailAddress) Validate() error {
	return validation.Validate(a, validation.Required)
}

type SubscriptionData map[string]interface{}
