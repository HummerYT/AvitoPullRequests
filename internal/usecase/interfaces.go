package usecase

import "AvitoPullRequest/internal/models"

type UserUseCase interface {
	SetUserActive(userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(userID string) (*UserReviewResponse, error)
}

type TeamUseCase interface {
	CreateTeam(team *models.Team) (*models.Team, error)
	GetTeam(teamName string) (*models.Team, error)
}

type PullRequestUseCase interface {
	CreatePR(prID, prName, authorID string) (*models.PullRequest, error)
	MergePR(prID string) (*models.PullRequest, error)
	ReassignReviewer(prID, oldUserID string) (*models.ReassignResponse, error)
}

type StatsRepository interface {
	GetStats() (*models.StatsResponse, error)
}

type UserReviewResponse struct {
	UserID       string                     `json:"user_id"`
	PullRequests []*models.PullRequestShort `json:"pull_requests"`
}
