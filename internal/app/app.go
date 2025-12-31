package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/Adi-ty/betteruptime-go/internal/api"
	"github.com/Adi-ty/betteruptime-go/internal/config"
	"github.com/Adi-ty/betteruptime-go/internal/middleware"
	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/Adi-ty/betteruptime-go/migrations"
)

type Application struct {
	Logger         *log.Logger
	WebsiteHandler *api.WebsiteHandler
	UserHandler    *api.UserHandler
	Middleware     *middleware.UserMiddleware
	DB             *sql.DB
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
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	
	// Handlers
	websiteHandler := api.NewWebsiteHandler(websiteStore, logger)
	userHandler := api.NewUserHandler(userStore, tokenStore, logger)
	middleware := middleware.NewUserMiddleware(userStore)

	app := &Application{
		Logger: logger,
		WebsiteHandler: websiteHandler,
		UserHandler: userHandler,
		Middleware: middleware,
		DB: pgDB,
	}
	return app, nil
}