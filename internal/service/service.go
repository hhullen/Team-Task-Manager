package service

import (
	"context"
	ds "team-task-manager/internal/datastruct"

	"github.com/robfig/cron/v3"
)

const (
	serviceName = "team-tast-manager"
)

type ILogger interface {
	InfoKV(message string, argsKV ...any)
	WarnKV(message string, argsKV ...any)
	ErrorKV(message string, argsKV ...any)
	FatalKV(message string, argsKV ...any)
}

type ICache interface {
	Read(key string, v any) (bool, error)
	Write(key string, v any) error
}

type ICachedState interface {
	SetCached(bool)
}

type IAuthStorage interface {
	AddNewUser(*ds.DBRegisterRequest) (*ds.RegisterResponse, error)
	GetAuthIdentitiesByUserID(id int64) (*ds.AuthIdentities, bool, error)
	GetAuthIdentitiesByLogin(login string) (*ds.AuthIdentities, bool, error)
	AddRefreshToken(*ds.DBRefreshToken) error
	GetRefreshToken(token string) (*ds.DBRefreshToken, bool, error)
	UpdateRefreshToken(*ds.DBUpdateRefreshToken) error
	CleanupUslessTokens() error
	DeleteAllUserSession(userId int64) error
}

type ISecretProvider interface {
	ReadSecret(key string) (string, error)
}

type Service struct {
	ctx         context.Context
	logger      ILogger
	cache       ICache
	storageAuth IAuthStorage
	secret      ISecretProvider
	cron        *cron.Cron
}

func NewService(ctx context.Context, as IAuthStorage, l ILogger, c ICache, sp ISecretProvider) *Service {
	s := buildService(ctx, as, l, c, sp)

	s.cron = cron.New()
	s.cron.AddFunc("0 0 * * *", func() {
		err := s.storageAuth.CleanupUslessTokens()
		if err != nil {
			s.logger.ErrorKV("CleanupUslessTokens", "error", err)
		}
	})

	return s
}

func buildService(ctx context.Context, as IAuthStorage, l ILogger, c ICache, sp ISecretProvider) *Service {
	return &Service{
		ctx:         ctx,
		logger:      l,
		cache:       c,
		storageAuth: as,
		secret:      sp,
	}
}
