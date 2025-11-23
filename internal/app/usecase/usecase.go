package usecase

import (
	"github.com/imperialmelon/avito/internal/app/service"
	"github.com/imperialmelon/avito/internal/app/usecase/pr"
	"github.com/imperialmelon/avito/internal/app/usecase/team"
	"github.com/imperialmelon/avito/internal/app/usecase/user"
)

// Repository алиас для service.Repository
type Repository = service.Repository

type UseCase struct {
	service *service.Service
	UserUC  *user.UseCase
	PrUC    *pr.UseCase
	TeamUC  *team.UseCase
}

func NewUseCase(service *service.Service) *UseCase {
	prUC := pr.NewUseCase(service.Prsvc)
	userUC := user.NewUseCase(service.Usersvc)
	teamUC := team.NewUseCase(service.Teamsvc)

	return &UseCase{
		service: service,
		UserUC:  userUC,
		PrUC:    prUC,
		TeamUC:  teamUC,
	}
}
