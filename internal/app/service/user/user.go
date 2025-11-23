package user

import (
	pkgerrors "github.com/pkg/errors"

	"github.com/Imperialmelon/AvitoTechTest/internal/models"
)

type Service struct {
	userRepo UserRepository
	prRepo   PRRepository
}

func NewService(userRepo UserRepository, prRepo PRRepository) *Service {
	return &Service{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *Service) SetIsActive(req models.SetActiveRequest) (models.UserAPI, error) {
	user, err := s.userRepo.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		return models.UserAPI{}, pkgerrors.Wrap(err, "user.SetIsActive: failed to update user status")
	}
	return user, nil
}

func (s *Service) GetReview(userID string) (models.UserWithPRsResponse, error) {
	prs, err := s.prRepo.GetPRsByUserIDToReview(userID)
	if err != nil {
		return models.UserWithPRsResponse{}, pkgerrors.Wrap(err, "user.GetReview: failed to get PRs for user")
	}
	return models.UserWithPRsResponse{
		UserID:       userID,
		PullRequests: prs,
	}, nil
}
