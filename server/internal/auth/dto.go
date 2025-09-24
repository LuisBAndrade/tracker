package auth

import (
	"time"

	"github.com/google/uuid"
)


type UserDTO struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User UserDTO `json:"user"`
}