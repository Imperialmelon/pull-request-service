package team

import "github.com/Imperialmelon/AvitoTechTest/internal/models"

type TeamUseCase interface {
	Add(req models.TeamApi) (models.TeamApi, error)
	Get(teamName string) (models.TeamApi, error)
}
