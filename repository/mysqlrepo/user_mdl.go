package mysqlrepo

import (
	"tempo/model"
	"time"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

type User struct {
	Id           *string
	Email        *string
	FullName     *string
	Password     *string
	PasswordSalt *string
	CreatedAt    *time.Time
}

func (u User) FromModel(data model.User) *User {
	return &User{
		Id:           data.Id,
		Email:        data.Email,
		FullName:     data.FullName,
		Password:     data.Password,
		PasswordSalt: data.PasswordSalt,
		CreatedAt:    data.CreatedAt,
	}
}

func (u User) ToModel() *model.User {
	return &model.User{
		Id:           u.Id,
		Email:        u.Email,
		FullName:     u.FullName,
		Password:     u.Password,
		PasswordSalt: u.PasswordSalt,
		CreatedAt:    u.CreatedAt,
	}
}

func (u User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.Id == nil {
		db.Statement.SetColumn("id", ksuid.New().String())
	}

	return nil
}
