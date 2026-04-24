package main

import (
	"log"
	"log/slog"

	"github.com/luponetn/lexpay/internals/config"
	"github.com/luponetn/lexpay/internals/logger"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup global logger
	logger.SetupLogger(cfg.Env)

	slog.Info("Service started",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.Port),
	)

	//create router & start server
	router := CreateRouter()
	SetupRoutes(router)
	StartServer(router, cfg.Port)

}
