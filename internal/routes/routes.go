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

	router.Get("/website/{id}", app.WebsiteHandler.HandleGetWebsite)
	router.Post("/website", app.WebsiteHandler.HandleCreateWebsite)

	return router
}