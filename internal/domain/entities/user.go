package entities

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserJson struct {
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	// IsActive bool   `json:"is_active"`
}

type UserIsActiveJson struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func NewUser(userId string, username string, teamName string, isActive bool) *User {
	return &User{UserId: userId, Username: username, TeamName: teamName, IsActive: isActive}
}
