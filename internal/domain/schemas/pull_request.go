package schemas

import "time"

type PullRequest struct {
    ID                string    `json:"pull_request_id" db:"pull_request_id"`
    Name              string    `json:"pull_request_name" db:"pull_request_name"`
    AuthorID          string    `json:"author_id" db:"author_id"`
    Status            string    `json:"status" db:"status"` // OPEN or MERGED
    AssignedReviewers []string  `json:"assigned_reviewers"` // Не в БД напрямую, вычисляется из pr_reviewers
    CreatedAt         *time.Time `json:"createdAt,omitempty" db:"created_at"`
    MergedAt          *time.Time `json:"mergedAt,omitempty" db:"merged_at"`
  }

  type PullRequestShort struct {
    ID       string `json:"pull_request_id"`
    Name     string `json:"pull_request_name"`
    AuthorID string `json:"author_id"`
    Status   string `json:"status"`
  }

  type PRStats struct {
	TotalPRs    int `json:"total_prs"`
	OpenPRs     int `json:"open_prs"`
	MergedPRs   int `json:"merged_prs"`
	ReviewCount int `json:"review_count"`
}
