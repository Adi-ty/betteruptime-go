package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Adi-ty/betteruptime-go/internal/middleware"
	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/go-chi/chi/v5"
)

type WebsiteHandler struct {
	logger *log.Logger
	websiteStore store.WebsiteStore
}

func NewWebsiteHandler(websiteStore store.WebsiteStore, logger *log.Logger) *WebsiteHandler {
	return &WebsiteHandler{
		logger: logger,
		websiteStore: websiteStore,
	}
}

func (h *WebsiteHandler) HandleGetWebsiteStatus(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "website_id")
	if idParam == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	website, err := h.websiteStore.GetWebsiteStatusByID(currentUser.ID, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get website: %v", err), http.StatusInternalServerError)
		return
	}
	if website == nil {
		http.Error(w, "website not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(website)
}

func (h *WebsiteHandler) HandleCreateWebsite(w http.ResponseWriter, r *http.Request) {
	var website store.Website
	decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&website)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	website.UserID = currentUser.ID
	website.TimeAdded =  time.Now()

	err = h.websiteStore.CreateWebsite(&website)
	if err != nil {
		http.Error(w, "Failed to create website", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": website.ID,
		"url": website.Url,
		"time_added": website.TimeAdded,
	})
}