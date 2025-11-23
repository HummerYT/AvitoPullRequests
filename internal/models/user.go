package models

import "time"

type User struct {
	UserID    string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	TeamName  string    `json:"team_name" db:"team_name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
