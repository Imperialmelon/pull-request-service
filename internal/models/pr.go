package models

import "time"

type PrStatus string

const (
	OPEN   PrStatus = "OPEN"
	MERGED PrStatus = "MERGED"
)

type PullRequestDB struct {
	ID              int
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	Status          PrStatus `json:"status"`
	CreatedAt       *string  `json:"created_at,omitempty"`
	MergedAt        *string  `json:"merged_at,omitempty"`
}

type PullRequestAPI struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            PrStatus `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type PullRequestAPIShort struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          PrStatus `json:"status"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestAPIShortWithReviewers struct {
	PullRequestAPIShort
	Reviewers []string `json:"assigned_reviewer"`
}

type PRID struct {
	PullRequestID string `json:"pull_request_id"`
}

type PullRequestAPIShortWithReviewersMerged struct {
	PullRequestAPIShortWithReviewers
	MergedAt time.Time `json:"merged_at"`
}

type PullRequestAPIShortWithReviewersReassigned struct {
	PullRequestAPIShortWithReviewers
	ReplacedBy string `json:"replaced_by"`
}

type ReassignRequest struct {
	PRID
	UserID string `json:"old_reviewer_id"`
}
