package handlers

import (
	"encoding/json"
	"net/http"

	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userUsecase.SetUserActive(req.UserID, req.IsActive)
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

	WriteJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	response, err := h.userUsecase.GetUserReviewPRs(userID)
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

	WriteJSON(w, http.StatusOK, response)
}
