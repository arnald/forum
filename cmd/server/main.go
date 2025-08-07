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
	userRepo := sqlite.NewRepo(db)
	logger := logger.New(os.Stdout, logger.LevelInfo)
	infraProviders := infra.NewInfraProviders(userRepo.DB)
	appServices := app.NewServices(infraProviders.UserRepository)
	infraHTTPServer := infra.NewHTTPServer(cfg, db, logger, appServices)
	infraHTTPServer.ListenAndServe()
}
