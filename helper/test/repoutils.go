package test

import (
	"context"
	"testing"

	"tempo/helper"
	"tempo/model"
	"tempo/repository/mysqlrepo"

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
