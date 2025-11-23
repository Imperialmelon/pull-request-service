package pr

import "github.com/Imperialmelon/AvitoTechTest/internal/models"

type PRUseCase interface {
	Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error)
	Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error)
	Reassign(prID string, oldRevID string) (models.PullRequestAPIShortWithReviewersReassigned, error)
}
