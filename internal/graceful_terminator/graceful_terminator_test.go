package gracefulterminator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGracefulTerminator(t *testing.T) {
	t.Parallel()

	arr := []int{}
	Add(func() {
		arr = append(arr, 1)
	})

	Add(func() {
		arr = append(arr, 2)
	})

	Add(func() {
		arr = append(arr, 3)
	})

	Stop()

	require.Equal(t, len(arr), int(3))
	require.Equal(t, arr[0], int(3))
	require.Equal(t, arr[1], int(2))
	require.Equal(t, arr[2], int(1))
}
