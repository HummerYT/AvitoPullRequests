package models

type StatsResponse struct {
	TotalPRs        int         `json:"total_prs"`
	OpenPRs         int         `json:"open_prs"`
	MergedPRs       int         `json:"merged_prs"`
	UserAssignments []UserStats `json:"user_assignments"`
	TeamStats       []TeamStat  `json:"team_stats"`
}

type UserStats struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Assignments int    `json:"assignments"`
}

type TeamStat struct {
	TeamName string `json:"team_name"`
	PRCount  int    `json:"pr_count"`
}
