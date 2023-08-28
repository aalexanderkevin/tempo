//go:build integration
// +build integration

package mysqlrepo_test

import (
	"context"
	"testing"

	"tempo/helper"
	"tempo/helper/test"
	"tempo/model"
	"tempo/repository/mysqlrepo"
	"tempo/storage"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func TestNewsRepository_Add(t *testing.T) {
	t.Run("ShouldInsertNews", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		fakeNews := test.FakeNews(t, nil)

		//-- code under test
		newsRepo := mysqlrepo.NewNewsRepository(db)
		addedNews, err := newsRepo.Add(context.TODO(), &fakeNews)

		//-- assert
		require.NoError(t, err)
		require.NotNil(t, addedNews)
		require.Equal(t, fakeNews.Title, addedNews.Title)
		require.Equal(t, fakeNews.Description, addedNews.Description)
		require.NotNil(t, addedNews.CreatedAt)
		require.NotNil(t, addedNews.UpdatedAt)
	})

	t.Run("ShouldReturnError_WhenInsertIdThatAlreadyExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		fakeNews := test.FakeNewsCreate(t, db, nil)

		//-- code under test
		newsRepo := mysqlrepo.NewNewsRepository(db)
		addedNews, err := newsRepo.Add(context.TODO(), fakeNews)

		//-- assert
		require.Error(t, err)
		require.EqualError(t, err, model.NewDuplicateError().Error())
		require.Nil(t, addedNews)
	})

}

func TestNewsRepository_Get(t *testing.T) {
	t.Run("ShouldReturnNotFoundError_WhenTheIdIsNotExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		//-- code under test
		newsRepo := mysqlrepo.NewNewsRepository(db)
		user, err := newsRepo.Get(context.TODO(), helper.Pointer("invalid-id"))
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, user)
	})

	t.Run("ShouldGet_WhenTheIdExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		fakeNews := test.FakeNewsCreate(t, db, nil)

		//-- code under test
		newsRepo := mysqlrepo.NewNewsRepository(db)
		user, err := newsRepo.Get(context.TODO(), fakeNews.Id)
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, user)
		require.Equal(t, *fakeNews.UserId, *user.UserId)
		require.Equal(t, *fakeNews.Title, *user.Title)
		require.Equal(t, *fakeNews.Description, *user.Description)
	})

}

func TestNewsRepository_Update(t *testing.T) {
	t.Run("ShouldNotFoundError_WhenIdNotExist", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)
		id := helper.Pointer("invalid-id")

		//-- code under test
		newsRepo := mysqlrepo.NewNewsRepository(db)
		user, err := newsRepo.Update(context.TODO(), id, &model.News{})
		require.Error(t, err)

		//-- assert
		require.EqualError(t, err, model.NewNotFoundError().Error())
		require.Nil(t, user)
	})

	t.Run("ShouldUpdateUser", func(t *testing.T) {
		//-- init
		db := storage.MySqlDbConn(&dbName)
		defer cleanDB(t, db)

		news := test.FakeNewsCreate(t, db, nil)
		updateNews := &model.News{
			Title: helper.Pointer(fake.Word()),
		}

		//-- code under test
		newsRepo := mysqlrepo.NewNewsRepository(db)
		res, err := newsRepo.Update(context.TODO(), news.Id, updateNews)
		require.NoError(t, err)

		//-- assert
		require.NotNil(t, res)
		require.NotEqual(t, *news.Title, *res.Title)
		require.Equal(t, *updateNews.Title, *res.Title)
		require.Equal(t, *news.UserId, *res.UserId)
		require.Equal(t, *news.Description, *res.Description)
	})

}
