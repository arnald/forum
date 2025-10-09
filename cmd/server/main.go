package main

import (
	"log"
	"os"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra"
	"github.com/arnald/forum/internal/infra/logger"
	"github.com/arnald/forum/internal/infra/storage/sqlite"
)

func main() {
	// 1. Load configuration first
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// 2. Initialize DB connection
	db, err := sqlite.InitializeDB(*cfg)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	// 3. Create repository with injected DB
	logger := logger.New(os.Stdout, logger.LevelInfo)
	infraProviders := infra.NewInfraProviders(db)
	appServices := app.NewServices(infraProviders.Repositories.UserRepo, infraProviders.Repositories.CategoryRepo, infraProviders.Repositories.TopicRepo)
	infraHTTPServer := infra.NewHTTPServer(cfg, db, logger, appServices)
	infraHTTPServer.ListenAndServe()
}
