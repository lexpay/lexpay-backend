package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	authGroup := router.Group("/auth")

	authGroup.POST("/signup", handler.HandleSignUp)
	authGroup.POST("/signin", handler.HandleSignIn)
	authGroup.POST("/refresh", handler.HandleRefreshToken)
}
