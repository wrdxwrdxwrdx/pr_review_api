package customerrors

import "fmt"

type ErrorCode string

const (
	TeamExists   ErrorCode = "TEAM_EXISTS"
	PRExists     ErrorCode = "PR_EXISTS"
	PRMerged     ErrorCode = "PR_MERGED"
	NotAssigned  ErrorCode = "NOT_ASSIGNED"
	NoCandidate  ErrorCode = "NO_CANDIDATE"
	NotFound     ErrorCode = "NOT_FOUND"
	Unauthorized ErrorCode = "UNAUTHORIZED"
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewDomainError(code ErrorCode, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

func NewDomainErrorWithErr(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{Code: code, Message: message, Err: err}
}
