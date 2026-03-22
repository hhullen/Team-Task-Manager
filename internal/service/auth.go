package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"team-task-manager/internal/supports"
	"time"

	ds "team-task-manager/internal/datastruct"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) RegisterUser(req *ds.RegisterRequest) *ds.RegisterResponse {
	hashPassword, err := supports.ArgonHash(req.Password)
	if err != nil {
		s.logger.ErrorKV("RegisterUser.ArgonHash", "error", err.Error())
		return nil
	}

	req.Password = hashPassword

	res, err := s.storageAuth.AddNewUser(&ds.DBRegisterRequest{
		UserCreds: req.UserCreds,
		UserInfo:  req.UserInfo,
		Role:      ds.RoleUser,
	})
	if err != nil {
		s.logger.ErrorKV("RegisterUser.InsertNewUser", "error", err.Error())
		return nil
	}

	return res
}

func (s *Service) LoginUser(req *ds.LoginRequest) *ds.LoginResponse {
	identities, exist, err := s.storageAuth.GetAuthIdentitiesByLogin(req.Login)
	if err != nil {
		s.logger.ErrorKV("LoginUser.GetAuthIdentitiesByLogin", "error", err.Error())
	}

	if !exist {
		return &ds.LoginResponse{
			Status: ds.Status{Message: ds.StatusWrongLoginOrPassword},
		}
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

	at, rt, err := s.createTokensPair(identities)
	if err != nil {
		s.logger.ErrorKV("LoginUser.createTokensPair", "error", err.Error())
		return nil
	}

	return &ds.LoginResponse{
		AccessToken:  ds.AccessToken{AccessToken: at},
		Status:       ds.Status{Message: ds.StatusSuccess},
		RefreshToken: ds.RefreshToken{RefreshToken: rt},
	}
}

func (s *Service) Refresh(req *ds.RefreshRequest) *ds.RefreshResponse {
	oldRT := hashRefreshToken(req.RefreshToken.RefreshToken)
	dbRT, exist, err := s.storageAuth.GetRefreshToken(oldRT)
	if err != nil {
		s.logger.ErrorKV("Refresh.GetRefreshToken", "error", err.Error())
		return nil
	}

	if !exist || dbRT.Revoked {
		return &ds.RefreshResponse{
			Status: ds.Status{Message: ds.StatusInvalidToken},
		}
	}

	expired := time.Now().After(dbRT.ExpiresAt)
	if expired && !dbRT.Used {
		return &ds.RefreshResponse{
			Status: ds.Status{Message: ds.StatusInvalidToken},
		}
	}

	identities, exist, err := s.storageAuth.GetAuthIdentitiesByUserID(dbRT.UserID)
	if err != nil {
		s.logger.ErrorKV("Refresh.GetAuthIdentitiesByUserID", "error", err.Error())
		return nil
	}

	if !exist {
		s.logger.ErrorKV("Refresh.GetAuthIdentitiesByUserID",
			"error", fmt.Sprintf("no found identities for user %d", dbRT.UserID))
		return nil
	}

	if expired && dbRT.Used {
		err := s.storageAuth.DeleteAllUserSession(identities.UserID)
		if err != nil {
			s.logger.ErrorKV("Refresh.DeleteAllUserSession", "error", err.Error())
			return nil
		}
		return &ds.RefreshResponse{
			Status: ds.Status{Message: ds.StatusSessionReset},
		}
	}

	err = s.storageAuth.UpdateRefreshToken(&ds.DBUpdateRefreshToken{
		RefreshToken: dbRT.RefreshToken,
		ExpiresAt:    time.Now().Add(ds.DefaultRefreshTokenGracePeriod),
		Used:         true,
		Revoked:      false,
	})
	if err != nil {
		s.logger.ErrorKV("Refresh.UpdateRefreshToken", "error", err.Error())
		return nil
	}

	at, rt, err := s.createTokensPair(identities)
	if err != nil {
		s.logger.ErrorKV("Refresh.createTokensPair", "error", err.Error())
		return nil
	}

	return &ds.RefreshResponse{
		AccessToken:  ds.AccessToken{AccessToken: at},
		Status:       ds.Status{Message: ds.StatusSuccess},
		RefreshToken: ds.RefreshToken{RefreshToken: rt},
	}
}

func (s *Service) createTokensPair(identities *ds.AuthIdentities) (at string, rt string, err error) {
	rt = rand.Text()
	err = s.storageAuth.AddRefreshToken(&ds.DBRefreshToken{
		RefreshToken: ds.RefreshToken{RefreshToken: hashRefreshToken(rt)},
		ExpiresAt:    time.Now().Add(ds.DefaultRefreshTokenTTL),
		UserID:       identities.UserID,
		Revoked:      false,
		Used:         false,
	})
	if err != nil {
		return
	}

	now := time.Now()
	claim := ds.RegisteredClaims{
		UserId: identities.UserID,
		Role:   identities.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    serviceName,
			ExpiresAt: jwt.NewNumericDate(now.Add(ds.DefaultAccessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	sign, err := s.secret.ReadSecret(ds.JWTSecretKey)
	if err != nil {
		return
	}

	atRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	at, err = atRaw.SignedString([]byte(sign))
	if err != nil {
		return
	}

	return
}

func hashRefreshToken(t string) string {
	hash := sha256.Sum256([]byte(t))
	return hex.EncodeToString(hash[:])
}
