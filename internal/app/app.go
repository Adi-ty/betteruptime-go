package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Adi-ty/betteruptime-go/internal/api"
	"github.com/Adi-ty/betteruptime-go/internal/config"
	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/Adi-ty/betteruptime-go/migrations"
)

type Application struct {
	Logger *log.Logger
	WebsiteHandler *api.WebsiteHandler
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	config, err := config.Load()
	if err != nil {
		return nil, err
	}

	pgDB, err := store.Open(config.DB)
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	// Stores
	websiteStore := store.NewPostgresWebsiteStore(pgDB)
	
	// Handlers
	websiteHandler := api.NewWebsiteHandler(websiteStore, logger)

	app := &Application{
		Logger: logger,
		WebsiteHandler: websiteHandler,
		DB: pgDB,
	}
	return app, nil
}