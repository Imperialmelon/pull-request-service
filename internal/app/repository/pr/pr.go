package pr

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	repoerrors "github.com/Imperialmelon/AvitoTechTest/internal/app/repository/errors"
	svcerrors "github.com/Imperialmelon/AvitoTechTest/internal/errors"
	"github.com/Imperialmelon/AvitoTechTest/internal/models"
	"github.com/lib/pq"
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

func (r *PostgresRepository) Create(req models.CreatePRRequest) (models.PullRequestAPIShortWithReviewers, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var authorTMID int
	var authorUserID string
	err = tx.QueryRow(
		`SELECT tm._id, u.user_id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE u.user_id = $1`,
		req.AuthorID,
	).Scan(&authorTMID, &authorUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequestAPIShortWithReviewers{}, svcerrors.ErrorNotFound
		}
		return models.PullRequestAPIShortWithReviewers{}, err
	}

	var prInternalID int
	_, err = tx.Exec(
		`INSERT INTO pull_request (pr_id, req_title, author_id, status)
         VALUES ($1, $2, $3, 'open')`,
		req.PullRequestID, req.PullRequestName, authorTMID,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == repoerrors.UniqueViolation {
				return models.PullRequestAPIShortWithReviewers{}, svcerrors.ErrorPRExists
			}
		}
		return models.PullRequestAPIShortWithReviewers{}, err
	}

	err = tx.QueryRow(
		`SELECT _id FROM pull_request WHERE pr_id = $1`,
		req.PullRequestID,
	).Scan(&prInternalID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}

	var teamID int
	err = tx.QueryRow(
		`SELECT team_id FROM team_member WHERE _id = $1`,
		authorTMID,
	).Scan(&teamID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}

	rows, err := tx.Query(
		`SELECT tm._id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE tm.team_id = $1
           AND tm._id != $2
           AND u.is_active = true
         ORDER BY random()
         LIMIT 2`,
		teamID, authorTMID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}
	defer rows.Close()

	var candidateTMIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return models.PullRequestAPIShortWithReviewers{}, err
		}
		candidateTMIDs = append(candidateTMIDs, id)
	}
	if err := rows.Err(); err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}

	now := time.Now().UTC()
	for _, tmid := range candidateTMIDs {
		_, err := tx.Exec(
			`INSERT INTO reviewer (req_id, user_id, created_at)
             VALUES ($1, $2, $3)`,
			prInternalID, tmid, now,
		)
		if err != nil {
			return models.PullRequestAPIShortWithReviewers{}, err
		}
	}

	rows, err = tx.Query(
		`SELECT u.user_id
         FROM reviewer r
         JOIN team_member tm ON tm._id = r.user_id
         JOIN "user" u ON u._id = tm.user_id
         WHERE r.req_id = $1`,
		prInternalID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}
	defer rows.Close()

	var assigned []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return models.PullRequestAPIShortWithReviewers{}, err
		}
		assigned = append(assigned, uid)
	}
	if err := rows.Err(); err != nil {
		return models.PullRequestAPIShortWithReviewers{}, err
	}

	out := models.PullRequestAPIShortWithReviewers{
		PullRequestAPIShort: models.PullRequestAPIShort{
			PullRequestID:   req.PullRequestID,
			PullRequestName: req.PullRequestName,
			AuthorID:        authorUserID,
			Status:          models.PrStatus("OPEN"),
		},
		Reviewers: assigned,
	}

	return out, nil
}

