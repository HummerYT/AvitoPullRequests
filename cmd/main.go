package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"AvitoPullRequest/internal/handlers"
	"AvitoPullRequest/internal/middleware"
	"AvitoPullRequest/internal/repository/postgres"
	"AvitoPullRequest/internal/usecase"
)

func main() {
	time.Sleep(5 * time.Second)

	pgHost := getEnv("POSTGRES_HOST", "postgres")
	pgPort := 5432
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPassword := getEnv("POSTGRES_PASSWORD", "7549")
	pgDBName := getEnv("POSTGRES_DB", "AvitoPullRequest")
	appHost := getEnv("APP_HOST", "0.0.0.0")
	appPort := "8080"

	log.Printf("Connecting to database: %s@%s:%d/%s", pgUser, pgHost, pgPort, pgDBName)

	pg, err := postgres.New(pgHost, pgPort, pgUser, pgPassword, pgDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pg.Close()

	log.Println("Successfully connected to database")

	userRepo := postgres.NewUserRepo(pg)
	teamRepo := postgres.NewTeamRepo(pg)
	prRepo := postgres.NewPullRequestRepo(pg)
	statsRepo := postgres.NewStatsRepo(pg)

	userUsecase := usecase.NewUserUseCase(userRepo, prRepo)
	teamUsecase := usecase.NewTeamUseCase(teamRepo)
	prUsecase := usecase.NewPullRequestUseCase(prRepo, userRepo, teamRepo)
	statsUsecase := usecase.NewStatsUseCase(statsRepo)

	healthHandler := handlers.NewHealthHandler()
	teamHandler := handlers.NewTeamHandler(teamUsecase)
	userHandler := handlers.NewUserHandler(userUsecase)
	prHandler := handlers.NewPullRequestHandler(prUsecase)
	statsHandler := handlers.NewStatsHandler(statsUsecase)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.HealthCheck)
	mux.HandleFunc("POST /team/add", teamHandler.AddTeam)
	mux.HandleFunc("GET /team/get", teamHandler.GetTeam)
	mux.HandleFunc("POST /users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("GET /users/getReview", userHandler.GetReview)
	mux.HandleFunc("POST /pullrequest/create", prHandler.CreatePR)
	mux.HandleFunc("POST /pullrequest/merge", prHandler.MergePR)
	mux.HandleFunc("POST /pullrequest/reassign", prHandler.Reassign)
	mux.HandleFunc("GET /stats", statsHandler.GetStats)

	handler := middleware.Recovery(middleware.Logging(mux))

	addr := appHost + ":" + appPort
	log.Printf("Server starting on %s", addr)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /health")
	log.Printf("  POST /team/add")
	log.Printf("  GET  /team/get?team_name=name")
	log.Printf("  POST /users/setIsActive")
	log.Printf("  GET  /users/getReview?user_id=id")
	log.Printf("  POST /pullrequest/create")
	log.Printf("  POST /pullrequest/merge")
	log.Printf("  POST /pullrequest/reassign")
	log.Printf("  GET  /stats")
	log.Fatal(http.ListenAndServe(addr, handler))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
