package main

import (
	"log"
	"os"

	"ReviewAssigner/internal/delivery/http"
	"ReviewAssigner/internal/delivery/middleware"
	"ReviewAssigner/internal/repository/postgres"
	"ReviewAssigner/internal/usecase/pr"
	"ReviewAssigner/internal/usecase/team"
	"ReviewAssigner/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/slog"
)

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

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (публичный)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Login (публичный) - временная заглушка
	r.POST("/auth/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"token": "test-jwt-token-12345"})
	})

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		handlers := http.NewHandlers(teamUsecase, userUsecase, prUsecase)
		
		// Регистрируем routes в protected группе
		protected.POST("/team/add", handlers.CreateTeam)
		protected.GET("/team/get", handlers.GetTeam)
		protected.POST("/users/setIsActive", handlers.SetUserActive)
		protected.GET("/users/getReview", handlers.GetUserReviews)
		protected.POST("/pullRequest/create", handlers.CreatePR)
		protected.POST("/pullRequest/merge", handlers.MergePR)
		protected.POST("/pullRequest/reassign", handlers.ReassignPR)
		protected.GET("/stats", handlers.GetStats)
	}

	// Логирование
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	_ = logger

	return r
}
