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
	cfg, error := config.LoadConfig()
	if error != nil {
		log.Fatalf("Configuration error: %v", error)
	}

	// 2. Initialize DB connection
	db, error := sqlite.InitializeDB(*cfg)
	if error != nil {
		log.Fatalf("Database error: %v", error)

	}
	defer db.Close()

	// 3. Create repository with injected DB
	userRepo := sqlite.NewRepo(db)
	infraProviders := infra.NewInfraProviders(userRepo.DB)
	appServices := app.NewServices(infraProviders.UserRepository)
	infraHTTPServer := infra.NewHTTPServer(cfg, db, appServices)
	infraHTTPServer.ListenAndServe()
}
