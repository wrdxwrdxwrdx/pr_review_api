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
	PullRequestId     string     `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name" db:"pull_request_name"`
	AuthorId          string     `json:"author_id" db:"author_id"`
	Status            Status     `json:"status" db:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers" db:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt" db:"created_at"`
	MergedAt          *time.Time `json:"mergedAt" db:"merged_at"`
}

type PullRequestShort struct {
	PullRequestId   string `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string `json:"pull_request_name" db:"pull_request_name"`
	AuthorId        string `json:"author_id" db:"author_id"`
	Status          Status `json:"status" db:"status"`
}

type PullRequestJson struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type PullRequestStatistics struct {
	Stats []*UserPr `json:"stats"`
}

type UserPr struct {
	UserId   string `json:"user_id"`
	PrNumber int    `json:"pull_request_number"`
}

type MergePullRequestJson struct {
	PullRequestId string `json:"pull_request_id"`
}

type ReassignPullRequestJson struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

func NewPullRequest(pullRequestId string, pullRequestName string, authorId string, status Status, assignedReviewers []string, createdAt *time.Time) *PullRequest {
	return &PullRequest{PullRequestId: pullRequestId,
		PullRequestName:   pullRequestName,
		AuthorId:          authorId,
		Status:            status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         createdAt,
	}
}

func NewPullRequestShort(pullRequestId string, pullRequestName string, authoeId string, status Status) *PullRequestShort {
	return &PullRequestShort{
		PullRequestId:   pullRequestId,
		PullRequestName: pullRequestName,
		AuthorId:        authoeId,
		Status:          status,
	}
}

func NewPullRequestShortFromPr(pullRequest PullRequest) *PullRequestShort {
	return &PullRequestShort{
		PullRequestId:   pullRequest.PullRequestId,
		PullRequestName: pullRequest.PullRequestName,
		AuthorId:        pullRequest.AuthorId,
		Status:          pullRequest.Status,
	}
}

func NewUserPr(userId string, prNumber int) *UserPr {
	return &UserPr{UserId: userId, PrNumber: prNumber}
}

func NewPullRequestStatistics(stats []*UserPr) *PullRequestStatistics {
	return &PullRequestStatistics{Stats: stats}
}
