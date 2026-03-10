package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"team-task-manager/internal/supports"
	"time"

	ds "team-task-manager/internal/datastruct"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultAccessTokenTTL  = time.Minute * 15
	defaultRefreshTokenTTL = time.Hour * 24 * 7
	defaultDeleteRTDelay   = time.Second * 30

	jwtSecretKey = "JWT_SECRET"
)

type RegisteredClaims struct {
	jwt.RegisteredClaims
	Role   string
	UserId int64
}

func (s *Service) RegisterUser(req *ds.RegisterRequest) *ds.RegisterResponse {
	hashPassword, err := supports.ArgonHash(req.Password)
	if err != nil {
		s.logger.ErrorKV("RegisterUser.ArgonHash", "error", err.Error())
		return nil
	}

	req.Password = hashPassword

	res, err := s.storage.AddNewUser(req)
	if err != nil {
		s.logger.ErrorKV("RegisterUser.InsertNewUser", "error", err.Error())
		return nil
	}

	return res
}

func (s *Service) LoginUser(req *ds.LoginRequest) *ds.LoginResponse {
	identities, err := s.storage.GetAuthIdentitiesByLogin(req.Login)
	if err != nil {
		s.logger.ErrorKV("LoginUser.GetAuthIdentitiesByLogin", "error", err.Error())
	}

	match, err := supports.IsStringArgonHash(req.Password, identities.Password)
	if err != nil {
		s.logger.ErrorKV("LoginUser.IsStringArgonHash", "error", err.Error())
		return nil
	}

	if !match {
		return &ds.LoginResponse{
			Status: ds.Status{
				Message: ds.StatusWrongLoginOrPassword,
			},
		}
	}

	at, rt, err := s.getTokensPair(identities)
	if err != nil {
		s.logger.ErrorKV("LoginUser.getTokensPair", "error", err.Error())
		return nil
	}

	return &ds.LoginResponse{
		AccessToken:  ds.AccessToken{AccessToken: at},
		RefreshToken: ds.RefreshToken{RefreshToken: rt},
	}
}

func (s *Service) Refresh(req *ds.RefreshRequest) *ds.RefreshResponse {
	oldRT := hashRefreshToken(req.RefreshToken.RefreshToken)
	dbRT, exist, err := s.storage.GetRefreshToken(oldRT)
	if err != nil {
		s.logger.ErrorKV("Refresh.GetRefreshToken", "error", err.Error())
		return nil
	}

	if !exist || dbRT.Revoked || time.Now().After(dbRT.ExpiresAt) {
		return &ds.RefreshResponse{
			Status: ds.Status{Message: ds.StatusInvalidToken},
		}
	}

	identities, err := s.storage.GetAuthIdentitiesByUserID(dbRT.UserID)
	if err != nil {
		s.logger.ErrorKV("Refresh.GetAuthIdentitiesByUserID", "error", err.Error())
		return nil
	}

	at, rt, err := s.getTokensPair(identities)
	if err != nil {
		s.logger.ErrorKV("Refresh.getTokensPair", "error", err.Error())
		return nil
	}

	go func() {
		tt := time.NewTicker(defaultDeleteRTDelay)
		<-tt.C
		err := s.storage.DeleteRefreshToken(oldRT)
		if err != nil {
			s.logger.ErrorKV("Refresh.DeleteRefreshToken", "error", err.Error())
		}
		tt.Stop()
	}()

	return &ds.RefreshResponse{
		AccessToken:  ds.AccessToken{AccessToken: at},
		RefreshToken: ds.RefreshToken{RefreshToken: rt},
	}
}

func (s *Service) getTokensPair(identities *ds.DBAuthIdentities) (at string, rt string, err error) {
	rt = rand.Text()
	err = s.storage.AddRefreshToken(identities,
		hashRefreshToken(rt), time.Now().Add(defaultRefreshTokenTTL))
	if err != nil {
		return
	}

	now := time.Now()
	claim := RegisteredClaims{
		UserId: identities.UserID,
		Role:   identities.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    serviceName,
			ExpiresAt: jwt.NewNumericDate(now.Add(defaultAccessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	sing, err := s.secret.ReadSecret(jwtSecretKey)
	if err != nil {
		return
	}

	atRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	at, err = atRaw.SignedString([]byte(sing))
	if err != nil {
		return
	}

	return
}

func hashRefreshToken(t string) string {
	hash := sha256.Sum256([]byte(t))
	return hex.EncodeToString(hash[:])
}
