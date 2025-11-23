package postgres

import (
	"log"

	"AvitoPullRequest/internal/models"
)

type StatsRepo struct {
	pg *Postgres
}

func NewStatsRepo(pg *Postgres) *StatsRepo {
	return &StatsRepo{pg: pg}
}
func (r *StatsRepo) GetStats() (*models.StatsResponse, error) {
	stats := &models.StatsResponse{}

	err := r.pg.DB.QueryRow("SELECT COUNT(*) FROM pull_requests").Scan(&stats.TotalPRs)
	if err != nil {
		log.Printf("Error getting total PRs: %v", err)
		return nil, err
	}

	err = r.pg.DB.QueryRow("SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN'").Scan(&stats.OpenPRs)
	if err != nil {
		log.Printf("Error getting open PRs: %v", err)
		return nil, err
	}

	err = r.pg.DB.QueryRow("SELECT COUNT(*) FROM pull_requests WHERE status = 'MERGED'").Scan(&stats.MergedPRs)
	if err != nil {
		log.Printf("Error getting merged PRs: %v", err)
		return nil, err
	}

	rows, err := r.pg.DB.Query(`
        SELECT u.user_id, u.username, 
               COUNT(pr.pull_request_id) as assignment_count
        FROM users u
        LEFT JOIN pull_requests pr ON u.user_id::text = pr.author_id::text
        WHERE u.is_active = true
        GROUP BY u.user_id, u.username
        ORDER BY assignment_count DESC
    `)
	if err != nil {
		log.Printf("Error querying user stats: %v", err)
		return nil, err
	}
	defer rows.Close()

	stats.UserAssignments = []models.UserStats{}
	for rows.Next() {
		var userStat models.UserStats
		var userID string // меняем на string!
		err := rows.Scan(&userID, &userStat.Username, &userStat.Assignments)
		if err != nil {
			log.Printf("Error scanning user stats: %v", err)
			continue
		}
		userStat.UserID = userID
		stats.UserAssignments = append(stats.UserAssignments, userStat)
	}

	teamRows, err := r.pg.DB.Query(`
        SELECT t.team_name, COUNT(DISTINCT pr.pull_request_id) as pr_count
        FROM teams t
        LEFT JOIN users u ON t.team_name = u.team_name  
        LEFT JOIN pull_requests pr ON u.user_id::text = pr.author_id::text
        GROUP BY t.team_name
        ORDER BY pr_count DESC
    `)
	if err != nil {
		log.Printf("Error querying team stats: %v", err)
		return nil, err
	}
	defer teamRows.Close()

	stats.TeamStats = []models.TeamStat{}
	for teamRows.Next() {
		var teamStat models.TeamStat
		err := teamRows.Scan(&teamStat.TeamName, &teamStat.PRCount)
		if err != nil {
			log.Printf("Error scanning team stats: %v", err)
			continue
		}
		stats.TeamStats = append(stats.TeamStats, teamStat)
	}

	return stats, nil
}
