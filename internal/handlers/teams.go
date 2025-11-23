package handlers

import (
	"encoding/json"
	"net/http"

	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/usecase"
)

type TeamHandler struct {
	teamUsecase usecase.TeamUseCase
}

func NewTeamHandler(teamUsecase usecase.TeamUseCase) *TeamHandler {
	return &TeamHandler{
		teamUsecase: teamUsecase,
	}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdTeam, err := h.teamUsecase.CreateTeam(&team)
	if err != nil {
		if appErr, ok := err.(*models.AppError); ok {
			errorResp := models.ToErrorResponse(appErr)
			WriteJSON(w, http.StatusBadRequest, errorResp)
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]interface{}{"team": createdTeam})
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		WriteError(w, http.StatusBadRequest, "team_name is required")
		return
	}

	team, err := h.teamUsecase.GetTeam(teamName)
	if err != nil {
		if appErr, ok := err.(*models.AppError); ok {
			errorResp := models.ToErrorResponse(appErr)
			if appErr.Code == models.NotFound {
				WriteJSON(w, http.StatusNotFound, errorResp)
				return
			}
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, team)
}
