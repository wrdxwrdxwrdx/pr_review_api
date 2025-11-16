package entities

type User struct {
	UserId   string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	TeamName string `json:"team_name" db:"teamname"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type UserJson struct {
	Username string `json:"username"`
	TeamName string `json:"team_name"`
}

type UserIsActiveJson struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type Reviews struct {
	UserId       string              `json:"user_id"`
	PullRequests []*PullRequestShort `json:"pull_requests"`
}

func NewReviews(userId string, pullRequests []*PullRequestShort) *Reviews {
	return &Reviews{UserId: userId, PullRequests: pullRequests}
}

func NewUser(userId string, username string, teamName string, isActive bool) *User {
	return &User{UserId: userId, Username: username, TeamName: teamName, IsActive: isActive}
}
