package postgres

import (
	"database/sql"

	"AvitoPullRequest/internal/models"
)

type TeamRepo struct {
	pg *Postgres
}

func NewTeamRepo(pg *Postgres) *TeamRepo {
	return &TeamRepo{pg: pg}
}

func (r *TeamRepo) CreateTeam(team *models.Team) error {
	tx, err := r.pg.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO teams (team_name) VALUES ($1)", team.TeamName)
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		_, err = tx.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET 
			username = EXCLUDED.username, 
			team_name = EXCLUDED.team_name, 
			is_active = EXCLUDED.is_active
		`, member.UserID, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TeamRepo) GetTeam(teamName string) (*models.Team, error) {
	var team models.Team
	team.TeamName = teamName

	query := `SELECT user_id, username, is_active FROM users WHERE team_name = $1`
	rows, err := r.pg.DB.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	team.Members = members
	return &team, nil
}

func (r *TeamRepo) TeamExists(teamName string) (bool, error) {
	query := `SELECT 1 FROM teams WHERE team_name = $1`
	var exists int
	err := r.pg.DB.QueryRow(query, teamName).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
