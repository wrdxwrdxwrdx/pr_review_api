package entities

import (
	"time"
)

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

type PullRequest struct {
	PullRequestId     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorId          string     `json:"author_id"`
	Status            Status     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
}

type PullRequestJson struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type MergePullRequestJson struct {
	PullRequestId string `json:"pull_request_id"`
}

type ReassignPullRequestJson struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

func NewPullRequest(pullRequestId string, pullRequestName string, authoeId string, status Status, assignedReviewers []string, createdAt *time.Time) *PullRequest {
	return &PullRequest{PullRequestId: pullRequestId,
		PullRequestName:   pullRequestName,
		AuthorId:          authoeId,
		Status:            status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         createdAt,
	}
}
