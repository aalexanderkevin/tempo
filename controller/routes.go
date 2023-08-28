package controller

import (
	"github.com/gin-gonic/gin"
)

func (h *httpServer) setupRouting() {
	router := h.engine

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "Ok")
	})

	// API
	router.POST("/user/register", h.controllers.user.Register)
	router.POST("/user/login", h.controllers.user.Login)
	router.PUT("/user", h.controllers.user.UpdateUser)
}
