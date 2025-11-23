package pr

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	svcerrors "github.com/Imperialmelon/AvitoTechTest/internal/errors"
	"github.com/Imperialmelon/AvitoTechTest/internal/models"
	"github.com/Imperialmelon/AvitoTechTest/internal/utils"
)

type Handler struct {
	prUC PRUseCase ``
}

func NewHandler(prUC PRUseCase) *Handler {
	return &Handler{prUC: prUC}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		utils.WriteErrorJSON(w, http.StatusBadRequest, models.ErrorNotFound, "invalid request body")
		return
	}

	createdPR, err := h.prUC.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, svcerrors.ErrorNotFound):
			log.Println("author/team not found:", err)
			utils.WriteErrorJSON(w, http.StatusNotFound, models.ErrorNotFound, "resource not found")
			return
		case errors.Is(err, svcerrors.ErrorPRExists):
			log.Println("pr exists:", err)
			utils.WriteErrorJSON(w, http.StatusConflict, models.ErrorPRExists, "pr exists")
			return
		default:
			log.Println("internal error:", err)
			utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create pr")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": createdPR,
	})
}

func (h *Handler) Merge(w http.ResponseWriter, r *http.Request) {
	var req models.PRID
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		utils.WriteErrorJSON(w, http.StatusBadRequest, models.ErrorNotFound, "invalid request body")
		return
	}
	mergedPR, err := h.prUC.Merge(req.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, svcerrors.ErrorNotFound):
			log.Println("pr:", err)
			utils.WriteErrorJSON(w, http.StatusNotFound, models.ErrorNotFound, "resource not found")
			return
		default:
			log.Println("internal error:", err)
			utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create pr")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": mergedPR,
	})
}

func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req models.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		utils.WriteErrorJSON(w, http.StatusBadRequest, models.ErrorNotFound, "invalid request body")
		return
	}
	reassignedPR, err := h.prUC.Reassign(req.PullRequestID, req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, svcerrors.ErrorNotFound):
			log.Println("pr:", err)
			utils.WriteErrorJSON(w, http.StatusNotFound, models.ErrorNotFound, "resource not found")
			return
		case errors.Is(err, svcerrors.ErrorPRMerged):
			log.Println("pr merged:", err)
			utils.WriteErrorJSON(w, http.StatusConflict, models.ErrorPRMerged, "pr merged")
			return
		default:
			log.Println("internal error:", err)
			utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create pr")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reassignedPR)
}
