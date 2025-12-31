package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/Adi-ty/betteruptime-go/internal/tokens"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Token    string `json:"token"`
}


type UserHandler struct {
	tokenStore store.TokenStore
	userStore store.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		tokenStore: tokenStore,
		logger: logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if (req.Password) == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (h *UserHandler) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding the body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		h.logger.Printf("Validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &store.User{
		Username: req.Username,
		Email: req.Email,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error setting password hash: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error creating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("Error creating token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	UserResponse := UserResponse{
		ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        Token:    token.Plaintext,
	}
	json.NewEncoder(w).Encode(UserResponse)
}

func (h *UserHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest
	decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    err := decoder.Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding the body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil {
		h.logger.Printf("Error fetching user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	match, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("Error checking password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("Error creating token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	UserResponse := UserResponse{
		ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        Token:    token.Plaintext,
	}
	json.NewEncoder(w).Encode(UserResponse)
}