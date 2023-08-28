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

type News struct {
	appContainer *container.Container
}

func NewNews(appContainer *container.Container) *News {
	return &News{appContainer: appContainer}
}

// Add New News
// @Summary 	Add New News
// @Description Add New News
// @Accept 			json
// @Produce 		json
// @Param 			body 	body 		request.News 			true 	" "
// @Success 		200		{object}	model.News				"Return the news model"
// @Failure 		401 	{object}	response.ErrorResponse 	"When	the auth token is missing or invalid"
// @Failure 		422 	{object}	response.ErrorResponse 	"When request validation failed"
// @Failure 		500 	{object}	response.ErrorResponse 	"When server encountered unhandled error"
// @Security 		BearerAuth
// @Router /news [post]
func (w *News) Add(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.Add")

	// auth
	user, err := middleware.GetJWTData(c)
	if err != nil {
		response.WriteFailResponse(c, http.StatusUnauthorized, err)
		return
	}

	// Validation
	var req request.News
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		response.WriteFailResponse(c, http.StatusBadRequest, err)
		return
	}

	// Action
	newsUseCase := usecase.NewNews(w.appContainer)
	res, err := newsUseCase.Add(c, &model.News{
		UserId:      user.Id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		var e model.Error
		if !errors.As(err, &e) {
			logger.WithError(err).Warning("error add news")
			response.WriteFailResponse(c, http.StatusInternalServerError, err)
		} else {
			response.WriteFailResponse(c, e.Code, e)
		}
		return
	}

	response.WriteSuccessResponse(c, res)
}

// Get News
// @Summary 	Get News
// @Description Get News
// @Produce 		json
// @Param id path string true "news id"
// @Success 		200		{object}	model.News				"Return the news model"
// @Failure 		401 	{object}	response.ErrorResponse 	"When	the auth token is missing or invalid"
// @Failure 		422 	{object}	response.ErrorResponse 	"When request validation failed"
// @Failure 		500 	{object}	response.ErrorResponse 	"When server encountered unhandled error"
// @Security 		BearerAuth
// @Router /news/:id [get]
func (w *News) Get(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.Add")

	// auth
	_, err := middleware.GetJWTData(c)
	if err != nil {
		response.WriteFailResponse(c, http.StatusUnauthorized, err)
		return
	}

	// Validation
	id := c.Param("id")
	if id == "" {
		response.WriteFailResponse(c, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	// Action
	newsUseCase := usecase.NewNews(w.appContainer)
	res, err := newsUseCase.Get(c, &id)
	if err != nil {
		var e model.Error
		if !errors.As(err, &e) {
			logger.WithError(err).Warning("error get news")
			response.WriteFailResponse(c, http.StatusInternalServerError, err)
		} else {
			response.WriteFailResponse(c, e.Code, e)
		}
		return
	}

	response.WriteSuccessResponse(c, res)
}

// Update News
// @Summary 	Update News
// @Description Update News
// @Accept 			json
// @Produce 		json
// @Param id path string true "news id"
// @Param 			body 	body 		request.News 			true 	" "
// @Success 		200		{object}	model.News				"Return the news model"
// @Failure 		401 	{object}	response.ErrorResponse 	"When	the auth token is missing or invalid"
// @Failure 		422 	{object}	response.ErrorResponse 	"When request validation failed"
// @Failure 		500 	{object}	response.ErrorResponse 	"When server encountered unhandled error"
// @Security 		BearerAuth
// @Router /news/:id [put]
func (w *News) Update(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Controller.Handler.Update")

	// auth
	_, err := middleware.GetJWTData(c)
	if err != nil {
		response.WriteFailResponse(c, http.StatusUnauthorized, err)
		return
	}

	// Validation
	id := c.Param("id")

	var req request.News
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Warning("bad request error")
		response.WriteFailResponse(c, http.StatusBadRequest, err)
		return
	}

	// Action
	newsUseCase := usecase.NewNews(w.appContainer)
	res, err := newsUseCase.Update(c, &id, &model.News{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		var e model.Error
		if !errors.As(err, &e) {
			logger.WithError(err).Warning("error update news")
			response.WriteFailResponse(c, http.StatusInternalServerError, err)
		} else {
			response.WriteFailResponse(c, e.Code, e)
		}
		return
	}

	response.WriteSuccessResponse(c, res)
}
