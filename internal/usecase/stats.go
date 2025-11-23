package usecase

import (
	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/repository"
)

type StatsUseCase interface {
	GetStats() (*models.StatsResponse, error)
}

type statsUseCase struct {
	statsRepo repository.StatsRepository
}

func NewStatsUseCase(statsRepo repository.StatsRepository) StatsUseCase {
	return &statsUseCase{
		statsRepo: statsRepo,
	}
}

func (uc *statsUseCase) GetStats() (*models.StatsResponse, error) {
	return uc.statsRepo.GetStats()
}
