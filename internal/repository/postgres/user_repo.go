package postgres

import (
	"database/sql"

	"AvitoPullRequest/internal/models"
)

type UserRepo struct {
	pg *Postgres
}

func NewUserRepo(pg *Postgres) *UserRepo {
	return &UserRepo{pg: pg}
}

func (r *UserRepo) CreateUser(user *models.User) error {
	query := `INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4)`
	_, err := r.pg.DB.Exec(query, user.UserID, user.Username, user.TeamName, user.IsActive)
	return err
}

func (r *UserRepo) GetUser(userID string) (*models.User, error) {
	query := `SELECT user_id, username, team_name, is_active, created_at FROM users WHERE user_id = $1`
	row := r.pg.DB.QueryRow(query, userID)

	var user models.User
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(user *models.User) error {
	query := `UPDATE users SET username = $1, team_name = $2, is_active = $3 WHERE user_id = $4`
	_, err := r.pg.DB.Exec(query, user.Username, user.TeamName, user.IsActive, user.UserID)
	return err
}

func (r *UserRepo) GetActiveUsersByTeam(teamName string) ([]*models.User, error) {
	query := `SELECT user_id, username, team_name, is_active, created_at 
	          FROM users WHERE team_name = $1 AND is_active = true`
	rows, err := r.pg.DB.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepo) SetUserActive(userID string, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE user_id = $2`
	_, err := r.pg.DB.Exec(query, isActive, userID)
	return err
}
