package test

import (
	"context"
	"testing"
	"time"

	"tempo/config"
	"tempo/controller/middleware"
	"tempo/helper"
	"tempo/model"
	"tempo/repository/mysqlrepo"

	"github.com/dgrijalva/jwt-go"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func FakeUser(t *testing.T, cb func(user model.User) model.User) model.User {
	t.Helper()

	fakeRp := model.User{
		Id:           helper.Pointer(fake.CharactersN(7)),
		Email:        helper.Pointer(fake.EmailAddress()),
		FullName:     helper.Pointer(fake.FullName()),
		Password:     helper.Pointer(fake.CharactersN(10)),
		PasswordSalt: helper.Pointer(fake.CharactersN(7)),
	}
	if cb != nil {
		fakeRp = cb(fakeRp)
	}
	return fakeRp
}

func FakeUserCreate(t *testing.T, mysqlDB *gorm.DB, callback func(user model.User) model.User) *model.User {
	t.Helper()

	fakeData := FakeUser(t, callback)

	repo := mysqlrepo.NewUserRepository(mysqlDB)
	user, err := repo.Add(context.TODO(), &fakeData)
	require.NoError(t, err)

	return user
}

func FakeJwtToken(t *testing.T, data *model.User) (string, model.User) {
	if data == nil {
		fakeUser := FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			user.Password = nil
			user.PasswordSalt = nil
			return user
		})
		data = &fakeUser
	}

	jwtClaims := middleware.JWTData{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		User: model.User{
			Id:       data.Id,
			Email:    data.Email,
			FullName: data.FullName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessToken, err := token.SignedString([]byte(config.Instance().JwtSecret))
	if err != nil {
		t.Fatalf("Failed generating access token")
	}
	return accessToken, *data
}
