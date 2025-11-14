package entities

type TeamMember struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func NewTeam(teamName string, members []TeamMember) *Team {
	return &Team{TeamName: teamName, Members: members}
}

func NewTeamMember(userId string, username string, isActive bool) *TeamMember {
	return &TeamMember{UserId: userId, Username: username, IsActive: isActive}
}
