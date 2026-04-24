package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luponetn/lexpay/internal/auth"
	"github.com/luponetn/lexpay/internal/db"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "OK",
        })
    })

	return router
}

func SetupRoutesAndServices(router *gin.Engine, queries *db.Queries) {
	//setup services
	authService := auth.NewService(queries)

	//setup handlers
	authHandler := auth.NewHandler(authService)

	//register routes
	auth.RegisterRoutes(router, authHandler)

}

func StartServer(router *gin.Engine, port string) {
	server := &http.Server{
		Addr: ":" + port,
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	slog.Info("Server started", "port", port)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
