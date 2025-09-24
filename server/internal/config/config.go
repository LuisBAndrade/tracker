package server

import (
	"github.com/LuisBAndrade/tracker/server/db/internal/database"
)

type ApiConfig struct {
	db *database.Queries
}