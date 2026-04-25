package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/luponetn/lexpay/internal/middleware"
)

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	authGroup := router.Group("/auth")

	authGroup.POST("/signup", handler.HandleSignUp)
	authGroup.POST("/signin", handler.HandleSignIn)
	authGroup.POST("/refresh", handler.HandleRefreshToken)
	authGroup.POST("/logout", handler.HandleLogout)

	// Protected routes
	authGroup.GET("/me", middleware.AuthMiddleware(), handler.HandleMe)
}
