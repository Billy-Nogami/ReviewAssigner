package interfaces

import "ReviewAssigner/internal/domain/schemas"

type TeamRepository interface {
    Create(team *schemas.Team) error
    GetByName(name string) (*schemas.Team, error)
    Exists(name string) (bool, error)
}
