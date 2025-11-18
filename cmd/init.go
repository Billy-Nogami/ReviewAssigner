// cmd/init.go
package main

import (
	"log"
	"os"

	"ReviewAssigner/internal/delivery/http"
	"ReviewAssigner/internal/repository/postgres"
	"ReviewAssigner/internal/usecase/pr"
	"ReviewAssigner/internal/usecase/team"
	"ReviewAssigner/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitApp() *gin.Engine {
	// === Подключение к БД ===
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "review_assigner")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASS", "pass")

	dsn := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" password=" + dbPass +
		" dbname=" + dbName +
		" sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// === Репозитории и UseCase ===
	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	prRepo := postgres.NewPullRequestRepository(db)

	teamUsecase := team.NewUsecase(teamRepo)
	userUsecase := user.NewUsecase(userRepo, prRepo)
	prUsecase := pr.NewUsecase(userRepo, prRepo)

	handlers := http.NewHandlers(teamUsecase, userUsecase, prUsecase)

	// === Gin ===
	r := gin.Default()

	// Публичные роуты
	r.GET("/health", handlers.Health)
	r.POST("/auth/login", handlers.Login)

	// Все остальные роуты — защищённые
	handlers.RegisterRoutes(r)

	log.Println("Server initialized successfully")
	return r
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
