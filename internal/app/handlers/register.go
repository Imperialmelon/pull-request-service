package handlers

import (
	"github.com/gorilla/mux"

	"github.com/imperialmelon/avito/internal/app/handlers/pr"
	"github.com/imperialmelon/avito/internal/app/handlers/team"
	"github.com/imperialmelon/avito/internal/app/handlers/user"
	"github.com/imperialmelon/avito/internal/app/usecase"
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
