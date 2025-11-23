package handlers

import (
	"github.com/gorilla/mux"

	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers/pr"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers/team"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/handlers/user"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/usecase"
)

type Registrator struct {
	uc *usecase.UseCase
}

func NewRegistrator(uc *usecase.UseCase) *Registrator {
	return &Registrator{
		uc: uc,
	}

}

func (r *Registrator) RegisterAll(router *mux.Router, uc *usecase.UseCase) {
	team.Register(router, uc.TeamUC)
	user.Register(router, uc.UserUC)
	pr.Register(router, uc.PrUC)
}
