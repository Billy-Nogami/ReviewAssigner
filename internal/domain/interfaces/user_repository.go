package interfaces

import "ReviewAssigner/internal/domain/schemas"

type UserRepository interface {
    GetByID(userID string) (*schemas.User, error)
    UpdateIsActive(userID string, isActive bool) (*schemas.User, error)
    GetActiveByTeam(teamName string, excludeUserID string) ([]schemas.User, error) // Для выбора ревьюверов
}
  