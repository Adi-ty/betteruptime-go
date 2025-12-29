package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WebsiteHandler struct {
	logger *log.Logger
}

func NewWebsiteHandler(logger *log.Logger) *WebsiteHandler {
	return &WebsiteHandler{
		logger: logger,
	}
}

type CreateWebsiteRequest struct {
	URL string `json:"url"`
}

func (h *WebsiteHandler) HandleGetWebsite(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "user id=%d", id)
}

func (h *WebsiteHandler) HandleCreateWebsite(w http.ResponseWriter, r *http.Request) {
	var website CreateWebsiteRequest
	err := json.NewDecoder(r.Body).Decode(&website)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Add logic to create the website in the store

	w.Write([]byte("Website created"))
}