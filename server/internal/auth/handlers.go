// internal/auth/handlers.go
package auth

import (
    "encoding/json"
    "net/http"
    "time"
    
   "github.com/LuisBAndrade/etracker/internal/utils"
)

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    User    UserResponse `json:"user"`
    Message string       `json:"message"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func (s *Service) HandleRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    if err := utils.ValidateStruct(req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    user, err := s.Register(r.Context(), req.Email, req.Password)
    if err != nil {
        if err == ErrUserExists {
            utils.RespondWithError(w, http.StatusConflict, "User already exists")
            return
        }
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }

    utils.RespondWithJSON(w, http.StatusCreated, AuthResponse{
        User: UserResponse{
            ID:        user.ID.String(),
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
        Message: "User created successfully",
    })
}

func (s *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    if err := utils.ValidateStruct(req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    user, token, err := s.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        if err == ErrInvalidCredentials {
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
            return
        }
        utils.RespondWithError(w, http.StatusInternalServerError, "Login failed")
        return
    }

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure: false,                // must be false for localhost HTTP
		SameSite: http.SameSiteLaxMode, // use Lax for local dev
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
    http.SetCookie(w, cookie)

    utils.RespondWithJSON(w, http.StatusOK, AuthResponse{
        User: UserResponse{
            ID:        user.ID.String(),
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
        Message: "Login successful",
    })
}

func (s *Service) HandleLogout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "No session found")
        return
    }

    if err := s.Logout(r.Context(), cookie.Value); err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Logout failed")
        return
    }

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Logged out successfully",
    })
}

func (s *Service) HandleMe(w http.ResponseWriter, r *http.Request) {
    user, ok := GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, UserResponse{
        ID:        user.ID.String(),
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    })
}

func (s *Service) HandleLogoutAll(w http.ResponseWriter, r *http.Request) {
    user, ok := GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
        return
    }

    if err := s.LogoutAll(r.Context(), user.ID); err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to logout all sessions")
        return
    }

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Logged out from all devices",
    })
}