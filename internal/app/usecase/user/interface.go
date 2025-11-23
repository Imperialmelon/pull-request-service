package user

import (
	"github.com/imperialmelon/avito/internal/models"
)

type UserService interface {
	SetIsActive(req models.SetActiveRequest) (models.UserAPI, error)
	GetReview(userID string) (models.UserWithPRsResponse, error)
}
