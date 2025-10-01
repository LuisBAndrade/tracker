package categories

import (
	"context"

	"github.com/LuisBAndrade/etracker/internal/database"
	"github.com/google/uuid"
)

type Service struct {
    queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
    return &Service{queries: queries}
}

func (s *Service) CreateCategory(ctx context.Context, userID uuid.UUID, name, color string) (*database.Category, error) {
    if color == "" {
        color = "#6B7280" // Default gray color
    }
    
    category, err := s.queries.CreateCategory(ctx, database.CreateCategoryParams{
        UserID: userID,
        Name:   name,
        Color:  color,
    })
    return &category, err
}

func (s *Service) GetUserCategories(ctx context.Context, userID uuid.UUID) ([]database.Category, error) {
    return s.queries.GetCategoriesByUser(ctx, userID)
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID, userID uuid.UUID) (*database.Category, error) {
    category, err := s.queries.GetCategoryByID(ctx, database.GetCategoryByIDParams{
        ID:     categoryID,
        UserID: userID,
    })
    return &category, err
}

func (s *Service) UpdateCategory(ctx context.Context, categoryID, userID uuid.UUID, name, color string) (*database.Category, error) {
    category, err := s.queries.UpdateCategory(ctx, database.UpdateCategoryParams{
        ID:     categoryID,
        Name:   name,
        Color:  color,
        UserID: userID,
    })
    return &category, err
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID, userID uuid.UUID) error {
    return s.queries.DeleteCategory(ctx, database.DeleteCategoryParams{
        ID:     categoryID,
        UserID: userID,
    })
}

// CreateDefaultCategories creates default categories for new users
func (s *Service) CreateDefaultCategories(ctx context.Context, userID uuid.UUID) error {
    defaults := []struct {
        name  string
        color string
    }{
        {"Food & Dining", "#EF4444"},
        {"Transportation", "#3B82F6"},
        {"Shopping", "#8B5CF6"},
        {"Entertainment", "#F59E0B"},
        {"Bills & Utilities", "#10B981"},
        {"Healthcare", "#EC4899"},
        {"Other", "#6B7280"},
    }
    
    for _, cat := range defaults {
        _, err := s.CreateCategory(ctx, userID, cat.name, cat.color)
        if err != nil {
            return err
        }
    }
    return nil
}