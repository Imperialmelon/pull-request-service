package handlers

import (
	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers/pr"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers/team"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers/user"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/usecase"
	"github.com/gorilla/mux"
)

type Handler struct {
	userHandler *user.Handler
	teamHandler *team.Handler
	prHandler   *pr.Handler
	registrator *Registrator
}

func NewHandler(uc *usecase.UseCase) *Handler {
	return &Handler{
		userHandler: user.NewHandler(uc.UserUC),
		teamHandler: team.NewHandler(uc.TeamUC),
		prHandler:   pr.NewHandler(uc.PrUC),
		registrator: NewRegistrator(uc),
	}
}

func (h *Handler) Register(r *mux.Router) {
	h.registrator.RegisterAll(r, h.registrator.uc)
}
