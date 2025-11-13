package entities

import "github.com/google/uuid"

type TeamMember struct {
	User_id   uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Is_active bool      `json:"is_active"`
}

type Team struct {
	Team_name string       `json:"team_name"`
	Members   []TeamMember `json:"members"`
}
