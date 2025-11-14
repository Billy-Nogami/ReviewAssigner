package inmemory

import (
	"errors"
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
)

type teamRepository struct {
	teams map[string]*schema.Team
}

func NewTeamRepository() interfaces.TeamRepository {
	return &teamRepository{teams: make(map[string]*schema.Team)}
}

func (r *teamRepository) Create(team *schema.Team) error {
	if _, exists := r.teams[team.Name]; exists {
		return errors.New("team already exists")
	}
	r.teams[team.Name] = team
	return nil
}

func (r *teamRepository) GetByName(name string) (*schema.Team, error) {
	team, exists := r.teams[name]
	if !exists {
		return nil, nil
	}
	return team, nil
}

func (r *teamRepository) Exists(name string) (bool, error) {
	_, exists := r.teams[name]
	return exists, nil
}

// Методы для тестов: AddTeam для инициализации
func (r *teamRepository) AddTeam(team *schema.Team) {
	r.teams[team.Name] = team
}
