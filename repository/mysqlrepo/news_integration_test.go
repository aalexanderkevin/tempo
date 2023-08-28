//go:build integration
// +build integration

package mysqlrepo_test

import (
	"context"
	"testing"

	"tempo/helper/test"
	"tempo/model"
	"tempo/repository/mysqlrepo"
	"tempo/storage"

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

