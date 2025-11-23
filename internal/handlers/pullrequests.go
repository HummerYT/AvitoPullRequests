package handlers

import (
	"encoding/json"
	"net/http"

	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/usecase"
)

type PullRequestHandler struct {
	prUsecase usecase.PullRequestUseCase
}

func NewPullRequestHandler(prUsecase usecase.PullRequestUseCase) *PullRequestHandler {
	return &PullRequestHandler{
		prUsecase: prUsecase,
	}
}

func (h *PullRequestHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pr, err := h.prUsecase.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if appErr, ok := err.(*models.AppError); ok {
			errorResp := models.ToErrorResponse(appErr)
			switch appErr.Code {
			case models.NotFound:
				WriteJSON(w, http.StatusNotFound, errorResp)
			case models.PRExists:
				WriteJSON(w, http.StatusConflict, errorResp)
			default:
				WriteJSON(w, http.StatusBadRequest, errorResp)
			}
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]interface{}{"pr": pr})
}

func (h *PullRequestHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req models.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pr, err := h.prUsecase.MergePR(req.PullRequestID)
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

	WriteJSON(w, http.StatusOK, map[string]interface{}{"pr": pr})
}

func (h *PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req models.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.prUsecase.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		if appErr, ok := err.(*models.AppError); ok {
			errorResp := models.ToErrorResponse(appErr)
			switch appErr.Code {
			case models.NotFound:
				WriteJSON(w, http.StatusNotFound, errorResp)
			case models.PRMerged, models.NotAssigned, models.NoCandidate:
				WriteJSON(w, http.StatusConflict, errorResp)
			default:
				WriteJSON(w, http.StatusBadRequest, errorResp)
			}
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, response)
}
