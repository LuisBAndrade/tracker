package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/LuisBAndrade/tracker/server/db/internal/database"
	"github.com/google/uuid"
)

type Service struct {
	q *database.Queries
}

func NewService(q *database.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) RegisterUser(ctx context.Context, email, password string) (UserDTO, error) {
	hashed, err := HashPassword(password)
	if err != nil {
		return UserDTO{}, err
	}
	
	_, err = s.q.CreateUser(ctx, database.CreateUserParams{
		Email: email,
		HashedPassword: hashed,
	})
	if err != nil {
		return UserDTO{}, err
	}

	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return UserDTO{}, err
	}

	return UserDTO{
		ID: u.ID,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil 
	
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, UserDTO, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return "", UserDTO{}, err
	}

	if err := CheckPassword(password, u.HashedPassword); err != nil {
		return "", UserDTO{}, err
	}

	token, err := MakeJWT(uuid.MustParse(u.ID), os.Getenv("JWT_SECRET"), 24*time.Hour)
	if err != nil {
		return "", UserDTO{}, err
	}

	dto := UserDTO{
		ID: u.ID,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	return token, dto, nil
}

func (s *Service) RefreshSession(ctx context.Context, refreshToken string) (string, UserDTO, error) {
	rt, err := s.q.GetRefreshToken(ctx, refreshToken)
	if err != nil || rt.ExpiresAt.Before(time.Now()){
		return "", UserDTO{}, errors.New("invalid or expired refresh token")
	}

	u, err := s.q.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return "", UserDTO{}, err
	}

	newAccessToken, err := MakeJWT(uuid.MustParse(u.ID), os.Getenv("JWT_SECRET"), time.Minute)
	if err != nil {
		return "", UserDTO{}, err
	}

	dto := UserDTO{
		ID: u.ID,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	return newAccessToken, dto, nil
}