package user

import (
	"database/sql"
	"errors"
	"fmt"

	svcerrors "github.com/Imperialmelon/AvitoTechTest/internal/errors"
	"github.com/Imperialmelon/AvitoTechTest/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SetIsActive(UserID string, IsActive bool) (models.UserAPI, error) {
	query := `
        WITH updated AS (
            UPDATE "user"
            SET is_active = $1
            WHERE user_id = $2
            RETURNING user_id
        )
        SELECT 
            u.user_id,
            u.username,
            u.is_active,
            t.team_name
        FROM updated up
        JOIN "user" u ON u.user_id = up.user_id
        JOIN team_member tm ON tm.user_id = u._id
        JOIN team t ON t._id = tm.team_id
    `

	var u models.UserAPI

	err := r.db.QueryRow(
		query,
		IsActive,
		UserID,
	).Scan(
		&u.UserID,
		&u.Username,
		&u.IsActive,
		&u.TeamName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserAPI{}, svcerrors.ErrorNotFound
		}
		return models.UserAPI{}, err
	}

	return u, nil
}