func (r *PostgresRepository) Merge(prID string) (models.PullRequestAPIShortWithReviewersMerged, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.PullRequestAPIShortWithReviewersMerged{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var prInternalID int
	var title string
	var status string
	var authorTMID int
	err = tx.QueryRow(
		`SELECT _id, req_title, status, author_id FROM pull_request WHERE pr_id = $1`,
		prID,
	).Scan(&prInternalID, &title, &status, &authorTMID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequestAPIShortWithReviewersMerged{}, svcerrors.ErrorNotFound
		}
		return models.PullRequestAPIShortWithReviewersMerged{}, err
	}

	if status == string(models.MERGED) {
		var authorUserID string
		_ = tx.QueryRow(
			`SELECT u.user_id
             FROM team_member tm
             JOIN "user" u ON u._id = tm.user_id
             WHERE tm._id = $1`,
			authorTMID,
		).Scan(&authorUserID)

		rows, err := tx.Query(
			`SELECT u.user_id, r.merged_at
             FROM reviewer r
             JOIN team_member tm ON tm._id = r.user_id
             JOIN "user" u ON u._id = tm.user_id
             WHERE r.req_id = $1`,
			prInternalID,
		)
		if err != nil {
			return models.PullRequestAPIShortWithReviewersMerged{}, err
		}
		defer rows.Close()

		var reviewers []string
		var mergedAt time.Time
		for rows.Next() {
			var rv string
			var ma sql.NullTime
			if err := rows.Scan(&rv, &ma); err != nil {
				return models.PullRequestAPIShortWithReviewersMerged{}, err
			}
			reviewers = append(reviewers, rv)
			if ma.Valid {
				mergedAt = ma.Time
			}
		}
		out := models.PullRequestAPIShortWithReviewersMerged{
			PullRequestAPIShortWithReviewers: models.PullRequestAPIShortWithReviewers{
				PullRequestAPIShort: models.PullRequestAPIShort{
					PullRequestID:   prID,
					PullRequestName: title,
					AuthorID:        authorUserID,
					Status:          models.PrStatus("MERGED"),
				},
				Reviewers: reviewers,
			},
			MergedAt: mergedAt,
		}
		return out, nil
	}
	_, err = tx.Exec(
		`UPDATE pull_request SET status = 'merged' WHERE _id = $1`,
		prInternalID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersMerged{}, err
	}

	now := time.Now().UTC()
	_, err = tx.Exec(
		`UPDATE reviewer SET merged_at = $1 WHERE req_id = $2`,
		now, prInternalID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersMerged{}, err
	}

	var authorUserID string
	err = tx.QueryRow(
		`SELECT u.user_id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE tm._id = $1`,
		authorTMID,
	).Scan(&authorUserID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersMerged{}, err
	}

	rows, err := tx.Query(
		`SELECT u.user_id
         FROM reviewer r
         JOIN team_member tm ON tm._id = r.user_id
         JOIN "user" u ON u._id = tm.user_id
         WHERE r.req_id = $1`,
		prInternalID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersMerged{}, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return models.PullRequestAPIShortWithReviewersMerged{}, err
		}
		reviewers = append(reviewers, uid)
	}

	out := models.PullRequestAPIShortWithReviewersMerged{
		PullRequestAPIShortWithReviewers: models.PullRequestAPIShortWithReviewers{
			PullRequestAPIShort: models.PullRequestAPIShort{
				PullRequestID:   prID,
				PullRequestName: title,
				AuthorID:        authorUserID,
				Status:          models.PrStatus("MERGED"),
			},
			Reviewers: reviewers,
		},
		MergedAt: now,
	}

	return out, nil
}

