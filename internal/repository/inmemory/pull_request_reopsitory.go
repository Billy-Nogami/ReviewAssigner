package inmemory

import (
	"errors"
	"sync"
	"time"
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
)

type pullRequestRepository struct {
	mu        sync.RWMutex
	prs       map[string]*schemas.PullRequest
	reviewers map[string][]string // prID -> []userID
}

func NewPullRequestRepository() interfaces.PullRequestRepository {
	return &pullRequestRepository{
		prs:       make(map[string]*schemas.PullRequest),
		reviewers: make(map[string][]string),
	}
}

func (r *pullRequestRepository) Create(pr *schemas.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.ID]; exists {
		return errors.New("PR already exists")
	}
	r.prs[pr.ID] = pr
	r.reviewers[pr.ID] = pr.AssignedReviewers
	return nil
}

func (r *pullRequestRepository) GetByID(id string) (*schemas.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pr, exists := r.prs[id]
	if !exists {
		return nil, nil
	}
	pr.AssignedReviewers = r.reviewers[id]
	return pr, nil
}

func (r *pullRequestRepository) UpdateStatus(id string, status string, mergedAt *time.Time) (*schemas.PullRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pr, exists := r.prs[id]
	if !exists {
		return nil, errors.New("PR not found")
	}
	pr.Status = status
	pr.MergedAt = mergedAt
	return pr, nil
}

func (r *pullRequestRepository) UpdateReviewers(id string, reviewers []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[id]; !exists {
		return errors.New("PR not found")
	}
	r.reviewers[id] = reviewers
	return nil
}

func (r *pullRequestRepository) GetByReviewerID(userID string) ([]schemas.PullRequestShort, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var prs []schemas.PullRequestShort
	for prID, reviewers := range r.reviewers {
		for _, rID := range reviewers {
			if rID == userID {
				pr := r.prs[prID]
				prs = append(prs, schemas.PullRequestShort{
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.prs[id]
	return exists, nil
}

// GetStats возвращает статистику по pull requests
func (r *pullRequestRepository) GetStats() (map[string]int, map[string]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userStats := make(map[string]int)  // userID -> количество ревью
	prStats := make(map[string]int)    // prID -> количество ревьюверов

	// Собираем статистику по пользователям (сколько ревью у каждого)
	for prID, reviewers := range r.reviewers {
		// Статистика по PR
		prStats[prID] = len(reviewers)

		// Статистика по пользователям
		for _, userID := range reviewers {
			userStats[userID]++
		}
	}

	return userStats, prStats, nil
}

// Методы для тестов: AddPR для инициализации
func (r *pullRequestRepository) AddPR(pr *schemas.PullRequest) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.prs[pr.ID] = pr
	r.reviewers[pr.ID] = pr.AssignedReviewers
}
