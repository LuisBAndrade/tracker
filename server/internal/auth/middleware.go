// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"

	"github.com/LuisBAndrade/etracker/internal/database"
	"github.com/LuisBAndrade/etracker/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get session token from cookie
        cookie, err := r.Cookie("session_token")
        if err != nil {
            utils.RespondWithError(w, http.StatusUnauthorized, "Missing session token")
            return
        }

        // Get user by session token
        user, err := s.GetUserBySession(r.Context(), cookie.Value)
        if err != nil {
            // Clear invalid cookie
            http.SetCookie(w, &http.Cookie{
                Name:     "session_token",
                Value:    "",
                Path:     "/",
                MaxAge:   -1,
                HttpOnly: true,
            })
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid session")
            return
        }

        // Add user to context
        ctx := context.WithValue(r.Context(), UserContextKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetUserFromContext(ctx context.Context) (*database.User, bool) {
    user, ok := ctx.Value(UserContextKey).(*database.User)
    return user, ok
}