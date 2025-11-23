package models

import "time"

type Team struct {
	TeamName  string       `json:"team_name" db:"team_name"`
	Members   []TeamMember `json:"members"`
	CreatedAt time.Time    `json:"created_at,omitempty" db:"created_at"`
}
