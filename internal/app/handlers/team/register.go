package team

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Register(r *mux.Router, teamUC TeamUseCase) {
	h := NewHandler(teamUC)

	r.HandleFunc("/team/add", h.CreateTeam).Methods(http.MethodPost)
	r.HandleFunc("/team/get/{team_name}", h.GetTeam).Methods(http.MethodGet)
}
