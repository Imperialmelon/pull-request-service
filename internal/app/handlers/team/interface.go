package team

import "github.com/imperialmelon/avito/internal/models"

type TeamUseCase interface {
	Add(req models.TeamApi) (models.TeamApi, error)
	Get(teamName string) (models.TeamApi, error)
}
