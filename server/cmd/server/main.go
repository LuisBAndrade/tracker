package main

import (
	"log"
	"os"

	"github.com/LuisBAndrade/tracker/server/db/internal/database"
	"github.com/LuisBAndrade/tracker/server/internal/auth"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")

	dbConn, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbConn.Close()

	q := database.New(dbConn)

	authService := auth.NewService(q)
	authHandler := auth.NewHandler(authService)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)
	r.POST("/refresh", authHandler.Refresh)

	r.GET("/transactions", auth.JWTMiddleware())

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}