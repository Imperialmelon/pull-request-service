package team

import (
	pkgerrors "github.com/pkg/errors"

	"github.com/imperialmelon/avito/internal/models"
)

type Service struct {
	repo TeamRepository
}

func NewService(repo TeamRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(name string) (models.TeamApi, error) {
	team, err := s.repo.Get(name)
	if err != nil {
		return models.TeamApi{}, pkgerrors.Wrap(err, "team.GetTeam: failed to get team")
	}

	return team, nil
}

func (s *Service) Add(req models.TeamApi) (models.TeamApi, error) {
	createdTeam, err := s.repo.Add(req)
	if err != nil {
		return models.TeamApi{}, pkgerrors.Wrap(err, "team.CreateTeam: failed to create team")
	}

	return createdTeam, nil
}
