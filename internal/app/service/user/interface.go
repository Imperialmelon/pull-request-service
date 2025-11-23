package user

import (
	"github.com/imperialmelon/avito/internal/models"
)

type UserRepository interface {
	SetIsActive(userID string, isActive bool) (models.UserAPI, error)
}

type PRRepository interface {
	Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error)
	Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error)
	Reassign(prID string, oldRevID string) (models.PullRequestAPIShortWithReviewersReassigned, error)
	GetPRsByUserIDToReview(userID string) ([]models.PullRequestAPIShort, error)
}
