package usecase

import (
	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/repository"
)

type userUseCase struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

func NewUserUseCase(userRepo repository.UserRepository, prRepo repository.PullRequestRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (uc *userUseCase) SetUserActive(userID string, isActive bool) (*models.User, error) {
	user, err := uc.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, models.NewErrorResponse(models.NotFound, "user not found")
	}

	err = uc.userRepo.SetUserActive(userID, isActive)
	if err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

func (uc *userUseCase) GetUserReviewPRs(userID string) (*UserReviewResponse, error) {
	user, err := uc.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, models.NewErrorResponse(models.NotFound, "user not found")
	}

	prs, err := uc.prRepo.GetPRsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	var shortPRs []*models.PullRequestShort
	for _, pr := range prs {
		shortPRs = append(shortPRs, &models.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return &UserReviewResponse{
		UserID:       userID,
		PullRequests: shortPRs,
	}, nil
}
