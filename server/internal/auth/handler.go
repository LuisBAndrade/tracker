package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return 
	}

	user, err := h.svc.RegisterUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't register user"})
		return 
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return 
	}

	token, user, err := h.svc.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return 
	}

	c.SetCookie(
		"session",
		token,
		3600,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, user)
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil  || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return 
	}

	newAccess, user, err := h.svc.RefreshSession(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return 
	}

	c.SetCookie(
		"access_token",
		newAccess,
		int((15 * time.Minute).Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "access token refreshed",
		"user": user,
	})
}