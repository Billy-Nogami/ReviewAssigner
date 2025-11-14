package pr

import (
	"math/rand"
	"time"
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
	"ReviewAssigner/internal/pkg/errors"
)

type Usecase struct {
	userRepo interfaces.UserRepository
	prRepo   interfaces.PullRequestRepository
}

func NewUsecase(userRepo interfaces.UserRepository, prRepo interfaces.PullRequestRepository) *Usecase {
	return &Usecase{userRepo: userRepo, prRepo: prRepo}
}

func (u *Usecase) CreatePR(prID, name, authorID string) (*schemas.PullRequest, error) {
	exists, err := u.prRepo.Exists(prID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPRExists
	}

	author, err := u.userRepo.GetByID(authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.ErrNotFound
	}

	candidates, err := u.userRepo.GetActiveByTeam(author.TeamName, authorID)
	if err != nil {
		return nil, err
	}

	// Случайный выбор до 2 ревьюверов
	rand.New(rand.NewSource(time.Now().UnixNano())) // Исправлено: вместо rand.Seed
	selected := []string{}
	if len(candidates) > 0 {
		perm := rand.Perm(len(candidates))
		for i := 0; i < len(perm) && i < 2; i++ {
			selected = append(selected, candidates[perm[i]].ID)
		}
	}

	// Исправление: создаем переменную для времени
	createdAt := time.Now()
	pr := &schemas.PullRequest{
		ID:                prID,
		Name:              name,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: selected,
		CreatedAt:         &createdAt, // Исправлено: используем переменную
	}

	err = u.prRepo.Create(pr)
	return pr, err
}

func (u *Usecase) MergePR(prID string) (*schemas.PullRequest, error) {
	pr, err := u.prRepo.GetByID(prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, errors.ErrNotFound
	}
	if pr.Status == "MERGED" {
		return pr, nil // Идемпотентность
	}
	mergedAt := time.Now()
	return u.prRepo.UpdateStatus(prID, "MERGED", &mergedAt)
}

func (u *Usecase) ReassignPR(prID, oldUserID string) (*schemas.PullRequest, string, error) {
	pr, err := u.prRepo.GetByID(prID)
	if err != nil {
		return nil, "", err
	}
	if pr == nil {
		return nil, "", errors.ErrNotFound
	}
	if pr.Status == "MERGED" {
		return nil, "", errors.ErrPRMerged
	}

	// Проверить, что oldUserID назначен
	assigned := false
	for _, r := range pr.AssignedReviewers {
		if r == oldUserID {
			assigned = true
			break
		}
	}
	if !assigned {
		return nil, "", errors.ErrNotAssigned
	}

	// Найти команду oldUserID
	oldUser, err := u.userRepo.GetByID(oldUserID)
	if err != nil {
		return nil, "", err
	}
	if oldUser == nil {
		return nil, "", errors.ErrNotFound
	}

	// Кандидаты из команды oldUser (активные, исключая автора и уже назначенных)
	candidates, err := u.userRepo.GetActiveByTeam(oldUser.TeamName, pr.AuthorID)
	if err != nil {
		return nil, "", err
	}
	validCandidates := []string{}
	for _, c := range candidates {
		if c.ID != oldUserID { // Исключаем заменяемого
			alreadyAssigned := false
			for _, a := range pr.AssignedReviewers {
				if a == c.ID {
					alreadyAssigned = true
					break
				}
			}
			if !alreadyAssigned {
				validCandidates = append(validCandidates, c.ID)
			}
		}
	}
	if len(validCandidates) == 0 {
		return nil, "", errors.ErrNoCandidate
	}

	// Случайный выбор
	rand.New(rand.NewSource(time.Now().UnixNano())) // Исправлено: вместо rand.Seed
	newReviewer := validCandidates[rand.Intn(len(validCandidates))]

	// Обновить список
	newReviewers := []string{}
	for _, r := range pr.AssignedReviewers {
		if r != oldUserID {
			newReviewers = append(newReviewers, r)
		}
	}
	newReviewers = append(newReviewers, newReviewer)

	err = u.prRepo.UpdateReviewers(prID, newReviewers)
	if err != nil {
		return nil, "", err
	}
	pr, _ = u.prRepo.GetByID(prID)
	return pr, newReviewer, nil
}
