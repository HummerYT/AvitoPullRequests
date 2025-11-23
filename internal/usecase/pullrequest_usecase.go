package usecase

import (
	"math/rand"
	"time"

	"AvitoPullRequest/internal/models"
	"AvitoPullRequest/internal/repository"
)

type pullRequestUseCase struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPullRequestUseCase(prRepo repository.PullRequestRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) PullRequestUseCase {
	return &pullRequestUseCase{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (uc *pullRequestUseCase) CreatePR(prID, prName, authorID string) (*models.PullRequest, error) {
	author, err := uc.userRepo.GetUser(authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, models.NewErrorResponse(models.NotFound, "author not found")
	}

	existingPR, err := uc.prRepo.GetPR(prID)
	if err != nil {
		return nil, err
	}
	if existingPR != nil {
		return nil, models.NewErrorResponse(models.PRExists, "PR id already exists")
	}

	reviewers := uc.selectReviewers(author.TeamName, authorID)

	pr := &models.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
	}

	err = uc.prRepo.CreatePR(pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (uc *pullRequestUseCase) MergePR(prID string) (*models.PullRequest, error) {
	pr, err := uc.prRepo.GetPR(prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, models.NewErrorResponse(models.NotFound, "PR not found")
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	now := time.Now()
	pr.Status = "MERGED"
	pr.MergedAt = &now

	err = uc.prRepo.UpdatePR(pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (uc *pullRequestUseCase) ReassignReviewer(prID, oldUserID string) (*models.ReassignResponse, error) {
	pr, err := uc.prRepo.GetPR(prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, models.NewErrorResponse(models.NotFound, "PR not found")
	}

	if pr.Status == "MERGED" {
		return nil, models.NewErrorResponse(models.PRMerged, "cannot reassign on merged PR")
	}

	found := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, models.NewErrorResponse(models.NotAssigned, "reviewer is not assigned")
	}

	newReviewer, err := uc.findReplacementReviewer(oldUserID, pr.AssignedReviewers)
	if err != nil {
		return nil, err
	}

	newReviewers := make([]string, len(pr.AssignedReviewers))
	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			newReviewers[i] = newReviewer
		} else {
			newReviewers[i] = reviewer
		}
	}

	pr.AssignedReviewers = newReviewers
	err = uc.prRepo.UpdatePR(pr)
	if err != nil {
		return nil, err
	}

	return &models.ReassignResponse{
		PR:         pr,
		ReplacedBy: newReviewer,
	}, nil
}

func (uc *pullRequestUseCase) selectReviewers(teamName, excludeUserID string) []string {
	users, err := uc.userRepo.GetActiveUsersByTeam(teamName)
	if err != nil {
		return []string{}
	}

	var candidates []string
	for _, user := range users {
		if user.UserID != excludeUserID && user.IsActive {
			candidates = append(candidates, user.UserID)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if len(candidates) > 2 {
		return candidates[:2]
	}
	return candidates
}

func (uc *pullRequestUseCase) findReplacementReviewer(oldUserID string, excludeUsers []string) (string, error) {
	oldUser, err := uc.userRepo.GetUser(oldUserID)
	if err != nil {
		return "", err
	}
	if oldUser == nil {
		return "", models.NewErrorResponse(models.NotFound, "user not found")
	}

	users, err := uc.userRepo.GetActiveUsersByTeam(oldUser.TeamName)
	if err != nil {
		return "", err
	}

	var candidates []string
	for _, user := range users {
		if user.IsActive && !contains(excludeUsers, user.UserID) && user.UserID != oldUserID {
			candidates = append(candidates, user.UserID)
		}
	}

	if len(candidates) == 0 {
		return "", models.NewErrorResponse(models.NoCandidate, "no active replacement candidates")
	}

	rand.Seed(time.Now().UnixNano())
	return candidates[rand.Intn(len(candidates))], nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
