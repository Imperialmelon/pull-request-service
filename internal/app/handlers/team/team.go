package team

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
	teamUC TeamUseCase
}

func NewHandler(teamUC TeamUseCase) *Handler {
	return &Handler{teamUC: teamUC}
}

func (h *Handler) parseIDFromURL(r *http.Request, paramName string) string {
	vars := mux.Vars(r)
	idStr := vars[paramName]
	return idStr
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req models.TeamApi
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		utils.WriteErrorJSON(w, http.StatusBadRequest, models.ErrorNotFound, "invalid request body")
		return
	}

	createdTeam, err := h.teamUC.Add(req)
	if err != nil {
		switch {
		case errors.Is(err, svcerrors.ErrorTeamExists):
			log.Println("team already exists:", err)
			utils.WriteErrorJSON(w, http.StatusBadRequest, models.ErrorTeamExists, "team already exists")
			return
		default:
			log.Println("internal error:", err)
			utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create team")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"team": createdTeam,
	})
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := h.parseIDFromURL(r, "team_name")

	createdTeam, err := h.teamUC.Get(teamName)
	if err != nil {
		switch {
		case errors.Is(err, svcerrors.ErrorNotFound):
			log.Println("team not found:", err)
			utils.WriteErrorJSON(w, http.StatusNotFound, models.ErrorNotFound, "team not found")
			return
		default:
			log.Println("internal error:", err)
			utils.WriteErrorJSON(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get team")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createdTeam)
}
