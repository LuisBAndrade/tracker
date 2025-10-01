package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/LuisBAndrade/etracker/internal/auth"
	"github.com/LuisBAndrade/etracker/internal/categories"
	"github.com/LuisBAndrade/etracker/internal/config"
	"github.com/LuisBAndrade/etracker/internal/database"
	"github.com/LuisBAndrade/etracker/internal/expenses"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
)

func main() {
    cfg := config.Load()

    conn, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer conn.Close()

    if err := conn.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    log.Println("Connected to database successfully")

    queries := database.New(conn)

    authService := auth.NewService(queries)
    categoriesService := categories.NewService(queries)
    expensesService := expenses.NewService(queries)

    router := mux.NewRouter()

    // Auth routes
    router.HandleFunc("/api/auth/register", authService.HandleRegister).Methods("POST")
    router.HandleFunc("/api/auth/login", authService.HandleLogin).Methods("POST")
    router.HandleFunc("/api/auth/logout", authService.HandleLogout).Methods("POST")

    // Protected routes
    protected := router.PathPrefix("/api").Subrouter()
    protected.Use(authService.AuthMiddleware)

    protected.HandleFunc("/auth/me", authService.HandleMe).Methods("GET")
    protected.HandleFunc("/auth/logout-all", authService.HandleLogoutAll).Methods("POST")

    protected.HandleFunc("/categories", categoriesService.HandleCreateCategory).Methods("POST")
    protected.HandleFunc("/categories", categoriesService.HandleGetCategories).Methods("GET")
    protected.HandleFunc("/categories/{id}", categoriesService.HandleUpdateCategory).Methods("PUT")
    protected.HandleFunc("/categories/{id}", categoriesService.HandleDeleteCategory).Methods("DELETE")

    protected.HandleFunc("/expenses", expensesService.HandleCreateExpense).Methods("POST")
    protected.HandleFunc("/expenses", expensesService.HandleGetExpenses).Methods("GET")
    protected.HandleFunc("/expenses/{id}", expensesService.HandleUpdateExpense).Methods("PUT")
    protected.HandleFunc("/expenses/{id}", expensesService.HandleDeleteExpense).Methods("DELETE")
    protected.HandleFunc("/expenses/by-category", expensesService.HandleGetExpensesByCategory).Methods("GET")

    // ✅ Apply CORS *after* all routes are mounted
    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"http://localhost:5173", "http://3.91.219.223:5173"}), // your frontend dev server
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Accept", "X-Requested-With"}),
        handlers.AllowCredentials(),
    )(router)

    server := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      corsHandler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

    log.Printf("Server starting on port %s", cfg.Port)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}
