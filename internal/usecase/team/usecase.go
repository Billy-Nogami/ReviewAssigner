package team

import (
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
	"ReviewAssigner/internal/pkg/errors"
)

type Usecase struct {
	teamRepo interfaces.TeamRepository
}

func NewUsecase(teamRepo interfaces.TeamRepository) *Usecase {
	return &Usecase{teamRepo: teamRepo}
}

func (u *Usecase) CreateTeam(team *schemas.Team) (*schemas.Team, error) {
	exists, err := u.teamRepo.Exists(team.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrTeamExists
	}
	err = u.teamRepo.Create(team)
	if err != nil {
		return nil, err
	}
	return u.teamRepo.GetByName(team.Name)
}

func (u *Usecase) GetTeam(name string) (*schemas.Team, error) {
	team, err := u.teamRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.ErrNotFound
	}
	return team, nil
}
