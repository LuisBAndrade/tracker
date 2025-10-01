package expenses

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LuisBAndrade/etracker/internal/auth"
	"github.com/LuisBAndrade/etracker/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateExpenseRequest struct {
    CategoryID  *string `json:"category_id"`
    Amount      float64 `json:"amount" validate:"required"` // Keep as float64 for JSON
    Description string  `json:"description" validate:"required"`
    Date        string  `json:"date"` // YYYY-MM-DD format
}

type UpdateExpenseRequest struct {
    CategoryID  *string `json:"category_id"`
    Amount      float64 `json:"amount" validate:"required"` // Keep as float64 for JSON
    Description string  `json:"description" validate:"required"`
    Date        string  `json:"date"` // YYYY-MM-DD format
}

type ExpenseResponse struct {
    ID           string  `json:"id"`
    CategoryID   *string `json:"category_id"`
    CategoryName *string `json:"category_name"`
    CategoryColor *string `json:"category_color"`
    Amount       string  `json:"amount"` // String from database
    Description  string  `json:"description"`
    Date         string  `json:"date"`
    CreatedAt    string  `json:"created_at"`
    UpdatedAt    string  `json:"updated_at"`
}

type ExpenseSummaryResponse struct {
    Total       string `json:"total"` // String from database
    Count       int    `json:"count"`
    Expenses    []ExpenseResponse `json:"expenses"`
}

type CategorySummaryResponse struct {
    CategoryID    string `json:"category_id"`
    CategoryName  string `json:"category_name"`
    CategoryColor string `json:"category_color"`
    TotalAmount   string `json:"total_amount"` // String from database
    ExpenseCount  int64  `json:"expense_count"`
}

func (s *Service) HandleCreateExpense(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    var req CreateExpenseRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    if err := utils.ValidateStruct(req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    if req.Amount <= 0 {
        utils.RespondWithError(w, http.StatusBadRequest, "Amount must be greater than 0")
        return
    }
    
    // Convert float64 to string for database
    amountStr := strconv.FormatFloat(req.Amount, 'f', 2, 64)
    
    // Parse date
    var date time.Time
    var err error
    if req.Date != "" {
        date, err = time.Parse("2006-01-02", req.Date)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
            return
        }
    } else {
        date = time.Now()
    }
    
    // Parse category ID
    var categoryID *uuid.UUID
    if req.CategoryID != nil && *req.CategoryID != "" {
        parsedID, err := uuid.Parse(*req.CategoryID)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
            return
        }
        categoryID = &parsedID
    }
    
    expense, err := s.CreateExpense(r.Context(), user.ID, categoryID, amountStr, req.Description, date)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create expense")
        return
    }
    
    response := ExpenseResponse{
        ID:          expense.ID.String(),
        Amount:      expense.Amount,
        Description: expense.Description,
        Date:        expense.Date.Format("2006-01-02"),
        CreatedAt:   expense.CreatedAt.Format("2006-01-02T15:04:05Z"),
        UpdatedAt:   expense.UpdatedAt.Format("2006-01-02T15:04:05Z"),
    }
    
    if expense.CategoryID.Valid {
        categoryIDStr := expense.CategoryID.UUID.String()
        response.CategoryID = &categoryIDStr
    }
    
    utils.RespondWithJSON(w, http.StatusCreated, response)
}

func (s *Service) HandleGetExpenses(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    // Parse query parameters
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    startDateStr := r.URL.Query().Get("start_date")
    endDateStr := r.URL.Query().Get("end_date")
    
    limit := int32(20) // default
    offset := int32(0) // default
    
    if limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
            limit = int32(l)
        }
    }
    
    if offsetStr != "" {
        if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
            offset = int32(o)
        }
    }
    
    // If date range is provided, use date range query
    if startDateStr != "" && endDateStr != "" {
        startDate, err := time.Parse("2006-01-02", startDateStr)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid start_date format")
            return
        }
        
        endDate, err := time.Parse("2006-01-02", endDateStr)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid end_date format")
            return
        }
        
        expenses, err := s.GetExpensesByDateRange(r.Context(), user.ID, startDate, endDate)
        if err != nil {
            utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get expenses")
            return
        }
        
        response := make([]ExpenseResponse, len(expenses))
        
        for i, exp := range expenses {
            response[i] = ExpenseResponse{
                ID:          exp.ID.String(),
                Amount:      exp.Amount,
                Description: exp.Description,
                Date:        exp.Date.Format("2006-01-02"),
                CreatedAt:   exp.CreatedAt.Format("2006-01-02T15:04:05Z"),
                UpdatedAt:   exp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
            }
            
            if exp.CategoryID.Valid {
                categoryIDStr := exp.CategoryID.UUID.String()
                response[i].CategoryID = &categoryIDStr
            }
            if exp.CategoryName.Valid {
                categoryName := exp.CategoryName.String
                response[i].CategoryName = &categoryName
            }
            if exp.CategoryColor.Valid {
                categoryColor := exp.CategoryColor.String
                response[i].CategoryColor = &categoryColor
            }
        }
        
        // Get total for date range
        total, _ := s.GetExpenseTotalByDateRange(r.Context(), user.ID, startDate, endDate)
        
        utils.RespondWithJSON(w, http.StatusOK, ExpenseSummaryResponse{
            Total:    total,
            Count:    len(expenses),
            Expenses: response,
        })
        return
    }
    
    // Regular pagination query
    expenses, err := s.GetUserExpenses(r.Context(), user.ID, limit, offset)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get expenses")
        return
    }
    
    response := make([]ExpenseResponse, len(expenses))
    
    for i, exp := range expenses {
        response[i] = ExpenseResponse{
            ID:          exp.ID.String(),
            Amount:      exp.Amount,
            Description: exp.Description,
            Date:        exp.Date.Format("2006-01-02"),
            CreatedAt:   exp.CreatedAt.Format("2006-01-02T15:04:05Z"),
            UpdatedAt:   exp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
        }
        
        if exp.CategoryID.Valid {
            categoryIDStr := exp.CategoryID.UUID.String()
            response[i].CategoryID = &categoryIDStr
        }
        if exp.CategoryName.Valid {
            categoryName := exp.CategoryName.String
            response[i].CategoryName = &categoryName
        }
        if exp.CategoryColor.Valid {
            categoryColor := exp.CategoryColor.String
            response[i].CategoryColor = &categoryColor
        }
    }
    
    // Get total for user
    total, _ := s.GetExpenseTotal(r.Context(), user.ID)
    
    utils.RespondWithJSON(w, http.StatusOK, ExpenseSummaryResponse{
        Total:    total,
        Count:    len(expenses),
        Expenses: response,
    })
}

