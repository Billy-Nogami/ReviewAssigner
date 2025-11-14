package inmemory

import (
	"errors"
	"ReviewAssigner/internal/domain/interfaces"
	"ReviewAssigner/internal/domain/schemas"
)

type userRepository struct {
	users map[string]*schema.User
}

func NewUserRepository() interfaces.UserRepository {
	return &userRepository{users: make(map[string]*schema.User)}
}

func (r *userRepository) GetByID(userID string) (*schema.User, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *userRepository) UpdateIsActive(userID string, isActive bool) (*schema.User, error) {
	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	user.IsActive = isActive
	return user, nil
}

func (r *userRepository) GetActiveByTeam(teamName string, excludeUserID string) ([]schema.User, error) {
	var users []schema.User
	for _, user := range r.users {
		if user.TeamName == teamName && user.IsActive && user.ID != excludeUserID {
			users = append(users, *user)
		}
	}
	return users, nil
}

// Методы для тестов: AddUser для инициализации
func (r *userRepository) AddUser(user *schema.User) {
	r.users[user.ID] = user
}