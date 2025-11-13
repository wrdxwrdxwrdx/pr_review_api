package entities

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

type PullRequest struct {
	Pull_request_id    uuid.UUID  `json:"pull_request_id"`
	Pull_request_name  string     `json:"pull_request_name"`
	Author_id          uuid.UUID  `json:"author_id"`
	Status             Status     `json:"status"`
	Assigned_reviewers []string   `json:"assigned_reviewers"`
	CreatedAt          *time.Time `json:"createdAt"`
	MergedAt           *time.Time `json:"mergedAt"`
}
