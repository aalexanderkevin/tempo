package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"tempo/container"
	"tempo/controller/request"
	"tempo/controller/response"
	"tempo/helper"
	"tempo/helper/test"
	"tempo/model"
	"tempo/repository"
	"tempo/repository/mocks"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUser_Register(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnErrorUnprocessableEntity_WhenEmailIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, nil)
		reqBody := request.User{
			Email:    nil,
			FullName: fakeUser.FullName,
			Password: fakeUser.Password,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		router := test.SetupHttpHandler(t, nil)

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/user/register", &buf, nil, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("ShouldReturnError_WhenFailedAddNewUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		reqBody := request.User{
			Email:    fakeUser.Email,
			FullName: fakeUser.FullName,
			Password: fakeUser.Password,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, model.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			require.Equal(t, *fakeUser.Email, *u.Email)
			require.Equal(t, *fakeUser.FullName, *u.FullName)

			password := helper.Pointer(helper.Hash(*u.PasswordSalt, *fakeUser.Password))
			require.Equal(t, *password, *u.Password)
			return true
		})).Return(nil, errors.New("error insert")).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/user/register", &buf, nil, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusInternalServerError, w.Code)

		resBody := response.ErrorResponse{}
		err = json.NewDecoder(w.Body).Decode(&resBody)
		require.NoError(t, err)

		require.Equal(t, "error insert", resBody.Message)
	})

	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		reqBody := request.User{
			Email:    fakeUser.Email,
			FullName: fakeUser.FullName,
			Password: fakeUser.Password,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(nil, model.NewNotFoundError()).Once()
		userMock.On("Add", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			require.Equal(t, *fakeUser.Email, *u.Email)
			require.Equal(t, *fakeUser.FullName, *u.FullName)

			password := helper.Pointer(helper.Hash(*u.PasswordSalt, *fakeUser.Password))
			require.Equal(t, *password, *u.Password)
			return true
		})).Return(&fakeUser, nil).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/user/register", &buf, nil, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUser_Login(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnError_WhenPasswordIsMissing", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		reqBody := request.User{
			Email:    fakeUser.Email,
			Password: nil,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/user/login", &buf, nil, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("ShouldReturnError_WhenPasswordIsIncorrect", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		reqBody := request.User{
			Email:    fakeUser.Email,
			Password: fakeUser.Password,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/user/login", &buf, nil, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ShouldReturnToken_WhenPasswordCorrect", func(t *testing.T) {
		t.Parallel()
		// INIT
		password := fake.CharactersN(7)
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("fakeemail@gmail.com")
			user.PasswordSalt = helper.Pointer(fake.CharactersN(7))
			user.Password = helper.Pointer(helper.Hash(*user.PasswordSalt, password))
			return user
		})

		reqBody := request.User{
			Email:    fakeUser.Email,
			Password: &password,
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "POST", "/user/login", &buf, nil, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusOK, w.Code)

		resBody := response.Login{}
		err = json.NewDecoder(w.Body).Decode(&resBody)
		require.NoError(t, err)

		require.Equal(t, *fakeUser.Id, *resBody.Id)
		require.NotNil(t, resBody.JwtToken)
	})
}

func TestUser_UpdateUser(t *testing.T) {
	t.Parallel()
	t.Run("ShouldReturnErrorUnAuthorized_WhenRequestTokenIsInvalid", func(t *testing.T) {
		t.Parallel()
		// INIT
		token := "token"
		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "PUT", "/user", nil, map[string]string{
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
		w, err := performRequest(router, "PUT", "/user", nil, map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ShouldReturnErrorInternalError_WhenFailedToUpdateUser", func(t *testing.T) {
		t.Parallel()
		// INIT
		fakeUser := test.FakeUser(t, func(user model.User) model.User {
			user.Email = helper.Pointer("email@gmail.com")
			return user
		})
		token, _ := test.FakeJwtToken(t, &fakeUser)
		reqBody := request.User{
			FullName: helper.Pointer(fake.FullName()),
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()
		userMock.On("Update", mock.Anything, *fakeUser.Id, &model.User{
			FullName: reqBody.FullName,
		}).Return(nil, errors.New("error update")).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "PUT", "/user", &buf, map[string]string{
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
		reqBody := request.User{
			FullName: helper.Pointer(fake.FullName()),
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(reqBody)
		require.NoError(t, err)

		userMock := &mocks.User{}
		userMock.On("Get", mock.Anything, repository.UserGetFilter{
			Email: fakeUser.Email,
		}).Return(&fakeUser, nil).Once()
		userMock.On("Update", mock.Anything, *fakeUser.Id, &model.User{
			FullName: reqBody.FullName,
		}).Return(&model.User{
			Id:           fakeUser.Id,
			FullName:     reqBody.FullName,
			Email:        fakeUser.Email,
			Password:     fakeUser.Password,
			PasswordSalt: fakeUser.PasswordSalt,
		}, nil).Once()

		router := test.SetupHttpHandler(t, func(appContainer *container.Container) *container.Container {
			appContainer.SetUserRepo(userMock)
			return appContainer
		})

		// CODE UNDER TEST
		w, err := performRequest(router, "PUT", "/user", &buf, map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/json",
		}, nil)
		require.NoError(t, err)
		defer printOnFailed(t)(w.Body.String())

		// EXPECTATION
		require.Equal(t, http.StatusOK, w.Code)

		resBody := model.User{}
		err = json.NewDecoder(w.Body).Decode(&resBody)
		require.NoError(t, err)

		require.Equal(t, *fakeUser.Id, *resBody.Id)
		require.Equal(t, *fakeUser.Email, *resBody.Email)
		require.Equal(t, *reqBody.FullName, *resBody.FullName)
		require.Nil(t, resBody.Password)
		require.Nil(t, resBody.PasswordSalt)
	})

}
