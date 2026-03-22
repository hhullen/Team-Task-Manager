package metrics

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGracefulTerminator(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		err := ReportResponse("GET", "/test", 200, 120)
		require.Nil(t, err)
	})

	require.NotPanics(t, func() {
		RepordDBStats(sql.DBStats{})
	})
}
