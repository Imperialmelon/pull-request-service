package pr

import (
	"github.com/Imperialmelon/AvitoTechTest/internal/models"
)

type PRService interface {
	Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error)
	Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error)
	Reassign(prID, userID string) (models.PullRequestAPIShortWithReviewersReassigned, error)
}
