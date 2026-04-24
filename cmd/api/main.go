package main

import (
	"log"
	"log/slog"

	"github.com/luponetn/lexpay/internal/config"
	"github.com/luponetn/lexpay/internal/db"
	"github.com/luponetn/lexpay/pkg/logger"
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

	//create db conn
	connPool, err := db.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}

	defer connPool.Close()

	queries := db.New(connPool)

	//create router & start server
	router := CreateRouter()
	SetupRoutesAndServices(router, queries)
	StartServer(router, cfg.Port)

}
