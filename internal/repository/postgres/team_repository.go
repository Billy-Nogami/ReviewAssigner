package postgres

import (
    // "database/sql" // УДАЛЕНО: неиспользуемый импорт
    "ReviewAssigner/internal/domain/schemas"
    "ReviewAssigner/internal/domain/interfaces"

    "github.com/jmoiron/sqlx"
)

type teamRepository struct {
    db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) interfaces.TeamRepository {
    return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *schemas.Team) error {
    tx, err := r.db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("INSERT INTO teams (team_name) VALUES ($1)", team.Name)
    if err != nil {
        return err
    }

    for _, member := range team.Members {
        _, err = tx.Exec("INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active",
            member.ID, member.Username, team.Name, member.IsActive)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *teamRepository) GetByName(name string) (*schemas.Team, error) {
    var team schemas.Team
    team.Name = name
    err := r.db.Select(&team.Members, "SELECT user_id, username, team_name, is_active FROM users WHERE team_name = $1", name)
    if err != nil {
        return nil, err
    }
    return &team, nil
}

func (r *teamRepository) Exists(name string) (bool, error) {
    var count int
    err := r.db.Get(&count, "SELECT COUNT(*) FROM teams WHERE team_name = $1", name)
    return count > 0, err
}
