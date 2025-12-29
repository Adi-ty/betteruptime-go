package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Adi-ty/betteruptime-go/internal/api"
	"github.com/Adi-ty/betteruptime-go/internal/config"
	"github.com/Adi-ty/betteruptime-go/internal/store"
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

	logger := log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	websiteHandler := api.NewWebsiteHandler(logger)

	app := &Application{
		Logger: logger,
		WebsiteHandler: websiteHandler,
		DB: pgDB,
	}
	return app, nil
}