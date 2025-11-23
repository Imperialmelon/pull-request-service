package team

import (
	"log"

	pkgerrors "github.com/pkg/errors"

	"github.com/Imperialmelon/AvitoTechTest/internal/models"
)

type UseCase struct {
	teamSvc TeamService
}

func NewUseCase(teamService TeamService) *UseCase {
	return &UseCase{
		teamSvc: teamService,
	}
}

func (uc *UseCase) Add(req models.TeamApi) (models.TeamApi, error) {
	team, err := uc.teamSvc.Add(req)
	if err != nil {
		log.Println("Failed to add team", "error", err, "team_name", req.Name)
		return models.TeamApi{}, pkgerrors.Wrap(err, "team.Add")
	}
	return team, nil
}

func (uc *UseCase) Get(teamName string) (models.TeamApi, error) {
	team, err := uc.teamSvc.Get(teamName)
	if err != nil {
		log.Println("Failed to get team", "error", err, "team_name", teamName)
		return models.TeamApi{}, pkgerrors.Wrap(err, "team.Get")
	}
	return team, nil
}
