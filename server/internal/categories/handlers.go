// internal/categories/handlers.go
package categories

import (
	"encoding/json"
	"net/http"

	"github.com/LuisBAndrade/etracker/internal/auth"
	"github.com/LuisBAndrade/etracker/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateCategoryRequest struct {
    Name  string `json:"name" validate:"required"`
    Color string `json:"color"`
}

type UpdateCategoryRequest struct {
    Name  string `json:"name" validate:"required"`
    Color string `json:"color"`
}

type CategoryResponse struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Color     string `json:"color"`
    CreatedAt string `json:"created_at"`
}

func (s *Service) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    var req CreateCategoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    if err := utils.ValidateStruct(req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    category, err := s.CreateCategory(r.Context(), user.ID, req.Name, req.Color)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create category")
        return
    }
    
    utils.RespondWithJSON(w, http.StatusCreated, CategoryResponse{
        ID:        category.ID.String(),
        Name:      category.Name,
        Color:     category.Color,
        CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z"),
    })
}

func (s *Service) HandleGetCategories(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    categories, err := s.GetUserCategories(r.Context(), user.ID)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get categories")
        return
    }
    
    response := make([]CategoryResponse, len(categories))
    for i, cat := range categories {
        response[i] = CategoryResponse{
            ID:        cat.ID.String(),
            Name:      cat.Name,
            Color:     cat.Color,
            CreatedAt: cat.CreatedAt.Format("2006-01-02T15:04:05Z"),
        }
    }
    
    utils.RespondWithJSON(w, http.StatusOK, response)
}

func (s *Service) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    vars := mux.Vars(r)
    categoryID, err := uuid.Parse(vars["id"])
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
        return
    }
    
    var req UpdateCategoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    if err := utils.ValidateStruct(req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    category, err := s.UpdateCategory(r.Context(), categoryID, user.ID, req.Name, req.Color)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update category")
        return
    }
    
    utils.RespondWithJSON(w, http.StatusOK, CategoryResponse{
        ID:        category.ID.String(),
        Name:      category.Name,
        Color:     category.Color,
        CreatedAt: category.CreatedAt.Format("2006-01-02T15:04:05Z"),
    })
}

func (s *Service) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUserFromContext(r.Context())
    if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
        return
    }
    
    vars := mux.Vars(r)
    categoryID, err := uuid.Parse(vars["id"])
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
        return
    }
    
    if err := s.DeleteCategory(r.Context(), categoryID, user.ID); err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete category")
        return
    }
    
    utils.RespondWithJSON(w, http.StatusOK, map[string]string{
        "message": "Category deleted successfully",
    })
}