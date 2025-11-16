// @title Review Assigner API
// @version 1.0
// @description API for assigning code reviewers to pull requests

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @security BearerAuth

package http

import (
	"ReviewAssigner/internal/domain/schemas"
	"ReviewAssigner/internal/pkg/errors"
	"ReviewAssigner/internal/usecase/pr"
	"ReviewAssigner/internal/usecase/team"
	"ReviewAssigner/internal/usecase/user"
	"ReviewAssigner/internal/delivery/middleware"
	"ReviewAssigner/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	teamUsecase *team.Usecase
	userUsecase *user.Usecase
	prUsecase   *pr.Usecase
}

func NewHandlers(teamUsecase *team.Usecase, userUsecase *user.Usecase, prUsecase *pr.Usecase) *Handlers {
	return &Handlers{
		teamUsecase: teamUsecase,
		userUsecase: userUsecase,
		prUsecase:   prUsecase,
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)

	// Public routes
	r.POST("/auth/login", h.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/team/add", h.CreateTeam)
		protected.GET("/team/get", h.GetTeam)
		protected.POST("/users/setIsActive", h.SetUserActive)
		protected.GET("/users/getReview", h.GetUserReviews)
		protected.POST("/pullRequest/create", h.CreatePR)
		protected.POST("/pullRequest/merge", h.MergePR)
		protected.POST("/pullRequest/reassign", h.ReassignPR)
		protected.GET("/stats", h.GetStats)
	}
}

// Модели запросов для Swagger

// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// CreateTeamRequest представляет запрос на создание команды
type CreateTeamRequest struct {
	// @Example: backend
	Name string `json:"name" example:"backend"`
	Members []schemas.User `json:"members"`
}

// @Summary Создать команду
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body CreateTeamRequest true "Team data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /team/add [post]
func (h *Handlers) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	teamData := &schemas.Team{
		Name:    req.Name,
		Members: req.Members,
	}

	team, err := h.teamUsecase.CreateTeam(teamData)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(201, gin.H{"team": team})
}

// @Summary Получить команду
// @Tags Teams
// @Produce json
// @Param team_name query string true "Team name" Example("backend")
// @Success 200 {object} schemas.Team
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /team/get [get]
func (h *Handlers) GetTeam(c *gin.Context) {
	name := c.Query("team_name")
	team, err := h.teamUsecase.GetTeam(name)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, team)
}

// SetUserActiveRequest представляет запрос на изменение активности пользователя
type SetUserActiveRequest struct {
	// @Example: user123
	UserID string `json:"user_id" example:"user123"`
	// @Example: true
	IsActive bool `json:"is_active" example:"true"`
}

// @Summary Установить активность пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body SetUserActiveRequest true "User active data"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users/setIsActive [post]
func (h *Handlers) SetUserActive(c *gin.Context) {
	var req SetUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	user, err := h.userUsecase.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"user": user})
}

// @Summary Получить PR пользователя
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID" Example("user123")
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users/getReview [get]
func (h *Handlers) GetUserReviews(c *gin.Context) {
	userID := c.Query("user_id")
	user, prs, err := h.userUsecase.GetUserReviews(userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"user_id": user.ID, "pull_requests": prs})
}

// CreatePRRequest представляет запрос на создание PR
type CreatePRRequest struct {
	// @Example: pr-123
	PRID string `json:"pull_request_id" example:"pr-123"`
	// @Example: "Add new feature"
	Name string `json:"pull_request_name" example:"Add new feature"`
	// @Example: user123
	Author string `json:"author_id" example:"user123"`
}

// @Summary Создать PR
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body CreatePRRequest true "PR data"
// @Success 201 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pullRequest/create [post]
func (h *Handlers) CreatePR(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	pr, err := h.prUsecase.CreatePR(req.PRID, req.Name, req.Author)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(201, gin.H{"pr": pr})
}

// MergePRRequest представляет запрос на мерж PR
type MergePRRequest struct {
	// @Example: pr-123
	PRID string `json:"pull_request_id" example:"pr-123"`
}

// @Summary Merge PR
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body MergePRRequest true "PR ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pullRequest/merge [post]
func (h *Handlers) MergePR(c *gin.Context) {
	var req MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	pr, err := h.prUsecase.MergePR(req.PRID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"pr": pr})
}

// ReassignPRRequest представляет запрос на перераспределение ревьювера
type ReassignPRRequest struct {
	// @Example: pr-123
	PRID string `json:"pull_request_id" example:"pr-123"`
	// @Example: user123
	OldUserID string `json:"old_user_id" example:"user123"`
}

// @Summary Reassign reviewer
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body ReassignPRRequest true "Reassign data"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pullRequest/reassign [post]
func (h *Handlers) ReassignPR(c *gin.Context) {
	var req ReassignPRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	pr, newReviewer, err := h.prUsecase.ReassignPR(req.PRID, req.OldUserID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"pr": pr, "replaced_by": newReviewer})
}

// @Summary Получить статистику назначений
// @Tags Statistics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /stats [get]
func (h *Handlers) GetStats(c *gin.Context) {
	userStats, prStats, err := h.prUsecase.GetStats()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"user_assignments": userStats, "pr_assignments": prStats})
}

// LoginRequest представляет запрос на аутентификацию
type LoginRequest struct {
	// @Example: admin
	UserID string `json:"user_id" example:"admin"`
	// @Example: admin
	Password string `json:"password" example:"admin"`
}

// LoginResponse представляет ответ на аутентификацию
type LoginResponse struct {
	// @Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// @Example: admin
	Role string `json:"role" example:"admin"`
}

// @Summary Логин для получения JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login data"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	// Простая проверка (в проде — из БД)
	var role string
	if req.UserID == "admin" && req.Password == "admin" {
		role = "admin"
	} else if req.UserID == "user" && req.Password == "user" {
		role = "user"
	} else {
		c.JSON(401, gin.H{"error": gin.H{"code": "INVALID_CREDENTIALS", "message": "Invalid user_id or password"}})
		return
	}
	token, err := jwt.GenerateToken(req.UserID, role)
	if err != nil {
		c.JSON(500, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(200, LoginResponse{
		Token: token,
		Role:  role,
	})
}

// ErrorResponse представляет стандартный ответ об ошибке
type ErrorResponse struct {
	Error struct {
		// @Example: NOT_FOUND
		Code    string `json:"code" example:"NOT_FOUND"`
		// @Example: Resource not found
		Message string `json:"message" example:"Resource not found"`
	} `json:"error"`
}

func handleError(c *gin.Context, err error) {
	switch err {
	case errors.ErrTeamExists:
		c.JSON(400, gin.H{"error": gin.H{"code": "TEAM_EXISTS", "message": "team_name already exists"}})
	case errors.ErrPRExists:
		c.JSON(409, gin.H{"error": gin.H{"code": "PR_EXISTS", "message": "PR id already exists"}})
	case errors.ErrPRMerged:
		c.JSON(409, gin.H{"error": gin.H{"code": "PR_MERGED", "message": "cannot reassign on merged PR"}})
	case errors.ErrNotAssigned:
		c.JSON(409, gin.H{"error": gin.H{"code": "NOT_ASSIGNED", "message": "reviewer is not assigned to this PR"}})
	case errors.ErrNoCandidate:
		c.JSON(409, gin.H{"error": gin.H{"code": "NO_CANDIDATE", "message": "no active replacement candidate in team"}})
	case errors.ErrNotFound:
		c.JSON(404, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "resource not found"}})
	default:
		c.JSON(500, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
	}
}
