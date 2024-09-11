package router

import (
	"crud_user/handler"
	"crud_user/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", handler.LoginHandler)
	r.POST("/register", handler.RegistrationHandler)

	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.GET("/users", handler.GetAllUserHandler)
		protected.PUT("/users", handler.UpdatePasswordHandler)
		protected.DELETE("/users", handler.DeleteAllUserHandler)
	}

	return r
}
