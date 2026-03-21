package service

import (
	"context"
	"errors"
	"strings"
	ds "team-task-manager/internal/datastruct"
	gracefulterminator "team-task-manager/internal/graceful_terminator"
	"team-task-manager/internal/supports"

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

type IAppStorage interface {
	AddNewTeam(*ds.CreateTeamRequest) (*ds.CreateTeamResponse, error)
	GetUserTeams(userId int64) (*ds.ListUserTeamsResponse, error)
	AddUserToUserTeam(*ds.DBInviteUserToTeamRequest) (*ds.InviteUserToTeamResponse, error)

	AddNewTask(*ds.DBCreateTaskRequest) (*ds.CreateTaskResponse, error)
	GetTasks(*ds.GetTasksRequest) (*ds.GetTasksResponse, error)
	UpdateTask(req *ds.DBUpdateTaskRequest) (*ds.UpdateTaskResponse, error)
	GetTaskHistory(*ds.GetTaskHistoryRequest) (*ds.GetTaskHistoryResponse, error)
	AddTaskComment(req *ds.AddTaskCommentRequest) (*ds.AddTaskCommentResponse, error)
}

type ISecretProvider interface {
	ReadSecret(key string) (string, error)
}

type Service struct {
	ctx         context.Context
	logger      ILogger
	cache       ICache
	storageAuth IAuthStorage
	storageApp  IAppStorage
	secret      ISecretProvider
	cron        *cron.Cron
}

func NewService(ctx context.Context, aus IAuthStorage, aps IAppStorage, l ILogger, c ICache, sp ISecretProvider) *Service {
	s := buildService(ctx, aus, aps, l, c, sp)

	s.cron = cron.New()
	_, err := s.cron.AddFunc("0 0 * * *", func() {
		err := s.storageAuth.CleanupUslessTokens()
		if err != nil {
			s.logger.ErrorKV("CleanupUslessTokens", "error", err)
		}
	})
	if err != nil {
		s.logger.FatalKV("NewService.AddFunc", "error", err)
	} else {
		s.cron.Start()
	}

	gracefulterminator.Add(func() {
		<-s.cron.Stop().Done()
	})

	return s
}

func buildService(ctx context.Context, aus IAuthStorage, aps IAppStorage, l ILogger, c ICache, sp ISecretProvider) *Service {
	return &Service{
		ctx:         ctx,
		logger:      l,
		cache:       c,
		storageAuth: aus,
		storageApp:  aps,
		secret:      sp,
	}
}

func makeCacheKey(vv ...string) string {
	length := 0
	for i := range vv {
		length += len(vv[i])
	}

	var b strings.Builder
	b.Grow(length)

	for i := range vv {
		b.WriteString(vv[i])
		b.WriteByte('_')
	}

	return b.String()
}

func execWithCache[RespT ICachedState](s *Service, key string, avoidCache bool, fetch func() (RespT, error)) (RespT, error) {
	var response RespT
	var cached bool
	var err error

	if !avoidCache {
		cached, err = s.cache.Read(key, &response)
		if err != nil {
			s.logger.ErrorKV("failed reading cache", "message", err.Error())
		}
	}

	if cached {
		response.SetCached(true)
		return response, nil
	}

	response, err = fetch()
	if err != nil {
		var empty RespT
		return empty, err
	}

	err = s.cache.Write(key, response)
	if err != nil {
		s.logger.ErrorKV("failed writing cache", "message", err.Error())
	}
	response.SetCached(false)

	return response, nil
}

func (s *Service) getAuthIdentitiesByLogin(login string, avoidCache bool) (*ds.AuthIdentities, bool, error) {
	notExist := errors.New("not exists")
	ident, err := execWithCache(s,
		makeCacheKey("IdByLogin", supports.FNV1Hash([]byte(login))),
		avoidCache,
		func() (*ds.AuthIdentities, error) {
			ident, ok, err := s.storageAuth.GetAuthIdentitiesByLogin(login)
			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, notExist
			}
			return ident, err
		})

	if err != nil {
		if errors.Is(err, notExist) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return ident, true, nil
}
