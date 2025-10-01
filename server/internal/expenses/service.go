package expenses

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/LuisBAndrade/etracker/internal/database"
)

type Service struct {
    queries *database.Queries
}

func NewService(queries *database.Queries) *Service {
    return &Service{queries: queries}
}

func (s *Service) CreateExpense(ctx context.Context, userID uuid.UUID, categoryID *uuid.UUID, amount string, description string, date time.Time) (*database.Expense, error) {
    var nullCategoryID uuid.NullUUID 
    if categoryID != nil {
        nullCategoryID = uuid.NullUUID{UUID: *categoryID, Valid: true}
    }
    
    expense, err := s.queries.CreateExpense(ctx, database.CreateExpenseParams{
        UserID:      userID,
        CategoryID:  nullCategoryID, 
        Amount:      amount,       
        Description: description,
        Date:        date,
    })
    return &expense, err
}

func (s *Service) GetUserExpenses(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]database.GetExpensesByUserRow, error) {
    return s.queries.GetExpensesByUser(ctx, database.GetExpensesByUserParams{
        UserID: userID,
        Limit:  limit,
        Offset: offset,
    })
}

func (s *Service) GetExpensesByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]database.GetExpensesByUserAndDateRangeRow, error) {
    return s.queries.GetExpensesByUserAndDateRange(ctx, database.GetExpensesByUserAndDateRangeParams{
        UserID: userID,
        Date:   startDate,
        Date_2: endDate,
    })
}

func (s *Service) GetExpenseByID(ctx context.Context, expenseID, userID uuid.UUID) (*database.GetExpenseByIDRow, error) {
    expense, err := s.queries.GetExpenseByID(ctx, database.GetExpenseByIDParams{
        ID:     expenseID,
        UserID: userID,
    })
    return &expense, err
}

func (s *Service) UpdateExpense(ctx context.Context, expenseID, userID uuid.UUID, categoryID *uuid.UUID, amount string, description string, date time.Time) (*database.Expense, error) {
    var nullCategoryID uuid.NullUUID
    if categoryID != nil {
        nullCategoryID = uuid.NullUUID{UUID: *categoryID, Valid: true}
    }
    
    expense, err := s.queries.UpdateExpense(ctx, database.UpdateExpenseParams{
        ID:          expenseID,
        UserID:      userID,
        Amount:      amount,
        Description: description,
        CategoryID:  nullCategoryID,
        Date:        date,
    })
    return &expense, err
}

func (s *Service) DeleteExpense(ctx context.Context, expenseID, userID uuid.UUID) error {
    return s.queries.DeleteExpense(ctx, database.DeleteExpenseParams{
        ID:     expenseID,
        UserID: userID,
    })
}

func (s *Service) GetExpenseTotal(ctx context.Context, userID uuid.UUID) (string, error) {
    result, err := s.queries.GetExpenseTotalByUser(ctx, userID)
    if err != nil {
        return "", err
    }
    
    // Type assert to string
    total, ok := result.(string)
    if !ok {
        return "0.00", nil // Return default value if assertion fails
    }
    
    return total, nil
}

func (s *Service) GetExpenseTotalByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (string, error) {
    result, err := s.queries.GetExpenseTotalByUserAndDateRange(ctx, database.GetExpenseTotalByUserAndDateRangeParams{
        UserID: userID,
        Date:   startDate,
        Date_2: endDate,
    })
    if err != nil {
        return "", err
    }
    
    // Type assert to string
    total, ok := result.(string)
    if !ok {
        return "0.00", nil // Return default value if assertion fails
    }
    
    return total, nil
}
func (s *Service) GetExpensesByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]database.GetExpensesByCategoryRow, error) {
    return s.queries.GetExpensesByCategory(ctx, database.GetExpensesByCategoryParams{
        UserID: userID,
        Date:   startDate,
        Date_2: endDate,
    })
}