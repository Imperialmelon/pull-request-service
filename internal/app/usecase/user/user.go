package user

import (
	"log"

	"github.com/imperialmelon/avito/internal/models"
	pkgerrors "github.com/pkg/errors"
)

type UseCase struct {
	userSvc UserService
}

func NewUseCase(userService UserService) *UseCase {
	return &UseCase{
		userSvc: userService,
	}
}

func (uc *UseCase) SetIsActive(req models.SetActiveRequest) (models.UserAPI, error) {
	updatedUser, err := uc.userSvc.SetIsActive(req)
	if err != nil {
		log.Println("Failed to set user active status", "error", err, "user_id", req.UserID)
		return models.UserAPI{}, pkgerrors.Wrap(err, "user.SetIsActive")
	}
	return updatedUser, nil
}

func (uc *UseCase) GetReview(userID string) (models.UserWithPRsResponse, error) {
	userWithPRs, err := uc.userSvc.GetReview(userID)
	if err != nil {
		log.Println("Failed to get user review", "error", err, "user_id", userID)
		return models.UserWithPRsResponse{}, pkgerrors.Wrap(err, "user.GetReview")
	}
	return userWithPRs, nil
}
