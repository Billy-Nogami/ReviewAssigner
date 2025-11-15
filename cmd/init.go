package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"
	"ReviewAssigner/internal/delivery/http"
	"ReviewAssigner/internal/delivery/middleware"
	"ReviewAssigner/internal/repository/postgres"
	"ReviewAssigner/internal/usecase/pr"
	"ReviewAssigner/internal/usecase/team"
	"ReviewAssigner/internal/usecase/user"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "ReviewAssigner/docs"
)

// @title ReviewAssigner
// @version 1.0
// @description Service for assigning reviewers to Pull Requests
func InitApp() *gin.Engine {
	// Конфиг из env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// Подключение к БД
	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Репозитории
	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	prRepo := postgres.NewPullRequestRepository(db)

	// Usecase
	teamUsecase := team.NewUsecase(teamRepo)
	userUsecase := user.NewUsecase(userRepo, prRepo)
	prUsecase := pr.NewUsecase(userRepo, prRepo)

	// Gin-сервер
	r := gin.Default()
	r.Use(middleware.AuthMiddleware()) // Middleware

	// Handlers
	handlers := http.NewHandlers(teamUsecase, userUsecase, prUsecase)
	
	// УДАЛЕНО: второй вызов RegisterRoutes
	handlers.RegisterRoutes(r) // ← ТОЛЬКО ОДИН РАЗ!

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Логирование
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	_ = logger

	return r
}