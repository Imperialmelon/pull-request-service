package models

type UserDB struct {
	ID       int    `json:"id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type UserAPI struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserWithPRsResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestAPIShort `json:"pull_requests"`
}
