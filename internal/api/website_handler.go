package api

import (
	"log"
	"net/http"
)

type WebsiteHandler struct {
	logger *log.Logger
}

func NewWebsiteHandler(logger *log.Logger) *WebsiteHandler {
	return &WebsiteHandler{
		logger: logger,
	}
}

func (h *WebsiteHandler) HandleGetWebsite(w http.ResponseWriter, r *http.Request) {
	
}

func (h *WebsiteHandler) HandleCreateWebsite(w http.ResponseWriter, r *http.Request) {
	
}