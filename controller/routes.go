package controller

import (
	"tempo/controller/middleware"
	_ "tempo/docs/api/rest/swag"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title User API
// @version 1.0
// @description User service REST API specification
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (h *httpServer) setupRouting() {
	router := h.engine

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "ok")
	})

	router.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API
	router.POST("/user/register", h.controllers.user.Register)
	router.POST("/user/login", h.controllers.user.Login)

	router.Use(middleware.NewHmacJwtMiddleware([]byte(h.config.JwtSecret)))
	{
		router.PUT("/user", h.controllers.user.UpdateUser)

		router.POST("/news", h.controllers.news.Add)
	}

}
