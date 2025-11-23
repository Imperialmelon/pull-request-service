package utils

import (
	"encoding/json"
	"net/http"

	"github.com/imperialmelon/avito/internal/models"
)

func WriteErrorJSON(w http.ResponseWriter, statusCode int, code models.ErrorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:    string(code),
			Message: message,
		},
	})
}
