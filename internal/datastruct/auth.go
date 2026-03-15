package datastruct

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTSecretKey = "JWT_SECRET"

	UserRoleKey = "role"
	UserIDKey   = "user_id"

	RoleUser  = "user"
	RoleAdmin = "admin"

	DefaultAccessTokenTTL          = time.Minute * 15
	DefaultRefreshTokenTTL         = time.Hour * 24 * 7
	DefaultRefreshTokenGracePeriod = time.Second * 30
)

type UserCreds struct {
	Login    string `json:"login" validate:"required" example:"VaKadyk359"`
	Password string `json:"password" validate:"required" example:"Xldf32Q"`
}

type UserInfo struct {
	Name string `json:"name" validate:"required" example:"Vasilisa"`
}

type JWTCreds struct {
	Role   string `swaggerignore:"true"`
	UserID int64  `swaggerignore:"true"`
}

func (c *JWTCreds) SetUserRole(r string) {
	c.Role = r
}

func (c *JWTCreds) SetUserId(id int64) {
	c.UserID = id
}

type AuthIdentities struct {
	JWTCreds
	UserCreds
	UserInfo
}

type AccessToken struct {
	AccessToken string `json:"access_token,omitempty" example:"-"`
}

type RefreshToken struct {
	RefreshToken string `json:",omitempty" validate:"required" example:"-" swaggerignore:"true"`
}

type DBRefreshToken struct {
	RefreshToken
	ExpiresAt time.Time
	UserID    int64
	Revoked   bool
	Used      bool
}

type DBUpdateRefreshToken struct {
	RefreshToken
	ExpiresAt time.Time
	Revoked   bool
	Used      bool
}

type DBRegisterRequest struct {
	UserCreds
	UserInfo
	Role string
}

type RegisterRequest struct {
	UserCreds
	UserInfo
}

type RegisterResponse struct {
	Status
}

type LoginRequest struct {
	UserCreds
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

type RegisteredClaims struct {
	jwt.RegisteredClaims
	Role   string
	UserId int64
}
