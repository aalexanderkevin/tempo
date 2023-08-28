package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	Password *string `json:"password"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 64)),
	)
}
