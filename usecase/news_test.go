package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"tempo/container"
	"tempo/helper"
	"tempo/helper/test"
	"tempo/model"
	"tempo/repository/mocks"
	"tempo/usecase"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNews_Add(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnError_WhenTitleIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		appContainer := container.Container{}

		fakeNews := test.FakeNews(t, func(news model.News) model.News {
			news.Title = nil
			return news
		})

		// CODE UNDER TEST
		uc := usecase.NewNews(&appContainer)
		res, err := uc.Add(context.Background(), &fakeNews)
		require.Error(t, err)
		require.True(t, model.IsParameterError(err))
		require.Nil(t, res)
	})

	t.Run("ShouldReturnError_WhenErrorAddUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeNews := test.FakeNews(t, nil)

		newsMock := &mocks.News{}
		newsMock.On("Add", mock.Anything, &fakeNews).Return(nil, errors.New("error insert")).Once()

		appContainer := container.Container{}
		appContainer.SetNewsRepo(newsMock)

		// CODE UNDER TEST
		uc := usecase.NewNews(&appContainer)
		res, err := uc.Add(context.Background(), &fakeNews)
		require.Error(t, err)
		require.EqualError(t, err, "error insert")
		require.Nil(t, res)

		newsMock.AssertExpectations(t)
	})

	t.Run("ShouldNotReturnError_WhenSuccessInsertNews", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeNews := test.FakeNews(t, nil)

		newsMock := &mocks.News{}
		newsMock.On("Add", mock.Anything, &fakeNews).Return(&model.News{
			Id:          helper.Pointer(fake.CharactersN(6)),
			Title:       fakeNews.Title,
			Description: fakeNews.Description,
			CreatedAt:   helper.Pointer(time.Now()),
			UpdatedAt:   helper.Pointer(time.Now()),
		}, nil).Once()

		appContainer := container.Container{}
		appContainer.SetNewsRepo(newsMock)

		// CODE UNDER TEST
		uc := usecase.NewNews(&appContainer)
		res, err := uc.Add(context.Background(), &fakeNews)
		require.NoError(t, err)
		require.NotNil(t, res.Id)
		require.Equal(t, *fakeNews.Title, *res.Title)
		require.Equal(t, *fakeNews.Description, *res.Description)
		require.NotNil(t, res.CreatedAt)
		require.NotNil(t, res.UpdatedAt)

		newsMock.AssertExpectations(t)
	})
}

