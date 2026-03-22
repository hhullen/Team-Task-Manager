package ratelimiter

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis_rate/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	mc := gomock.NewController(t)
	lim := NewMockILimiter(mc)
	rrl := buildRedisRateLimiter(context.Background(), lim)

	t.Run("Ok", func(t *testing.T) {
		lim.EXPECT().Allow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&redis_rate.Result{Allowed: 1}, nil)

		ok, err := rrl.Allow("key", 1)
		require.Nil(t, err)
		require.True(t, ok)
	})

	t.Run("Not allowed", func(t *testing.T) {
		lim.EXPECT().Allow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&redis_rate.Result{Allowed: 0}, nil)

		ok, err := rrl.Allow("key", 1)
		require.Nil(t, err)
		require.False(t, ok)
	})

	t.Run("error on Allow", func(t *testing.T) {
		lim.EXPECT().Allow(gomock.Any(), gomock.Any(), gomock.Any()).Return(&redis_rate.Result{}, errors.New("error"))

		ok, err := rrl.Allow("key", 1)
		require.NotNil(t, err)
		require.False(t, ok)
	})
}