func (s *Service) HandleUpdateExpense(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    vars := mux.Vars(r)
    expenseID, err := uuid.Parse(vars["id"])
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid expense ID")
        return
    }
    
    var req UpdateExpenseRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    if err := utils.ValidateStruct(req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    if req.Amount <= 0 {
        utils.RespondWithError(w, http.StatusBadRequest, "Amount must be greater than 0")
        return
    }
    
    // Convert float64 to string for database
    amountStr := strconv.FormatFloat(req.Amount, 'f', 2, 64)
    
    // Parse date
    var date time.Time
    if req.Date != "" {
        date, err = time.Parse("2006-01-02", req.Date)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
            return
        }
    } else {
        date = time.Now()
    }
    
    // Parse category ID
    var categoryID *uuid.UUID
    if req.CategoryID != nil && *req.CategoryID != "" {
        parsedID, err := uuid.Parse(*req.CategoryID)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
            return
        }
        categoryID = &parsedID
    }
    
    expense, err := s.UpdateExpense(r.Context(), expenseID, user.ID, categoryID, amountStr, req.Description, date)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update expense")
        return
    }
    
    response := ExpenseResponse{
        ID:          expense.ID.String(),
        Amount:      expense.Amount,
        Description: expense.Description,
        Date:        expense.Date.Format("2006-01-02"),
        CreatedAt:   expense.CreatedAt.Format("2006-01-02T15:04:05Z"),
        UpdatedAt:   expense.UpdatedAt.Format("2006-01-02T15:04:05Z"),
    }
    
    if expense.CategoryID.Valid {
        categoryIDStr := expense.CategoryID.UUID.String()
        response.CategoryID = &categoryIDStr
    }
    
    utils.RespondWithJSON(w, http.StatusOK, response)
}

func (s *Service) HandleDeleteExpense(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    vars := mux.Vars(r)
    expenseID, err := uuid.Parse(vars["id"])
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid expense ID")
        return
    }
    
    if err := s.DeleteExpense(r.Context(), expenseID, user.ID); err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete expense")
        return
    }
    
    utils.RespondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Expense deleted successfully",
    })
}

func (s *Service) HandleGetExpensesByCategory(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    startDateStr := r.URL.Query().Get("start_date")
    endDateStr := r.URL.Query().Get("end_date")
    
    var startDate, endDate time.Time
    var err error
    
    if startDateStr == "" || endDateStr == "" {
        // Default to current month
        now := time.Now()
        startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
        endDate = startDate.AddDate(0, 1, -1)
    } else {
        startDate, err = time.Parse("2006-01-02", startDateStr)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid start_date format")
            return
        }
        
        endDate, err = time.Parse("2006-01-02", endDateStr)
        if err != nil {
            utils.RespondWithError(w, http.StatusBadRequest, "Invalid end_date format")
            return
        }
    }
    
    categories, err := s.GetExpensesByCategory(r.Context(), user.ID, startDate, endDate)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get expenses by category")
        return
    }
    
    response := make([]CategorySummaryResponse, len(categories))
    for i, cat := range categories {
        // Type assert TotalAmount to string
        totalAmount, ok := cat.TotalAmount.(string)
        if !ok {
            utils.RespondWithError(w, http.StatusInternalServerError, "Invalid total amount format")
            return
        }
        
        response[i] = CategorySummaryResponse{
            CategoryID:    cat.CategoryID.String(),
            CategoryName:  cat.CategoryName,
            CategoryColor: cat.CategoryColor,
            TotalAmount:   totalAmount,
            ExpenseCount:  cat.ExpenseCount,
        }
    }
    
    utils.RespondWithJSON(w, http.StatusOK, response)
}