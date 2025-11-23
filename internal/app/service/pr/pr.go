package pr

import (
	pkgerrors "github.com/pkg/errors"

	"github.com/imperialmelon/avito/internal/models"
)

type Service struct {
	repo PRRepository
}

func NewService(repo PRRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error) {
	createdPR, err := s.repo.Create(req)
	if err != nil {
		return models.PullRequestAPIShortWithReviewers{}, pkgerrors.Wrap(err, "pr.Create: failed to create PR")
	}

	return createdPR, nil
}

func (s *Service) Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error) {
	mergedPR, err := s.repo.Merge(prID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersMerged{}, pkgerrors.Wrap(err, "pr.Merge: failed to merge PR")
	}

	return mergedPR, nil
}

func (s *Service) Reassign(prID, oldReviewerID string) (models.PullRequestAPIShortWithReviewersReassigned, error) {
	reassignedPR, err := s.repo.Reassign(prID, oldReviewerID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, pkgerrors.Wrap(err, "pr.Reassign: failed to reassign reviewer")
	}

	return reassignedPR, nil
}

func (s *Service) GetByUserID(userID string) ([]models.PullRequestAPIShort, error) {
	prs, err := s.repo.GetPRsByUserIDToReview(userID)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "pr.GetByUserID: failed to fetch PRs by user")
	}

	return prs, nil
}
