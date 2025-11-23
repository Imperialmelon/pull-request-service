package service

import (
	"github.com/Imperialmelon/AvitoTechTest/internal/app/usecase/pr"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/usecase/team"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/usecase/user"
	"github.com/Imperialmelon/AvitoTechTest/internal/models"

	prsvc "github.com/Imperialmelon/AvitoTechTest/internal/app/service/pr"
	teamsvc "github.com/Imperialmelon/AvitoTechTest/internal/app/service/team"
	usersvc "github.com/Imperialmelon/AvitoTechTest/internal/app/service/user"
)

type Repository interface {
	Close() error
	PRRepository
	UserRepository
	TeamRepository
}

type PRRepository interface {
	Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error)
	Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error)
	Reassign(prID string, oldRevID string) (models.PullRequestAPIShortWithReviewersReassigned, error)
	GetPRsByUserIDToReview(userID string) ([]models.PullRequestAPIShort, error)
}

type UserRepository interface {
	SetIsActive(userID string, isActive bool) (models.UserAPI, error)
}

type TeamRepository interface {
	Add(req models.TeamApi) (models.TeamApi, error)
	Get(teamName string) (models.TeamApi, error)
}

type Service struct {
	Usersvc user.UserService
	Teamsvc team.TeamService
	Prsvc   pr.PRService
}

func NewService(store Repository) *Service {

	userService := usersvc.NewService(store, store)
	teamService := teamsvc.NewService(store)
	prService := prsvc.NewService(store)

	return &Service{
		Usersvc: userService,
		Teamsvc: teamService,
		Prsvc:   prService,
	}
}
