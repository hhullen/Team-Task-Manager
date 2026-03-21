package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"team-task-manager/internal/api/v1"
	gracefulterminator "team-task-manager/internal/graceful_terminator"
	"team-task-manager/internal/logger"
	"team-task-manager/internal/migrator"
	ratelimiter "team-task-manager/internal/redis_rate_limiter"
	secretprovider "team-task-manager/internal/secret_provider"
	"team-task-manager/internal/service"
	"testing"
	"time"

	mysqlclient "team-task-manager/internal/clients/mysql"
	redisclient "team-task-manager/internal/clients/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	apiAddress = ":8080"
	apiPrefix  = "/api/v1"
)

type ServicesTestSuite struct {
	suite.Suite
	ctx            context.Context
	cancelCtx      context.CancelFunc
	mysqlContainer *mysql.MySQLContainer
	redisContainer testcontainers.Container

	db    *sql.DB
	cache *redis.Client
	api   *api.API
}

func (s *ServicesTestSuite) SetupSuite() {
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	go func() {
		<-s.ctx.Done()
		gracefulterminator.Stop()
	}()

	dbPass := "root_default"
	dbName := "db_default"

	initScriptPath, _ := filepath.Abs("../../init/init.sh")
	initRolesPath, _ := filepath.Abs("../../init/init_roles.sql")

	mysqlC, err := mysql.Run(s.ctx,
		"mysql:8.0.45-debian",
		mysql.WithDatabase(dbName),
		mysql.WithPassword(dbPass),

		mysql.WithScripts(initScriptPath),
		testcontainers.WithFiles(
			testcontainers.ContainerFile{
				HostFilePath:      initRolesPath,
				ContainerFilePath: "/sql_init/init_roles.sql",
				FileMode:          0755,
			},
			testcontainers.ContainerFile{
				Reader:            strings.NewReader("root_default"),
				ContainerFilePath: "/run/secrets/db_root_password",
				FileMode:          0444,
			},
			testcontainers.ContainerFile{
				Reader:            strings.NewReader("app_default"),
				ContainerFilePath: "/run/secrets/db_app_user",
				FileMode:          0444,
			},
			testcontainers.ContainerFile{
				Reader:            strings.NewReader("app_default"),
				ContainerFilePath: "/run/secrets/db_app_password",
				FileMode:          0444,
			},
			testcontainers.ContainerFile{
				Reader:            strings.NewReader("migrator_default"),
				ContainerFilePath: "/run/secrets/db_migrator_password",
				FileMode:          0444,
			},
			testcontainers.ContainerFile{
				Reader:            strings.NewReader("migrator_default"),
				ContainerFilePath: "/run/secrets/db_migrator_user",
				FileMode:          0444,
			},
			testcontainers.ContainerFile{
				Reader:            strings.NewReader("db_default"),
				ContainerFilePath: "/run/secrets/db_name",
				FileMode:          0444,
			},
		),
	)
	s.Require().NoError(err, "failed running MySQL")
	s.mysqlContainer = mysqlC
	gracefulterminator.Add(func() {
		s.mysqlContainer.Terminate(s.ctx)
	})

	dbConnStr, err := mysqlC.ConnectionString(s.ctx)
	s.Require().NoError(err)

	s.db, err = sql.Open("mysql", dbConnStr+"?parseTime=true")
	s.Require().NoError(err)
	gracefulterminator.Add(func() {
		s.db.Close()
	})

	redisPass := "redis_default"

	redisReq := testcontainers.ContainerRequest{
		Image:        "redis:8.2.4-alpine3.22",
		ExposedPorts: []string{"6379/tcp"},
		Cmd:          []string{"redis-server", "--requirepass", redisPass, "--appendonly", "yes"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(15 * time.Second),
	}

	redisC, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: redisReq,
		Started:          true,
	})
	s.Require().NoError(err, "failed running Redis")
	s.redisContainer = redisC
	gracefulterminator.Add(func() {
		s.redisContainer.Terminate(s.ctx)
	})

	redisHost, err := redisC.Host(s.ctx)
	s.Require().NoError(err)
	redisPort, err := redisC.MappedPort(s.ctx, "6379")
	s.Require().NoError(err)

	s.cache = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort.Port()),
		Password: redisPass,
	})
	gracefulterminator.Add(func() {
		s.cache.Close()
	})

	s.setupAPI()
}

