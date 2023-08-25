package model

import (
	"time"
)

type User struct {
	Id           *string    `json:"id"`
	PhoneNumber  *string    `json:"phone_number"`
	FullName     *string    `json:"full_name"`
	Password     *string    `json:"-"`
	PasswordSalt *string    `json:"-"`
	CreatedAt    *time.Time `json:"created_at"`
}

// func (u User) Validate() error {
// 	return validation.ValidateStruct(
// 		&u,
// 		validation.Field(&u.PhoneNumber, validation.Required, validation.Length(10, 13), validation.Match(regexp.MustCompile(`^\+62\d+$`))),
// 		validation.Field(&u.FullName, validation.Required, validation.Length(3, 60)),
// 		validation.Field(&u.Password, validation.Length(6, 64)),
// 	)
// }

// func (u User) ValidateLogin() error {
// 	return validation.ValidateStruct(
// 		&u,
// 		validation.Field(&u.PhoneNumber, validation.Required, validation.Length(10, 13), validation.Match(regexp.MustCompile(`^\+62\d+$`))),
// 		validation.Field(&u.Password, validation.Required, validation.Length(6, 64)),
// 	)
// }

// func IsPasswordValid(password string) bool {
// 	if len(password) < 8 {
// 		return false
// 	}

// 	hasUpperCase := false
// 	hasNumber := false
// 	hasSpecialChar := false

// 	for _, char := range password {
// 		if 'A' <= char && char <= 'Z' {
// 			hasUpperCase = true
// 		} else if '0' <= char && char <= '9' {
// 			hasNumber = true
// 		} else if !('a' <= char && char <= 'z') && !('A' <= char && char <= 'Z') && !('0' <= char && char <= '9') {
// 			hasSpecialChar = true
// 		}

// 		if hasUpperCase && hasNumber && hasSpecialChar {
// 			return true
// 		}
// 	}
// 	return false
// }
