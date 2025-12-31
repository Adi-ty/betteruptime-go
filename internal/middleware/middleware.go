package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Adi-ty/betteruptime-go/internal/store"
	"github.com/Adi-ty/betteruptime-go/internal/tokens"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

func NewUserMiddleware(userStore store.UserStore) *UserMiddleware {
	return &UserMiddleware{UserStore: userStore}
}

type contextKey string

const userContextKey = contextKey("userkey")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		return nil
	}
	return user
}

func (m *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(header, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		user, err := m.UserStore.GetUserByToken(tokens.ScopeAuth, token)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		if user == nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}