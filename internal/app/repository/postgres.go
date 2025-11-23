package repository

import (
	"database/sql"
	"fmt"

	"github.com/Imperialmelon/AvitoTechTest/internal/app/repository/pr"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/repository/team"
	"github.com/Imperialmelon/AvitoTechTest/internal/app/repository/user"
	"github.com/Imperialmelon/AvitoTechTest/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB       *sql.DB
	TeamRepo *team.PostgresRepository
	UserRepo *user.PostgresRepository
	PrRepo   *pr.PostgresRepository
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &PostgresStore{
		DB: db,
	}

	store.UserRepo = user.NewPostgresRepository(db)
	store.TeamRepo = team.NewPostgresRepository(db)
	store.PrRepo = pr.NewPostgresRepository(db)

	return store, nil
}

func (s *PostgresStore) Close() error {
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

func (s *PostgresStore) SetIsActive(userID string, isActive bool) (models.UserAPI, error) {
	return s.UserRepo.SetIsActive(userID, isActive)
}

func (s *PostgresStore) Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error) {
	return s.PrRepo.Create(req)
}

func (s *PostgresStore) Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error) {
	return s.PrRepo.Merge(prID)
}

func (s *PostgresStore) Reassign(prID string, oldRevID string) (models.PullRequestAPIShortWithReviewersReassigned, error) {
	return s.PrRepo.Reassign(prID, oldRevID)
}

func (s *PostgresStore) GetPRsByUserIDToReview(userID string) ([]models.PullRequestAPIShort, error) {
	return s.PrRepo.GetPRsByUserIDToReview(userID)
}

func (s *PostgresStore) Add(req models.TeamApi) (models.TeamApi, error) {
	return s.TeamRepo.Add(req)
}

func (s *PostgresStore) Get(teamName string) (models.TeamApi, error) {
	return s.TeamRepo.Get(teamName)
}
