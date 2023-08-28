package handler

import (
	"tempo/container"
	"tempo/controller/middleware"
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

// Register New User
// @Summary 	Register New User
// @Description Register New User
// @Accept 			json
// @Param 			body 	body 		request.User 			true 	" "
// @Success 		200
// @Failure 		401 	{object}	response.ErrorResponse 	"When the auth token is missing or invalid"
// @Failure 		422 	{object}	response.ErrorResponse 	"When request validation failed"
// @Failure 		500 	{object}	response.ErrorResponse 	"When server encountered unhandled error"
// @Router /user/register [post]
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

// Login User
// @Summary 	Login User
// @Description Login User, return the token
// @Accept 			json
// @Produce 		json
// @Param 			body 	body 		request.User 			true 	" "
// @Success 		200		{object}	response.Login			"Return the user model"
// @Failure 		401 	{object}	response.ErrorResponse 	"When	the auth token is missing or invalid"
// @Failure 		422 	{object}	response.ErrorResponse 	"When request validation failed"
// @Failure 		500 	{object}	response.ErrorResponse 	"When server encountered unhandled error"
// @Router /user/login [post]
func (w *User) Login(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.Login")
	config := w.appContainer.Config()

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

	token, err := middleware.GenerateJwt(*res, config.JwtSecret)
	if err != nil {
		response.WriteFailResponse(c, http.StatusInternalServerError, err)
		return
	}

	response.WriteSuccessResponse(c, response.Login{
		Id:       res.Id,
		JwtToken: token,
	})
}

// Updater User
// @Summary 	Updater User
// @Description Updater User, return the updated user
// @Accept 			json
// @Produce 		json
// @Param 			body 	body 		request.User 			true 	" "
// @Success 		200		{object}	model.User			"Return the user model"
// @Failure 		401 	{object}	response.ErrorResponse 	"When	the auth token is missing or invalid"
// @Failure 		422 	{object}	response.ErrorResponse 	"When request validation failed"
// @Failure 		500 	{object}	response.ErrorResponse 	"When server encountered unhandled error"
// @Security 		BearerAuth
// @Router /user [put]
func (w *User) UpdateUser(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.UpdateUser")

	// auth
	user, err := middleware.GetJWTData(c)
	if err != nil {
		response.WriteFailResponse(c, http.StatusUnauthorized, err)
		return
	}

	// Validation
	var req request.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		response.WriteFailResponse(c, http.StatusBadRequest, err)
		return
	}

	// Action
	userUseCase := usecase.NewUser(w.appContainer)
	res, err := userUseCase.Update(c, user.Email, &model.User{
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
