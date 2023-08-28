package handler

import (
	"tempo/container"
	"tempo/controller/request"
	"tempo/controller/response"
	"tempo/helper"
	"tempo/model"
	"tempo/usecase"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	appContainer *container.Container
}

func NewUser(appContainer *container.Container) *User {
	return &User{appContainer: appContainer}
}

func (w *User) Register(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.Register")

	// Validation
	var req request.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		response.WriteFailResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		logger.WithError(err).Warning("missing required field")
		response.WriteFailResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Action
	userUseCase := usecase.NewUser(w.appContainer)
	_, err := userUseCase.Register(c, model.User{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		var e model.Error
		if !errors.As(err, &e) {
			logger.WithError(err).Warning("error register")
			response.WriteFailResponse(c, http.StatusInternalServerError, err)
		} else {
			response.WriteFailResponse(c, e.Code, e)
		}
		return
	}

	response.WriteSuccessResponse(c, nil)
}

func (w *User) Login(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.Login")

	// Validation
	var req request.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		response.WriteFailResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		logger.WithError(err).Warning("missing required field")
		response.WriteFailResponse(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Action
	userUseCase := usecase.NewUser(w.appContainer)
	res, err := userUseCase.Login(c, &model.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		var e model.Error
		if !errors.As(err, &e) {
			logger.WithError(err).Warning("error login")
			response.WriteFailResponse(c, http.StatusInternalServerError, err)
		} else {
			response.WriteFailResponse(c, e.Code, e)
		}
		return
	}

	response.WriteSuccessResponse(c, res)
}

func (w *User) UpdateUser(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.UpdateUser")

	// Validation
	var req request.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		response.WriteFailResponse(c, http.StatusBadRequest, err)
		return
	}

	// Action
	userUseCase := usecase.NewUser(w.appContainer)
	res, err := userUseCase.Update(c, req.Email, &model.User{
		Email:    req.Email,
		FullName: req.FullName,
	})
	if err != nil {
		var e model.Error
		if !errors.As(err, &e) {
			logger.WithError(err).Warning("error update user")
			response.WriteFailResponse(c, http.StatusInternalServerError, err)
		} else {
			response.WriteFailResponse(c, e.Code, e)
		}
		return
	}

	response.WriteSuccessResponse(c, res)
}
