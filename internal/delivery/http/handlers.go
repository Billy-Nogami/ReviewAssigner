// internal/delivery/http/handlers.go
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

// Публичные
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

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

	c.JSON(200, gin.H{
		"token": token,
		"role":  role,
	})
}

// Защищённые хендлеры (все остальные — оставляем как были, но без @Summary и т.д.)

func (h *Handlers) CreateTeam(c *gin.Context) {
	var req struct {
		Name    string          `json:"name" binding:"required"`
		Members []schemas.User `json:"members" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	teamData := &schemas.Team{Name: req.Name, Members: req.Members}
	team, err := h.teamUsecase.CreateTeam(teamData)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(201, gin.H{"team": team})
}

func (h *Handlers) GetTeam(c *gin.Context) {
	name := c.Query("team_name")
	if name == "" {
		c.JSON(400, gin.H{"error": "team_name query param is required"})
		return
	}
	team, err := h.teamUsecase.GetTeam(name)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, team)
}

func (h *Handlers) SetUserActive(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		IsActive bool   `json:"is_active"`
	}
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

func (h *Handlers) GetUserReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "user_id query param is required"})
		return
	}
	user, prs, err := h.userUsecase.GetUserReviews(userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{"user_id": user.ID, "pull_requests": prs})
}

func (h *Handlers) CreatePR(c *gin.Context) {
	var req struct {
		PRID   string `json:"pull_request_id" binding:"required"`
		Name   string `json:"pull_request_name" binding:"required"`
		Author string `json:"author_id" binding:"required"`
	}
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

func (h *Handlers) MergePR(c *gin.Context) {
	var req struct {
		PRID string `json:"pull_request_id" binding:"required"`
	}
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

func (h *Handlers) ReassignPR(c *gin.Context) {
	var req struct {
		PRID      string `json:"pull_request_id" binding:"required"`
		OldUserID string `json:"old_user_id" binding:"required"`
	}
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

func (h *Handlers) GetStats(c *gin.Context) {
	userStats, prStats, err := h.prUsecase.GetStats()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"user_assignments": userStats,
		"pr_assignments":   prStats,
	})
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
