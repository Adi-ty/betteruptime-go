package routes

import (
	"net/http"

	"github.com/Adi-ty/betteruptime-go/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Group(func (r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Get("/status/{website_id}", app.WebsiteHandler.HandleGetWebsiteStatus)
		r.Post("/website", app.WebsiteHandler.HandleCreateWebsite)
	})

	router.Post("/user/register", app.UserHandler.HandleUserRegister)
	router.Post("/user/login", app.UserHandler.HandleUserLogin)

	return router
}