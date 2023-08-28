package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type News struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (n News) Validate() error {
	return validation.ValidateStruct(
		&n,
		validation.Field(&n.Title, validation.Required),
		validation.Field(&n.Description, validation.Required),
	)
}
