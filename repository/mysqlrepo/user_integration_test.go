//go:build integration
// +build integration

package mysqlrepo_test

import (
	"context"
	"testing"

	"tempo/helper"
	"tempo/helper/test"
	"tempo/model"
	"tempo/repository"
	"tempo/repository/mysqlrepo"
	"tempo/storage"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Add(t *testing.T) {
	t.Run("ShouldInsertUser", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUser(t, nil)

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		addedUser, err := userRepo.Add(context.TODO(), &fakeUser)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, addedUser)
		require.Equal(t, fakeUser.Email, addedUser.Email)
		require.Equal(t, fakeUser.FullName, addedUser.FullName)
		require.Equal(t, fakeUser.Password, addedUser.Password)
		require.Equal(t, fakeUser.PasswordSalt, addedUser.PasswordSalt)
		require.NotNil(t, addedUser.CreatedAt)
	})

	t.Run("ShouldReturnError_WhenInsertIdThatAlreadyExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUserCreate(t, db, nil)
		fakeUser.Id = nil

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		addedUser, err := userRepo.Add(context.TODO(), fakeUser)

		//-- assert
		require.Error(t, err)
		require.EqualError(t, err, model.NewDuplicateError().Error())
		require.Nil(t, addedUser)
	})

}

func TestUserRepository_Get(t *testing.T) {
	t.Run("ShouldReturnNotFoundError_WhenTheIdIsNotExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		user, err := userRepo.Get(context.TODO(), repository.UserGetFilter{
			Id: helper.Pointer("invalid-id"),
		})
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, user)
	})

	t.Run("ShouldGet_WhenTheIdExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		fakeUser := test.FakeUserCreate(t, db, nil)

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		user, err := userRepo.Get(context.TODO(), repository.UserGetFilter{
			Id: fakeUser.Id,
		})
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, user)
		require.Equal(t, *fakeUser.Email, *user.Email)
		require.Equal(t, *fakeUser.FullName, *user.FullName)
		require.Equal(t, *fakeUser.Password, *user.Password)
		require.Equal(t, *fakeUser.PasswordSalt, *user.PasswordSalt)
		require.NotNil(t, user.CreatedAt)
	})

}

func TestUserRepository_Update(t *testing.T) {
	t.Run("ShouldNotFoundError_WhenIdNotExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)
		invalidId := "invalid-id"

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		user, err := userRepo.Update(context.TODO(), invalidId, &model.User{
			Email: helper.Pointer(fake.EmailAddress()),
		})
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, user)
	})

	t.Run("ShouldUpdateUser", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		user := test.FakeUserCreate(t, db, nil)
		updateUser := &model.User{
			Email:    helper.Pointer(fake.EmailAddress()),
			FullName: helper.Pointer(fake.FullName()),
		}

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		res, err := userRepo.Update(context.TODO(), *user.Id, updateUser)
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, res)
		require.NotEqual(t, *user.Email, *res.Email)
		require.NotEqual(t, *user.FullName, *res.FullName)
		require.Equal(t, *updateUser.Email, *res.Email)
		require.Equal(t, *updateUser.FullName, *res.FullName)
		require.Equal(t, *user.Password, *res.Password)
		require.Equal(t, *user.PasswordSalt, *res.PasswordSalt)
	})

	t.Run("ShouldErrorDuplicate_WhenUpdateEmailThatAlreadyExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		userExist := test.FakeUserCreate(t, db, nil)
		user := test.FakeUserCreate(t, db, nil)
		updateUser := &model.User{
			Email: userExist.Email,
		}

		//-- code under test
		userRepo := mysqlrepo.NewUserRepository(db)
		res, err := userRepo.Update(context.TODO(), *user.Id, updateUser)

		//-- assert
		require.Error(t, err)
		require.EqualError(t, err, model.NewDuplicateError().Error())
		require.Nil(t, res)
	})

}
