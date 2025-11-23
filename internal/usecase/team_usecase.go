package usecase

import (
	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/repository"
)

type teamUseCase struct {
	teamRepo repository.TeamRepository
}

func NewTeamUseCase(teamRepo repository.TeamRepository) TeamUseCase {
	return &teamUseCase{
		teamRepo: teamRepo,
	}
}

func (uc *teamUseCase) CreateTeam(team *models.Team) (*models.Team, error) {
	exists, err := uc.teamRepo.TeamExists(team.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, models.NewErrorResponse(models.TeamExists, "team_name already exists")
	}

	err = uc.teamRepo.CreateTeam(team)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (uc *teamUseCase) GetTeam(teamName string) (*models.Team, error) {
	team, err := uc.teamRepo.GetTeam(teamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, models.NewErrorResponse(models.NotFound, "team not found")
	}

	return team, nil
}
