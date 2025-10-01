// internal/utils/response.go
package utils

import (
    "encoding/json"
    "log"
    "net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
    if code > 499 {
        log.Printf("Responding with 5XX error: %s", message)
    }
    
    type errorResponse struct {
        Error string `json:"error"`
    }
    
    RespondWithJSON(w, code, errorResponse{
        Error: message,
    })
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Failed to marshal JSON response: %v", err)
        w.WriteHeader(500)
        return
    }
    
    w.WriteHeader(code)
    w.Write(data)
}