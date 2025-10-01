package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/LuisBAndrade/etracker/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserExists        = errors.New("user already exists")
    ErrInvalidSession    = errors.New("invalid session")
)

type Service struct {
    queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
    return &Service{
        queries: queries,
    }
}

func (s *Service) Register(ctx context.Context, email, password string) (*database.User, error) {
    // Check if user exists
    _, err := s.queries.GetUserByEmail(ctx, email)
    if err == nil {
        return nil, ErrUserExists
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Create user
    user, err := s.queries.CreateUser(ctx, database.CreateUserParams{
        Email:          email,
        HashedPassword: string(hashedPassword),
    })
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*database.User, string, error) {
    // Get user by email
    user, err := s.queries.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, "", ErrInvalidCredentials
    }

    // Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
    if err != nil {
        return nil, "", ErrInvalidCredentials
    }

    // Create session token
    token, err := s.generateSessionToken()
    if err != nil {
        return nil, "", err
    }

    // Save session
    expiresAt := time.Now().Add(24 * time.Hour * 7) // 7 days
    err = s.queries.CreateSession(ctx, database.CreateSessionParams{
        Token:     token,
        UserID:    user.ID,
        ExpiresAt: expiresAt,
    })
    if err != nil {
        return nil, "", err
    }

    return &user, token, nil
}

func (s *Service) GetUserBySession(ctx context.Context, token string) (*database.User, error) {
    user, err := s.queries.GetUserBySessionToken(ctx, token)
    if err != nil {
        return nil, ErrInvalidSession
    }
    return &user, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
    return s.queries.RevokeSession(ctx, token)
}

func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
    return s.queries.RevokeAllUserSessions(ctx, userID)
}

func (s *Service) generateSessionToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
    return s.queries.CleanupExpiredSessions(ctx)
}