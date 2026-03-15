package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	ds "team-task-manager/internal/datastruct"
	_ "team-task-manager/internal/docs"
	"team-task-manager/internal/supports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/schema"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	address      = ":8080"
	readTimeout  = time.Second * 5
	writeTimeout = time.Second * 5

	globalRateLimit      = 500
	globalLimiterKey     = "global"
	overlimitMessage     = "service exhausted"
	userIdLimitPerSecond = 100
	limiterUserIdPrefix  = "user_id"

	apiPrefix     = "/api/v1"
	swaggerPrefix = "/swagger/"

	contentTypeKey = "Content-Type"
	contentLenKey  = "Content-Length"

	appJSONValue = "application/json"
)

var schemaDecoder = schema.NewDecoder()

type IWithStatus interface {
	GetStatus() string
}

type IWithUserId interface {
	SetUserId(int64)
}

type IWithUserRole interface {
	SetUserRole(string)
}

type IWithJWTUserCreds interface {
	IWithUserId
	IWithUserRole
}

type IAppService interface {
	CreateTeam(*ds.CreateTeamRequest) *ds.CreateTeamResponse
	ListUserTeams(*ds.ListUserTeamsRequest) *ds.ListUserTeamsResponse
	InviteUserToTeam(*ds.InviteUserToTeamRequest) *ds.InviteUserToTeamResponse

	AddNewTask(*ds.CreateTaskRequest) *ds.CreateTaskResponse
	GetTasks(*ds.GetTasksRequest) *ds.GetTasksResponse
}

type IAuthService interface {
	RegisterUser(req *ds.RegisterRequest) *ds.RegisterResponse
	LoginUser(req *ds.LoginRequest) *ds.LoginResponse
	Refresh(req *ds.RefreshRequest) *ds.RefreshResponse
}

type ISecretProvider interface {
	ReadSecret(key string) (string, error)
}

type IServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type IRouter interface {
	Handle(pattern string, handler http.Handler)
}

type ILogger interface {
	InfoKV(message string, argsKV ...any)
	WarnKV(message string, argsKV ...any)
	ErrorKV(message string, argsKV ...any)
	FatalKV(message string, argsKV ...any)
}

type ILimiter interface {
	Allow(key string, perSecond int) (bool, error)
}

type ExecArgs[ReqT any, RespT IWithStatus] struct {
	api              *API
	serviceFunc      func(*ReqT) *RespT
	requestExtractor func(r *http.Request, v *ReqT) error
	responseWriter   func(w *http.ResponseWriter, v *RespT) error
	httpRequest      *http.Request
	httpResponse     *http.ResponseWriter
	validator        func(s *ReqT) error
}

type ResponseWriterInterceptor struct {
	rw   http.ResponseWriter
	code int
}

func (rwi *ResponseWriterInterceptor) Header() http.Header {
	return rwi.rw.Header()
}

func (rwi *ResponseWriterInterceptor) Write(data []byte) (int, error) {
	return rwi.rw.Write(data)
}

func (rwi *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	rwi.code = statusCode
	rwi.rw.WriteHeader(statusCode)
}

type API struct {
	ctx         context.Context
	logger      ILogger
	appService  IAppService
	authService IAuthService
	secret      ISecretProvider
	limiter     ILimiter
	server      IServer
	router      IRouter

	jwtSecret string
}

func NewAPI(ctx context.Context,
	app IAppService,
	auth IAuthService,
	sec ISecretProvider,
	lim ILimiter,
	log ILogger) (*API, error) {
	router := http.NewServeMux()

	router.Handle(swaggerPrefix, httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         address,
		Handler:      mainMiddleware(router, log),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	jwtSec, err := sec.ReadSecret(ds.JWTSecretKey)
	if err != nil {
		return nil, err
	}

	return buildAPI(ctx, app, auth, log, sec, lim, server, router, jwtSec), nil
}

func buildAPI(ctx context.Context,
	app IAppService,
	auth IAuthService,
	log ILogger,
	sec ISecretProvider,
	lim ILimiter,
	srv IServer,
	rou IRouter,
	jwtSec string) *API {
	api := &API{
		ctx:         ctx,
		logger:      log,
		appService:  app,
		authService: auth,
		secret:      sec,
		limiter:     lim,
		server:      srv,
		router:      rou,
		jwtSecret:   jwtSec,
	}

	api.setupAuthHandlers()
	api.setupTeamsHandlers()
	api.setupTasksHandlers()

	return api
}

func (a *API) StartListening() error {
	a.logger.InfoKV("Server is listening")
	return a.server.ListenAndServe()
}

func (a *API) Stop() error {
	a.logger.InfoKV("Server Shutdown")
	return a.server.Shutdown(a.ctx)
}

func mainMiddleware(next http.Handler, log ILogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now()

		rwi := &ResponseWriterInterceptor{rw: w}
		next.ServeHTTP(rwi, r)

		te := time.Since(ts)
		log.InfoKV("Request", "method", r.Method, "url", r.URL.String(), "duration(ms)", te.Milliseconds(), "status_code", rwi.code)
	})
}

func globalRateLimitedMiddleware(a *API, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow, err := a.limiter.Allow(globalLimiterKey, globalRateLimit)
		if err != nil {
			http.Error(w, ds.StatusServiceError, http.StatusInternalServerError)
			return
		}

		if !allow {
			http.Error(w, overlimitMessage, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jwtBasedMiddleware(a *API, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)

		claims, err := a.validateJWT(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			a.logger.WarnKV("failed token validation", "error", err.Error())
			return
		}

		key := supports.Concat(strconv.FormatInt(claims.UserId, 10), ":", limiterUserIdPrefix)

		allow, err := a.limiter.Allow(key, userIdLimitPerSecond)
		if err != nil {
			http.Error(w, ds.StatusServiceError, http.StatusInternalServerError)
			a.logger.ErrorKV("failed asking rate limiter", "error", err.Error())
			return
		}

		if !allow {
			http.Error(w, overlimitMessage, http.StatusTooManyRequests)
			a.logger.WarnKV(overlimitMessage, "user_id", claims.UserId)
			return
		}

		ctx := context.WithValue(r.Context(), ds.UserIDKey, claims.UserId)
		ctx = context.WithValue(ctx, ds.UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) validateJWT(rawToken string) (*ds.RegisteredClaims, error) {
	claims := &ds.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func pattern(method, prefixPath string) string {
	return supports.Concat(method, " ", prefixPath)
}
