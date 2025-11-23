package models

type TeamDB struct {
	ID     int    `json:"id"`
	TeamID string `json:"team_id"`
	Name   string `json:"name"`
}

type TeamMemberApi struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamApi struct {
	Name    string          `json:"team_name"`
	Members []TeamMemberApi `json:"members"`
}

type GetTeamRequest struct {
	Name string `json:"team_name"`
}
