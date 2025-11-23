package pr

import (
	"log"

	pkgerrors "github.com/pkg/errors"

	"github.com/Imperialmelon/AvitoTechTest/internal/models"
)

type UseCase struct {
	prSvc PRService
}

func NewUseCase(prService PRService) *UseCase {
	return &UseCase{
		prSvc: prService,
	}
}

func (uc *UseCase) Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error) {
	pr, err := uc.prSvc.Create(req)
	if err != nil {
		log.Println("Failed to create PR", "error", err, "pr_title", req.PullRequestName)
		return models.PullRequestAPIShortWithReviewers{}, pkgerrors.Wrap(err, "pr.Create")
	}
	return pr, nil
}

func (uc *UseCase) Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error) {
	pr, err := uc.prSvc.Merge(prID)
	if err != nil {
		log.Println("Failed to merge PR", "error", err, "pr_id", prID)
		return models.PullRequestAPIShortWithReviewersMerged{}, pkgerrors.Wrap(err, "pr.Merge")
	}
	return pr, nil
}

func (uc *UseCase) Reassign(prID, userID string) (models.PullRequestAPIShortWithReviewersReassigned, error) {
	pr, err := uc.prSvc.Reassign(prID, userID)
	if err != nil {
		log.Println("Failed to reassign PR", "error", err, "pr_id", prID, "new_user_id", userID)
		return models.PullRequestAPIShortWithReviewersReassigned{}, pkgerrors.Wrap(err, "pr.Reassign")
	}
	return pr, nil
}
