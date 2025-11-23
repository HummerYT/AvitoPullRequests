package handlers

import (
	"net/http"

	"AvitoPullRequest/internal/usecase"
)

type StatsHandler struct {
	statsUsecase usecase.StatsUseCase
}

func NewStatsHandler(statsUsecase usecase.StatsUseCase) *StatsHandler {
	return &StatsHandler{
		statsUsecase: statsUsecase,
	}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsUsecase.GetStats()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get statistics")
		return
	}

	WriteJSON(w, http.StatusOK, stats)
}
