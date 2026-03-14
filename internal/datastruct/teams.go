package datastruct

type CreateTeamRequest struct {
	JWTCreds
	Name        string `json:"name" validate:"required" example:"core"`
	Description string `json:"description" validate:"required" example:"responsible for transactions"`
}

type CreateTeamResponse struct {
	Status
}

type ListUserTeamsRequest struct {
	JWTCreds
}

type TeamDescription struct {
	Name        string `json:"name" validate:"required" example:"core"`
	Description string `json:"description" validate:"required" example:"responsible for transactions"`
	TeamId      int64  `json:"team_id,omitempty" validate:"required"`
}

type ListUserTeamsResponse struct {
	Status
	List []TeamDescription
}

type DBInviteUserToTeamRequest struct {
	InviterId      int64
	TeamId         int64
	UserIdToInvite int64
}

type InviteUserToTeamRequest struct {
	JWTCreds
	UserLoginToInvite string `json:"login" validate:"required" example:"VaKadyk359"`
	TeamId            int64  `uri:"team_id" validate:"required" example:"1" swaggerignore:"true"`
}

type InviteUserToTeamResponse struct {
	Status
}
