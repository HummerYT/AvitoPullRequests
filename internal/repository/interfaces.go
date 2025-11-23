package repository

import (
	"AvitoPullRequest/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUser(userID string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetActiveUsersByTeam(teamName string) ([]*models.User, error)
	SetUserActive(userID string, isActive bool) error
}

type TeamRepository interface {
	CreateTeam(team *models.Team) error
	GetTeam(teamName string) (*models.Team, error)
	TeamExists(teamName string) (bool, error)
}

type PullRequestRepository interface {
	CreatePR(pr *models.PullRequest) error
	GetPR(prID string) (*models.PullRequest, error)
	UpdatePR(pr *models.PullRequest) error
	GetPRsByReviewer(userID string) ([]*models.PullRequest, error)
}

type StatsRepository interface {
	GetStats() (*models.StatsResponse, error)
}
