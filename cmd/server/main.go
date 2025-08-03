package main

import (
	"log"

	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/config"
	"github.com/arnald/forum/internal/infra"
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
	infraProviders := infra.NewInfraProviders(userRepo.DB)
	appServices := app.NewServices(infraProviders.UserRepository)
	infraHTTPServer := infra.NewHTTPServer(cfg, db, appServices)
	infraHTTPServer.ListenAndServe()
}
