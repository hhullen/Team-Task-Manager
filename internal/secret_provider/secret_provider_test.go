package secretprovider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testKey = "test_value"
)

func TestSecretProvider(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, testKey)

	err := os.WriteFile(filePath, []byte(testKey), 0644)
	require.NoError(t, err)

	t.Run("Ok", func(t *testing.T) {
		sp := NewSecretProvider(tempDir)
		v, err := sp.ReadSecret(testKey)
		require.Nil(t, err)
		require.Equal(t, v, testKey)

		v, err = sp.ReadSecret(testKey)
		require.Nil(t, err)
		require.Equal(t, v, testKey)
	})

	t.Run("Not existed", func(t *testing.T) {
		sp := NewSecretProvider(tempDir)
		_, err := sp.ReadSecret("not_exists")
		require.NotNil(t, err)
	})

}
