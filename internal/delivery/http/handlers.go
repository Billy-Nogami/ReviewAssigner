package http

import (
	"ReviewAssigner/internal/domain/schemas"
	"ReviewAssigner/internal/pkg/errors"
	"ReviewAssigner/internal/usecase/pr"
	"ReviewAssigner/internal/usecase/team"
	"ReviewAssigner/internal/usecase/user"

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
	r.POST("/team/add", h.CreateTeam)
	r.GET("/team/get", h.GetTeam)
	r.POST("/users/setIsActive", h.SetUserActive)
	r.GET("/users/getReview", h.GetUserReviews)
	r.POST("/pullRequest/create", h.CreatePR)
	r.POST("/pullRequest/merge", h.MergePR)
	r.POST("/pullRequest/reassign", h.ReassignPR)
}

// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// @Summary Создать команду
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body schemas.Team true "Team data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /team/add [post]
func (h *Handlers) CreateTeam(c *gin.Context) {
	var req schemas.Team
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}
	team, err := h.teamUsecase.CreateTeam(&req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(201, gin.H{"team": team})
}

// @Summary Получить команду
// @Tags Teams
// @Produce json
// @Param team_name query string true "Team name"
// @Success 200 {object} schemas.Team
// @Failure 404 {object} map[string]interface{}
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

// @Summary Установить активность пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "User active data"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/setIsActive [post]
func (h *Handlers) SetUserActive(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id"`
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

// @Summary Получить PR пользователя
// @Tags Users
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
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

// @Summary Создать PR
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "PR data"
// @Success 201 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /pullRequest/create [post]
func (h *Handlers) CreatePR(c *gin.Context) {
	var req struct {
		PRID   string `json:"pull_request_id"`
		Name   string `json:"pull_request_name"`
		Author string `json:"author_id"`
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

// @Summary Merge PR
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "PR ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /pullRequest/merge [post]
func (h *Handlers) MergePR(c *gin.Context) {
	var req struct {
		PRID string `json:"pull_request_id"`
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

// @Summary Reassign reviewer
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Reassign data"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /pullRequest/reassign [post]
func (h *Handlers) ReassignPR(c *gin.Context) {
	var req struct {
		PRID     string `json:"pull_request_id"`
		OldUserID string `json:"old_user_id"`
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