package service

import (
	"context"
	"errors"
	ds "team-task-manager/internal/datastruct"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	errTest = errors.New("error")
)

type Caching struct {
	ds.CachedStatus
	v int
}

type ServiceTest struct {
	ctx       context.Context
	logger    *MockILogger
	cache     *MockICache
	authStore *MockIAuthStorage
	appStore  *MockIAppStorage
	secret    *MockISecretProvider
	s         *Service
}

func newServiceTest(t *testing.T) *ServiceTest {
	mc := gomock.NewController(t)

	st := &ServiceTest{
		ctx:       context.Background(),
		logger:    NewMockILogger(mc),
		cache:     NewMockICache(mc),
		authStore: NewMockIAuthStorage(mc),
		appStore:  NewMockIAppStorage(mc),
		secret:    NewMockISecretProvider(mc),
	}

	s := buildService(st.ctx, st.authStore, st.appStore, st.logger, st.cache, st.secret)

	st.s = s

	return st
}

func TestMakeCacheKey(t *testing.T) {
	t.Parallel()

	res := makeCacheKey("key1", "key2")
	require.Equal(t, res, "key1_key2_")
}

func TestExecWithCache(t *testing.T) {
	t.Parallel()

	t.Run("Ok from cache", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				vv := v.(**Caching)
				*vv = &Caching{v: 2}
				return true, nil
			})

		res, err := execWithCache(st.s, "key", false, func() (*Caching, error) {
			return &Caching{}, nil
		})

		require.Nil(t, err)
		require.Equal(t, res.v, int(2))
		require.True(t, res.CachedStatus.Cached)
	})

	t.Run("Ok fetched", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})

		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		res, err := execWithCache(st.s, "key", false, func() (*Caching, error) {
			return &Caching{v: 2}, nil
		})

		require.Nil(t, err)
		require.Equal(t, res.v, int(2))
		require.False(t, res.CachedStatus.Cached)
	})

	t.Run("error on Read", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, errTest
			})
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		res, err := execWithCache(st.s, "key", false, func() (*Caching, error) {
			return &Caching{v: 2}, nil
		})

		require.Nil(t, err)
		require.Equal(t, res.v, int(2))
		require.False(t, res.CachedStatus.Cached)
	})

	t.Run("error on Read and Write", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, errTest
			})
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		res, err := execWithCache(st.s, "key", false, func() (*Caching, error) {
			return &Caching{v: 2}, nil
		})

		require.Nil(t, err)
		require.Equal(t, res.v, int(2))
		require.False(t, res.CachedStatus.Cached)
	})

	t.Run("error on Fetch", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})

		_, err := execWithCache(st.s, "key", false, func() (*Caching, error) {
			return nil, errTest
		})

		require.NotNil(t, err)
	})
}

func TestGetAuthIdentitiesByLogin(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		res, ok, err := st.s.getAuthIdentitiesByLogin("login", false)

		require.Nil(t, err)
		require.False(t, res.CachedStatus.Cached)
		require.True(t, ok)
	})

	t.Run("Not exists", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(nil, false, nil)

		_, ok, err := st.s.getAuthIdentitiesByLogin("login", false)

		require.Nil(t, err)
		require.False(t, ok)
	})

	t.Run("error on GetAuthIdentitiesByLogin", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(nil, false, errTest)

		_, ok, err := st.s.getAuthIdentitiesByLogin("login", false)

		require.NotNil(t, err)
		require.False(t, ok)
	})
}
