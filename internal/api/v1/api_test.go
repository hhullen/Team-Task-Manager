package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ds "team-task-manager/internal/datastruct"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	errTest       = errors.New("test")
	jwtSecretTest = "test"
	userIdTest    = int64(16)
	userRoleTest  = "testRole"
)

type TestAPI struct {
	ctx         context.Context
	logger      *MockILogger
	appService  *MockIAppService
	authService *MockIAuthService
	secret      *MockISecretProvider
	limiter     *MockILimiter
	server      *MockIServer
	router      *MockIRouter
	handler     *MockHandler
	accessToken string

	a *API
}

func newTestAPI(t *testing.T) *TestAPI {
	now := time.Now()
	claim := ds.RegisteredClaims{
		UserId: userIdTest,
		Role:   userRoleTest,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "tester",
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 5)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	atRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	at, err := atRaw.SignedString([]byte(jwtSecretTest))
	require.Nil(t, err)

	mc := gomock.NewController(t)
	ta := &TestAPI{
		ctx:         context.Background(),
		logger:      NewMockILogger(mc),
		appService:  NewMockIAppService(mc),
		authService: NewMockIAuthService(mc),
		secret:      NewMockISecretProvider(mc),
		limiter:     NewMockILimiter(mc),
		server:      NewMockIServer(mc),
		router:      NewMockIRouter(mc),
		handler:     NewMockHandler(mc),
		accessToken: at,
	}

	ta.router.EXPECT().Handle(gomock.Any(), gomock.Any()).MinTimes(1)

	a := buildAPI(ta.ctx, ta.appService, ta.authService, ta.logger, ta.secret, ta.limiter, ta.server, ta.router, jwtSecretTest)

	ta.a = a

	return ta
}

func TestStartListening(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.server.EXPECT().ListenAndServe().Return(nil)
		ta.logger.EXPECT().InfoKV(gomock.Any(), gomock.All())

		err := ta.a.StartListening()
		require.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.server.EXPECT().ListenAndServe().Return(errTest)
		ta.logger.EXPECT().InfoKV(gomock.Any(), gomock.All())

		err := ta.a.StartListening()
		require.NotNil(t, err)
	})
}

func TestStop(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.logger.EXPECT().InfoKV(gomock.Any(), gomock.All())
		ta.server.EXPECT().Shutdown(gomock.Any()).Return(nil)

		err := ta.a.Stop()
		require.Nil(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.logger.EXPECT().InfoKV(gomock.Any(), gomock.All())
		ta.server.EXPECT().Shutdown(gomock.Any()).Return(errTest)

		err := ta.a.Stop()
		require.NotNil(t, err)
	})
}

func TestServeHTTP(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.router.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

		ta.a.ServeHTTP(httptest.NewRecorder(), &http.Request{})

	})

}

func TestMainMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.handler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
		ta.logger.EXPECT().InfoKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		w := httptest.NewRecorder()
		mainMiddleware(ta.handler, ta.logger).ServeHTTP(w, r)
	})
}

func TestGlobalRateLimitedMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.limiter.EXPECT().Allow(gomock.Any(), gomock.Any()).Return(true, nil)
		ta.handler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		w := httptest.NewRecorder()
		globalRateLimitedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)
	})

	t.Run("Error on limiter", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.limiter.EXPECT().Allow(gomock.Any(), gomock.Any()).Return(false, errTest)

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		w := httptest.NewRecorder()
		globalRateLimitedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusInternalServerError)
	})

	t.Run("Not allowed", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.limiter.EXPECT().Allow(gomock.Any(), gomock.Any()).Return(false, nil)

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		w := httptest.NewRecorder()
		globalRateLimitedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusTooManyRequests)
	})
}

func TestJwtBasedMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.limiter.EXPECT().Allow(gomock.Any(), gomock.Any()).Return(true, nil)
		ta.handler.EXPECT().ServeHTTP(gomock.Any(), gomock.Any()).Do(func(_ http.ResponseWriter, r *http.Request) {
			require.NotPanics(t, func() {
				require.Equal(t, r.Context().Value(ds.UserIDKey).(int64), userIdTest)
				require.Equal(t, r.Context().Value(ds.UserRoleKey).(string), userRoleTest)
			})
		})

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		r.Header.Add("Authorization", "Bearer "+ta.accessToken)
		w := httptest.NewRecorder()
		jwtBasedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

	})

	t.Run("No authorization header", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		w := httptest.NewRecorder()
		jwtBasedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusUnauthorized)
	})

	t.Run("Wrong JWT", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.accessToken = "Bearer Shit"

		ta.logger.EXPECT().WarnKV(gomock.Any(), gomock.All())
		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		r.Header.Add("Authorization", "Bearer "+ta.accessToken)
		w := httptest.NewRecorder()
		jwtBasedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusUnauthorized)
	})

	t.Run("Error on limiter", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.limiter.EXPECT().Allow(gomock.Any(), gomock.Any()).Return(false, errTest)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		r.Header.Add("Authorization", "Bearer "+ta.accessToken)
		w := httptest.NewRecorder()
		jwtBasedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusInternalServerError)
	})

	t.Run("Not allowed", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		ta.limiter.EXPECT().Allow(gomock.Any(), gomock.Any()).Return(false, nil)
		ta.logger.EXPECT().WarnKV(gomock.Any(), gomock.All())

		r := httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("test"))
		r.Header.Add("Authorization", "Bearer "+ta.accessToken)
		w := httptest.NewRecorder()
		jwtBasedMiddleware(ta.a, ta.handler).ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusTooManyRequests)
	})
}
