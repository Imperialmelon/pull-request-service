package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	svcerrors "github.com/Imperialmelon/AvitoTechTest/internal/errors"
	"github.com/Imperialmelon/AvitoTechTest/internal/models"
	"github.com/Imperialmelon/AvitoTechTest/internal/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	userUC UserUseCase
}

func NewHandler(userUC UserUseCase) *Handler {
	return &Handler{userUC: userUC}
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) string {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return idStr
}

func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req models.SetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		utils.WriteErrorJSON(w, http.StatusBadRequest, models.ErrorNotFound, "invalid request body")
		return
	}

	updatedUser, err := h.userUC.SetIsActive(req)
	if err != nil {
		switch {
		case errors.Is(err, svcerrors.ErrorNotFound):
			log.Println("user not found:", err)
			utils.WriteErrorJSON(w, http.StatusNotFound, models.ErrorNotFound, "resource not found")
			return
		default:
			log.Println("internal error:", err)
			utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update user")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": updatedUser,
	})
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	id := h.parseIDFromURL(r, "id")
	userWithPRs, err := h.userUC.GetReview(id)
	if err != nil {
		log.Println("internal error:", err)
		utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get reviews")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userWithPRs)
}
