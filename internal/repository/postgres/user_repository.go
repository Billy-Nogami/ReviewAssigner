  package postgres

  import (
      "database/sql"
      "ReviewAssigner/internal/domain/schemas"
      "ReviewAssigner/internal/domain/interfaces"

      "github.com/jmoiron/sqlx"
  )

  type userRepository struct {
      db *sqlx.DB
  }

  func NewUserRepository(db *sqlx.DB) interfaces.UserRepository {
      return &userRepository{db: db}
  }

  func (r *userRepository) GetByID(userID string) (*schemas.User, error) {
      var user schemas.User
      err := r.db.Get(&user, "SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1", userID)
      if err == sql.ErrNoRows {
          return nil, nil
      }
      return &user, err
  }

  func (r *userRepository) UpdateIsActive(userID string, isActive bool) (*schemas.User, error) {
      _, err := r.db.Exec("UPDATE users SET is_active = $1 WHERE user_id = $2", isActive, userID)
      if err != nil {
          return nil, err
      }
      return r.GetByID(userID)
  }

  func (r *userRepository) GetActiveByTeam(teamName string, excludeUserID string) ([]schemas.User, error) {
      var users []schemas.User
      err := r.db.Select(&users, "SELECT user_id, username, team_name, is_active FROM users WHERE team_name = $1 AND is_active = true AND user_id != $2", teamName, excludeUserID)
      return users, err
  }
