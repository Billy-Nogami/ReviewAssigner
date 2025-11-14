package user

import (
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

func (u *Usecase) SetIsActive(userID string, isActive bool) (*schemas.User, error) {
	user, err := u.userRepo.UpdateIsActive(userID, isActive)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrNotFound
	}
	return user, nil
}

func (u *Usecase) GetUserReviews(userID string) (*schemas.User, []schemas.PullRequestShort, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.ErrNotFound
	}
	prs, err := u.prRepo.GetByReviewerID(userID)
	return user, prs, err
}
