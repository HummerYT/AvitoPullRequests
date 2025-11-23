package postgres

import (
	"database/sql"
	"encoding/json"

	"AvitoPullRequest/internal/models"
)

type PullRequestRepo struct {
	pg *Postgres
}

func NewPullRequestRepo(pg *Postgres) *PullRequestRepo {
	return &PullRequestRepo{pg: pg}
}

func (r *PullRequestRepo) CreatePR(pr *models.PullRequest) error {
	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return err
	}

	query := `INSERT INTO pull_requests 
	          (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.pg.DB.Exec(query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID,
		pr.Status, reviewersJSON, pr.CreatedAt)
	return err
}

func (r *PullRequestRepo) GetPR(prID string) (*models.PullRequest, error) {
	query := `SELECT pull_request_id, pull_request_name, author_id, status, 
	                 assigned_reviewers, created_at, merged_at 
	          FROM pull_requests WHERE pull_request_id = $1`
	row := r.pg.DB.QueryRow(query, prID)

	var pr models.PullRequest
	var reviewersJSON []byte
	var mergedAt sql.NullTime

	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status,
		&reviewersJSON, &pr.CreatedAt, &mergedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	err = json.Unmarshal(reviewersJSON, &pr.AssignedReviewers)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepo) UpdatePR(pr *models.PullRequest) error {
	reviewersJSON, err := json.Marshal(pr.AssignedReviewers)
	if err != nil {
		return err
	}

	query := `UPDATE pull_requests 
	          SET pull_request_name = $1, status = $2, assigned_reviewers = $3, merged_at = $4 
	          WHERE pull_request_id = $5`
	_, err = r.pg.DB.Exec(query, pr.PullRequestName, pr.Status, reviewersJSON, pr.MergedAt, pr.PullRequestID)
	return err
}

func (r *PullRequestRepo) GetPRsByReviewer(userID string) ([]*models.PullRequest, error) {
	query := `SELECT pull_request_id, pull_request_name, author_id, status, 
	                 assigned_reviewers, created_at, merged_at 
	          FROM pull_requests 
	          WHERE assigned_reviewers @> $1`

	userIDJSON, err := json.Marshal([]string{userID})
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.DB.Query(query, userIDJSON)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*models.PullRequest
	for rows.Next() {
		var pr models.PullRequest
		var reviewersJSON []byte
		var mergedAt sql.NullTime

		err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status,
			&reviewersJSON, &pr.CreatedAt, &mergedAt)
		if err != nil {
			return nil, err
		}

		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}

		err = json.Unmarshal(reviewersJSON, &pr.AssignedReviewers)
		if err != nil {
			return nil, err
		}

		prs = append(prs, &pr)
	}

	return prs, nil
}
