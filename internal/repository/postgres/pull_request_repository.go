package postgres

import (
    "database/sql"
    "time"
    "ReviewAssigner/internal/domain/schemas"
    "ReviewAssigner/internal/domain/interfaces"

    "github.com/jmoiron/sqlx"
)

type pullRequestRepository struct {
    db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) interfaces.PullRequestRepository {
    return &pullRequestRepository{db: db}
}

func (r *pullRequestRepository) Create(pr *schemas.PullRequest) error {
    tx, err := r.db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at) VALUES ($1, $2, $3, $4, $5)",
        pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt)
    if err != nil {
      return err
    }

    for _, reviewerID := range pr.AssignedReviewers {
        _, err = tx.Exec("INSERT INTO pr_reviewers (pull_request_id, user_id) VALUES ($1, $2)", pr.ID, reviewerID)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *pullRequestRepository) GetByID(id string) (*schemas.PullRequest, error) {
    var pr schemas.PullRequest
    err := r.db.Get(&pr, "SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at FROM pull_requests WHERE pull_request_id = $1", id)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var reviewers []string
    err = r.db.Select(&reviewers, "SELECT user_id FROM pr_reviewers WHERE pull_request_id = $1", id)
    pr.AssignedReviewers = reviewers
    return &pr, err
}

func (r *pullRequestRepository) UpdateStatus(id string, status string, mergedAt *time.Time) (*schemas.PullRequest, error) {
    _, err := r.db.Exec("UPDATE pull_requests SET status = $1, merged_at = $2 WHERE pull_request_id = $3", status, mergedAt, id)
    if err != nil {
        return nil, err
    }
    return r.GetByID(id)
}

func (r *pullRequestRepository) UpdateReviewers(id string, reviewers []string) error {
    tx, err := r.db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("DELETE FROM pr_reviewers WHERE pull_request_id = $1", id)
    if err != nil {
        return err
    }

    for _, reviewerID := range reviewers {
        _, err = tx.Exec("INSERT INTO pr_reviewers (pull_request_id, user_id) VALUES ($1, $2)", id, reviewerID)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *pullRequestRepository) GetByReviewerID(userID string) ([]schemas.PullRequestShort, error) {
    var prs []schemas.PullRequestShort
    err := r.db.Select(&prs, "SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status FROM pull_requests pr JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id WHERE prr.user_id = $1", userID)
    return prs, err
}

func (r *pullRequestRepository) Exists(id string) (bool, error) {
    var count int
    err := r.db.Get(&count, "SELECT COUNT(*) FROM pull_requests WHERE pull_request_id = $1", id)
    return count > 0, err
}

func (r *pullRequestRepository) GetStats() (map[string]int, map[string]int, error) {
    userStats := make(map[string]int)
    prStats := make(map[string]int)

    // Статистика по пользователям
    rows, err := r.db.Query("SELECT user_id, COUNT(*) FROM pr_reviewers GROUP BY user_id")
    if err != nil {
        return nil, nil, err
    }
    defer rows.Close()
    for rows.Next() {
        var userID string
        var count int
        if err := rows.Scan(&userID, &count); err != nil {
            return nil, nil, err
        }
        userStats[userID] = count
    }

    // Статистика по PR
    rows2, err := r.db.Query("SELECT pull_request_id, COUNT(*) FROM pr_reviewers GROUP BY pull_request_id")
    if err != nil {
        return nil, nil, err
    }
    defer rows2.Close()
    for rows2.Next() {
        var prID string
        var count int
        if err := rows2.Scan(&prID, &count); err != nil {
            return nil, nil, err
        }
        prStats[prID] = count
    }

    return userStats, prStats, nil
}
