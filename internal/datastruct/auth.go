package datastruct

import "time"

type AuthIdentities struct {
	Login    string `json:"login" validate:"required" example:"Vasilisa"`
	Password string `json:"password" validate:"required" example:"Xldf32Q"`
}

type AccessToken struct {
	AccessToken string `json:"access_token,omitempty" example:"-"`
}

type RefreshToken struct {
	RefreshToken string `cookie:"refresh_token" validate:"required" example:"-"`
}

type DBAuthIdentities struct {
	AuthIdentities
	Role   string
	UserID int64
}

type DBRefreshToken struct {
	RefreshToken
	ExpiresAt time.Time
	UserID    int64
	Revoked   bool
}

type RegisterRequest struct {
	AuthIdentities
}

type RegisterResponse struct {
	Status
}

type LoginRequest struct {
	AuthIdentities
}

type LoginResponse struct {
	Status
	AccessToken
	RefreshToken
}

type RefreshRequest struct {
	RefreshToken
}

type RefreshResponse struct {
	Status
	AccessToken
	RefreshToken
}
