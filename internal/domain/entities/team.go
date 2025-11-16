package entities

type TeamMember struct {
	UserId   string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name" db:"teamname"`
	Members  []TeamMember `json:"members"`
}

func NewTeam(teamName string, members []TeamMember) *Team {
	return &Team{TeamName: teamName, Members: members}
}

func NewTeamMember(userId string, username string, isActive bool) *TeamMember {
	return &TeamMember{UserId: userId, Username: username, IsActive: isActive}
}
