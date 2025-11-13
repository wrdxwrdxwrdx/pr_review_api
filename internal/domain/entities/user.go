package entities

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `json:"UserId"`
	Username string    `json:"username"`
	TeamName string    `json:"team_name"`
	IsActive bool      `json:"is_active"`
}

type UserJson struct {
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	// IsActive bool   `json:"is_active"`
}

func NewUser(userId uuid.UUID, username string, teamName string, isActive bool) *User {
	return &User{UserId: userId, Username: username, TeamName: teamName, IsActive: isActive}
}
