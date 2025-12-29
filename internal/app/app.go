package app

import (
	"log"
	"os"

	"github.com/Adi-ty/betteruptime-go/internal/api"
)

type Application struct {
	Logger *log.Logger
	WebsiteHandler *api.WebsiteHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

	websiteHandler := api.NewWebsiteHandler(logger)

	app := &Application{
		Logger: logger,
		WebsiteHandler: websiteHandler,

	}
	return app, nil
}