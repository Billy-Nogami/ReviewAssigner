package interfaces

import (
	"ReviewAssigner/internal/domain/schemas"
	"time"
)

type PullRequestRepository interface {
    Create(pr *schemas.PullRequest) error
    GetByID(id string) (*schemas.PullRequest, error)
    UpdateStatus(id string, status string, mergedAt *time.Time) (*schemas.PullRequest, error)
    UpdateReviewers(id string, reviewers []string) error
    GetByReviewerID(userID string) ([]schemas.PullRequestShort, error)
    Exists(id string) (bool, error)
    GetStats() (map[string]int, map[string]int, error) // userStats, prStats
}