func (r *PostgresRepository) Reassign(prID string, oldRevUserID string) (models.PullRequestAPIShortWithReviewersReassigned, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	log.Println("got1")

	var prInternalID int
	var status string
	var title string
	var authorTMID int
	err = tx.QueryRow(
		`SELECT _id, status, req_title, author_id FROM pull_request WHERE pr_id = $1`,
		prID,
	).Scan(&prInternalID, &status, &title, &authorTMID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequestAPIShortWithReviewersReassigned{}, svcerrors.ErrorNotFound
		}
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	log.Println("got2")

	if status == string(models.MERGED) {
		return models.PullRequestAPIShortWithReviewersReassigned{}, svcerrors.ErrorPRMerged
	}

	var oldRevTMID int
	err = tx.QueryRow(
		`SELECT tm._id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE u.user_id = $1`,
		oldRevUserID,
	).Scan(&oldRevTMID)
	log.Println("oldRevUserID:", oldRevUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequestAPIShortWithReviewersReassigned{}, svcerrors.ErrorNotFound
		}
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	log.Println("got3")
	var reviewerRowID int
	log.Println("prInternalID:", prInternalID, "oldRevTMID:", oldRevTMID)
	err = tx.QueryRow(
		`SELECT _id FROM reviewer WHERE req_id = $1 AND user_id = $2`,
		prInternalID, oldRevTMID,
	).Scan(&reviewerRowID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequestAPIShortWithReviewersReassigned{}, svcerrors.ErrorNotFound
		}
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	log.Println("got4")
	var teamID int
	err = tx.QueryRow(
		`SELECT team_id FROM team_member WHERE _id = $1`,
		oldRevTMID,
	).Scan(&teamID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	rows, err := tx.Query(`SELECT user_id FROM reviewer WHERE req_id = $1`, prInternalID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}
	defer rows.Close()
	var currentReviewerTMIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return models.PullRequestAPIShortWithReviewersReassigned{}, err
		}
		currentReviewerTMIDs = append(currentReviewerTMIDs, id)
	}
	if err := rows.Err(); err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	args := []interface{}{teamID, oldRevTMID}
	q := `
        SELECT tm._id
        FROM team_member tm
        JOIN "user" u ON u._id = tm.user_id
        WHERE tm.team_id = $1
          AND tm._id != $2
          AND u.is_active = true
    `
	if len(currentReviewerTMIDs) > 0 {
		q += " AND tm._id NOT IN ("
		for i, v := range currentReviewerTMIDs {
			args = append(args, v)
			q += fmt.Sprintf("$%d", len(args))
			if i < len(currentReviewerTMIDs)-1 {
				q += ","
			}
		}
		q += ")"
	}
	q += " ORDER BY random() LIMIT 1"

	var newTMID int
	err = tx.QueryRow(q, args...).Scan(&newTMID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequestAPIShortWithReviewersReassigned{}, svcerrors.ErrorNoCandidate
		}
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	var newUserID string
	err = tx.QueryRow(
		`SELECT u.user_id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE tm._id = $1`,
		newTMID,
	).Scan(&newUserID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	now := time.Now().UTC()
	_, err = tx.Exec(
		`UPDATE reviewer
         SET user_id = $1, created_at = $2, merged_at = NULL
         WHERE _id = $3`,
		newTMID, now, reviewerRowID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	var authorUserID string
	err = tx.QueryRow(
		`SELECT u.user_id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE tm._id = $1`,
		authorTMID,
	).Scan(&authorUserID)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}

	rows2, err := tx.Query(
		`SELECT u.user_id
         FROM reviewer r
         JOIN team_member tm ON tm._id = r.user_id
         JOIN "user" u ON u._id = tm.user_id
         WHERE r.req_id = $1`,
		prInternalID,
	)
	if err != nil {
		return models.PullRequestAPIShortWithReviewersReassigned{}, err
	}
	defer rows2.Close()

	var reviewers []string
	for rows2.Next() {
		var uid string
		if err := rows2.Scan(&uid); err != nil {
			return models.PullRequestAPIShortWithReviewersReassigned{}, err
		}
		reviewers = append(reviewers, uid)
	}

	out := models.PullRequestAPIShortWithReviewersReassigned{
		PullRequestAPIShortWithReviewers: models.PullRequestAPIShortWithReviewers{
			PullRequestAPIShort: models.PullRequestAPIShort{
				PullRequestID:   prID,
				PullRequestName: title,
				AuthorID:        authorUserID,
				Status:          models.PrStatus("OPEN"),
			},
			Reviewers: reviewers,
		},
		ReplacedBy: newUserID,
	}

	return out, nil
}

func (r *PostgresRepository) GetPRsByUserIDToReview(userID string) ([]models.PullRequestAPIShort, error) {
	var tmID int
	err := r.db.QueryRow(
		`SELECT tm._id
         FROM team_member tm
         JOIN "user" u ON u._id = tm.user_id
         WHERE u.user_id = $1`,
		userID,
	).Scan(&tmID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := r.db.Query(
		`SELECT pr.pr_id, pr.req_title,
                au.user_id AS author_user_id,
                pr.status
         FROM reviewer r
         JOIN pull_request pr ON pr._id = r.req_id
         JOIN team_member at ON at._id = pr.author_id
         JOIN "user" au ON au._id = at.user_id
         WHERE r.user_id = $1`,
		tmID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PullRequestAPIShort
	for rows.Next() {
		var p models.PullRequestAPIShort
		var statusStr string
		if err := rows.Scan(&p.PullRequestID, &p.PullRequestName, &p.AuthorID, &statusStr); err != nil {
			return nil, err
		}
		p.Status = models.PrStatus(statusStr)
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
