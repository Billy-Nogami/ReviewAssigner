package inmemory

import (
	"errors"
	"time"
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
)

type pullRequestRepository struct {
	prs       map[string]*schema.PullRequest
	reviewers map[string][]string // prID -> []userID
}

func NewPullRequestRepository() interfaces.PullRequestRepository {
	return &pullRequestRepository{
		prs:       make(map[string]*schema.PullRequest),
		reviewers: make(map[string][]string),
	}
}

func (r *pullRequestRepository) Create(pr *schema.PullRequest) error {
	if _, exists := r.prs[pr.ID]; exists {
		return errors.New("PR already exists")
	}
	r.prs[pr.ID] = pr
	r.reviewers[pr.ID] = pr.AssignedReviewers
	return nil
}

func (r *pullRequestRepository) GetByID(id string) (*schema.PullRequest, error) {
	pr, exists := r.prs[id]
	if !exists {
		return nil, nil
	}
	pr.AssignedReviewers = r.reviewers[id]
	return pr, nil
}

func (r *pullRequestRepository) UpdateStatus(id string, status string, mergedAt *time.Time) (*schema.PullRequest, error) {
	pr, exists := r.prs[id]
	if !exists {
		return nil, errors.New("PR not found")
	}
	pr.Status = status
	pr.MergedAt = mergedAt
	return pr, nil
}

func (r *pullRequestRepository) UpdateReviewers(id string, reviewers []string) error {
	if _, exists := r.prs[id]; !exists {
		return errors.New("PR not found")
	}
	r.reviewers[id] = reviewers
	return nil
}

func (r *pullRequestRepository) GetByReviewerID(userID string) ([]schema.PullRequestShort, error) {
	var prs []schema.PullRequestShort
	for prID, reviewers := range r.reviewers {
		for _, rID := range reviewers {
			if rID == userID {
				pr := r.prs[prID]
				prs = append(prs, schema.PullRequestShort{
					ID:       pr.ID,
					Name:     pr.Name,
					AuthorID: pr.AuthorID,
					Status:   pr.Status,
				})
				break
			}
		}
	}
	return prs, nil
}

func (r *pullRequestRepository) Exists(id string) (bool, error) {
	_, exists := r.prs[id]
	return exists, nil
}

// Методы для тестов: AddPR для инициализации
func (r *pullRequestRepository) AddPR(pr *schema.PullRequest) {
	r.prs[pr.ID] = pr
	r.reviewers[pr.ID] = pr.AssignedReviewers
}
