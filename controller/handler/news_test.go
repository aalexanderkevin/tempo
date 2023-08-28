package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"tempo/container"
	"tempo/controller/request"
	"tempo/helper"
	"tempo/helper/test"
	"tempo/model"
	"tempo/repository/mocks"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNews_AddNews(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnErrorUnAuthorized_WhenRequestTokenIsInvalid", func(t *testing.T) {
		t.Parallel()
		// INIT
		token := "token"
		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/news", nil, map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ShouldReturnErrorBadRequest_WhenRequestPayloadIsInvalidJson", func(t *testing.T) {
		t.Parallel()
		// INIT
		token, _ := test.FakeJwtToken(t, nil)
		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/news", nil, map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ShouldReturnErrorInternalError_WhenFailedToAddNews", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		token, _ := test.FakeJwtToken(t, &fakeUser)
		reqBody := request.News{
			Title:       helper.Pointer(fake.WordsN(3)),
			Description: helper.Pointer(fake.WordsN(10)),
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		newsMock := &mocks.News{}
		newsMock.On("Add", mock.Anything, &model.News{
			UserId:      fakeUser.Id,
			Title:       reqBody.Title,
			Description: reqBody.Description,
		}).Return(nil, errors.New("error add")).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetNewsRepo(newsMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/news", &buf, map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ShouldReturnUpdatedUser_WhenSuccessUpdate", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		token, _ := test.FakeJwtToken(t, &fakeUser)
		fakeNews := test.FakeNews(t, func(news model.News) model.News {
			news.UserId = fakeUser.Id
			return news
		})
		reqBody := request.News{
			Title:       fakeNews.Title,
			Description: fakeNews.Description,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		newsMock := &mocks.News{}
		newsMock.On("Add", mock.Anything, &model.News{
			UserId:      fakeUser.Id,
			Title:       reqBody.Title,
			Description: reqBody.Description,
		}).Return(&fakeNews, nil).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetNewsRepo(newsMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/news", &buf, map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusOK, w.Code)

		resBody := model.News{}
		err = json.NewDecoder(w.Body).Decode(&resBody)
		require.NoError(t, err)

		require.Equal(t, *fakeNews.Id, *resBody.Id)
		require.Equal(t, *fakeNews.UserId, *resBody.UserId)
		require.Equal(t, *fakeNews.Title, *resBody.Title)
		require.Equal(t, *fakeNews.Description, *resBody.Description)
		require.Nil(t, resBody.CreatedAt)
		require.Nil(t, resBody.UpdatedAt)
	})

}

func TestNews_GetNews(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnErrorUnAuthorized_WhenRequestTokenIsInvalid", func(t *testing.T) {
		t.Parallel()
		// INIT
		token := "token"
		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "GET", "/news/id", nil, map[string]string{
			"Authorization": "Bearer " + token,
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ShouldReturnErrorInternalError_WhenFailedToGetNews", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		token, _ := test.FakeJwtToken(t, &fakeUser)
		id := helper.Pointer("id")

		newsMock := &mocks.News{}
		newsMock.On("GET", mock.Anything, id).Return(nil, errors.New("error get")).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetNewsRepo(newsMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "GET", "/news/"+*id, nil, map[string]string{
			"Authorization": "Bearer " + token,
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ShouldReturnExistingNews", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		token, _ := test.FakeJwtToken(t, &fakeUser)
		fakeNews := test.FakeNews(t, func(news model.News) model.News {
			news.UserId = fakeUser.Id
			return news
		})

		newsMock := &mocks.News{}
		newsMock.On("Get", mock.Anything, fakeNews.Id).Return(&fakeNews, nil).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetNewsRepo(newsMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "GET", "/news/"+*fakeNews.Id, nil, map[string]string{
			"Authorization": "Bearer " + token,
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusOK, w.Code)

		resBody := model.News{}
		err = json.NewDecoder(w.Body).Decode(&resBody)
		require.NoError(t, err)

		require.Equal(t, *fakeNews.Id, *resBody.Id)
		require.Equal(t, *fakeNews.UserId, *resBody.UserId)
		require.Equal(t, *fakeNews.Title, *resBody.Title)
		require.Equal(t, *fakeNews.Description, *resBody.Description)
		require.Nil(t, resBody.CreatedAt)
		require.Nil(t, resBody.UpdatedAt)
	})

}
