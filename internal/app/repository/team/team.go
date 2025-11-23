package team

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	repoerrors "github.com/Imperialmelon/AvitoTechTest/internal/app/repository/errors"
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

func (r *PostgresRepository) Add(req models.TeamApi) (models.TeamApi, error) {

	var teamID int
	err := r.db.QueryRow(
		`INSERT INTO team (team_id, team_name)
         VALUES ($1, $2)
         ON CONFLICT (team_name) DO UPDATE SET team_name = EXCLUDED.team_name
         RETURNING _id`,
		uuid.New().String(),
		req.Name,
	).Scan(&teamID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == repoerrors.UniqueViolation {
				return models.TeamApi{}, svcerrors.ErrorTeamExists
			}
		}
		return models.TeamApi{}, err
	}

	result := models.TeamApi{Name: req.Name}

	for _, m := range req.Members {

		var userID int
		var member models.TeamMemberApi

		err := r.db.QueryRow(
			`INSERT INTO "user" (user_id, username, is_active)
             VALUES ($1, $2, $3)
             ON CONFLICT (user_id) DO UPDATE SET
                 username = EXCLUDED.username,
                 is_active = EXCLUDED.is_active
             RETURNING _id, user_id, username, is_active`,
			m.UserID, m.Username, m.IsActive,
		).Scan(&userID, &member.UserID, &member.Username, &member.IsActive)

		if err != nil {
			return models.TeamApi{}, err
		}

		_, err = r.db.Exec(
			`INSERT INTO team_member (team_id, user_id)
             VALUES ($1, $2)
             ON CONFLICT DO NOTHING`,
			teamID, userID,
		)
		if err != nil {
			return models.TeamApi{}, err
		}

		result.Members = append(result.Members, member)
	}

	return result, nil
}

func (r *PostgresRepository) Get(teamName string) (models.TeamApi, error) {

	var teamID int
	err := r.db.QueryRow(
		`SELECT _id FROM team WHERE team_name = $1`,
		teamName,
	).Scan(&teamID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TeamApi{}, svcerrors.ErrorNotFound
		}
		return models.TeamApi{}, err
	}

	rows, err := r.db.Query(
		`SELECT 
            u.user_id,
            u.username,
            u.is_active
         FROM team t
         JOIN team_member tm ON tm.team_id = t._id
         JOIN "user" u ON u._id = tm.user_id
         WHERE t._id = $1`,
		teamID,
	)
	if err != nil {
		return models.TeamApi{}, err
	}
	defer rows.Close()

	result := models.TeamApi{Name: teamName}

	for rows.Next() {
		var m models.TeamMemberApi
		err := rows.Scan(&m.UserID, &m.Username, &m.IsActive)
		if err != nil {
			return models.TeamApi{}, err
		}

		result.Members = append(result.Members, m)
	}

	if err = rows.Err(); err != nil {
		return models.TeamApi{}, err
	}

	return result, nil
}
