package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type News struct {
	Id          *string    `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (n News) Validate() error {
	return validation.ValidateStruct(
		&n,
		validation.Field(&n.Title, validation.Required),
		validation.Field(&n.Description, validation.Required),
	)
}
