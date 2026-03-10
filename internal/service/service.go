package service

import (
	"context"
	ds "team-task-manager/internal/datastruct"
	"time"
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

type IStorage interface {
	AddNewUser(*ds.RegisterRequest) (*ds.RegisterResponse, error)
	GetAuthIdentitiesByLogin(login string) (*ds.DBAuthIdentities, error)
	GetAuthIdentitiesByUserID(id int64) (*ds.DBAuthIdentities, error)
	AddRefreshToken(ident *ds.DBAuthIdentities, toket string, expire time.Time) error
	GetRefreshToken(token string) (*ds.DBRefreshToken, bool, error)
	DeleteRefreshToken(token string) error
}

type ISecretProvider interface {
	ReadSecret(key string) (string, error)
}

type Service struct {
	ctx     context.Context
	logger  ILogger
	cache   ICache
	storage IStorage
	secret  ISecretProvider
}

func NewService(ctx context.Context, s IStorage, l ILogger, c ICache, sp ISecretProvider) *Service {
	return buildService(ctx, s, l, c, sp)
}

func buildService(ctx context.Context, s IStorage, l ILogger, c ICache, sp ISecretProvider) *Service {
	return &Service{
		ctx:     ctx,
		logger:  l,
		cache:   c,
		storage: s,
		secret:  sp,
	}
}