func (s *ServicesTestSuite) setupAPI() {
	testLog := logger.NewLogger(os.Stdout, "TEST")
	gracefulterminator.Add(func() {
		testLog.Stop()
	})

	secret := secretprovider.NewSecretProvider("../../secrets/")
	dbClient := mysqlclient.NewClient(s.ctx, s.db, testLog)
	redisClient := redisclient.NewClient(s.cache)

	service := service.NewService(s.ctx,
		dbClient,
		dbClient,
		testLog,
		redisClient,
		secret,
	)

	limiter := ratelimiter.NewRedisRateLimiter(s.ctx, s.cache)

	var err error
	s.api, err = api.NewAPI(s.ctx, apiAddress, service, service, secret, limiter, testLog)
	s.Require().NoError(err)
	gracefulterminator.Add(func() {
		s.api.Stop()
	})

	migrator.ExecMigration(s.db, "up", "../../migrations/mysql")

	go s.api.StartListening()
}

func (s *ServicesTestSuite) TearDownSuite() {
	s.cancelCtx()
}

func (s *ServicesTestSuite) TestDatabaseConnection() {
	err := s.db.Ping()
	s.NoError(err)
}

func (s *ServicesTestSuite) TestRedisConnection() {
	pong, err := s.cache.Ping(context.Background()).Result()
	s.NoError(err)
	s.Equal("PONG", pong)
}

func TestServicesSuite(t *testing.T) {
	suite.Run(t, new(ServicesTestSuite))
}

func (s *ServicesTestSuite) JSONBodyRequest(method string, payload map[string]any, uri string, header [][2]string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(method, uri, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	for _, v := range header {
		req.Header.Set(v[0], v[1])
	}

	w := httptest.NewRecorder()

	s.api.ServeHTTP(w, req)

	return w
}

func (s *ServicesTestSuite) QueryRequest(method string, payloadQuery map[string]string, payloadBody map[string]any, uri string, header [][2]string) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payloadBody)

	req := httptest.NewRequest(method, uri, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)

	q := req.URL.Query()
	for k, v := range payloadQuery {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for _, v := range header {
		req.Header.Set(v[0], v[1])
	}

	w := httptest.NewRecorder()

	s.api.ServeHTTP(w, req)

	return w
}

func (s *ServicesTestSuite) register(name string) *httptest.ResponseRecorder {
	payload := map[string]any{
		"login":    name,
		"name":     name,
		"password": name,
	}

	return s.JSONBodyRequest(http.MethodPost, payload, apiPrefix+"/register", [][2]string{})
}

func (s *ServicesTestSuite) login(name string) (at, rt string) {
	payload := map[string]any{
		"login":    name,
		"password": name,
	}

	w := s.JSONBodyRequest(http.MethodPost, payload, apiPrefix+"/login", [][2]string{})
	v := map[string]string{}
	err := json.NewDecoder(w.Body).Decode(&v)
	s.Nil(err)

	at, ok := v["access_token"]
	s.True(ok)
	s.True(at != "")

	res := w.Result()
	var coockie *http.Cookie
	for _, v := range res.Cookies() {
		if v.Name == "refresh_token" {
			coockie = v
		}
	}
	s.NotNil(coockie)

	return at, coockie.Value
}

func (s *ServicesTestSuite) createTeam(name, at string) *httptest.ResponseRecorder {
	payload := map[string]any{
		"name":        name,
		"description": name,
	}

	uri := apiPrefix + "/teams"

	w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
		{"Authorization", "Bearer " + at},
	})

	s.Equal(http.StatusOK, w.Code)

	return w
}

func (s *ServicesTestSuite) createTaks(name, teamName, at string) *httptest.ResponseRecorder {
	uri := apiPrefix + "/tasks"

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", teamName).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	payload := map[string]any{
		"assignee_login": name,
		"subject":        name,
		"description":    name,
		"status":         "todo",
		"team_id":        teamId,
	}

	w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
		{"Authorization", "Bearer " + at},
	})
	s.Equal(http.StatusOK, w.Code)

	return w
}
