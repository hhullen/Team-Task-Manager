package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	errTest       = errors.New("test")
	errDuplicate  = &mysql.MySQLError{Number: mysqlerr.ER_DUP_ENTRY}
	errlongData   = &mysql.MySQLError{Number: mysqlerr.ER_DATA_TOO_LONG}
	errForeignKey = &mysql.MySQLError{Number: mysqlerr.ER_NO_REFERENCED_ROW_2}
	errNoRow      = sql.ErrNoRows

	patchTest = `{"status": "@@ -0,0 +1,4 @@\n+todo\n", "subject": "@@ -0,0 +1,16 @@\n+service endpoint\n", "team_id": "@@ -0,0 +1 @@\n+1\n", "assignee_id": "@@ -0,0 +1 @@\n+4\n", "description": "@@ -0,0 +1,24 @@\n+add new service endpoint\n"}`
)

type TestClient struct {
	ctx     context.Context
	querier *MockIQuerier
	db      *MockIDB
	log     *MockILogger
	c       *Client
	sqlRes  *MockResult
}

func newTestClient(t *testing.T) *TestClient {
	mc := gomock.NewController(t)
	tc := &TestClient{
		ctx:     context.Background(),
		querier: NewMockIQuerier(mc),
		db:      NewMockIDB(mc),
		log:     NewMockILogger(mc),
		sqlRes:  NewMockResult(mc),
	}

	c := buildClient(tc.db)

	tc.c = c

	return tc
}

func TestMySQLErrors(t *testing.T) {
	t.Parallel()

	t.Run("isDuplicate", func(t *testing.T) {
		t.Parallel()
		require.True(t, isDuplicate(errDuplicate))
		require.False(t, isDuplicate(errTest))
	})

	t.Run("isLongData", func(t *testing.T) {
		t.Parallel()
		require.True(t, isLongData(errlongData))
		require.False(t, isLongData(errTest))
	})

	t.Run("isNoRows", func(t *testing.T) {
		t.Parallel()
		require.True(t, isNoRows(errNoRow))
		require.False(t, isNoRows(errTest))
	})

	t.Run("isForeignKeyErr", func(t *testing.T) {
		t.Parallel()
		require.True(t, isForeignKeyErr(errForeignKey))
		require.False(t, isForeignKeyErr(errTest))
	})
}

func TestNullTypes(t *testing.T) {
	t.Parallel()

	t.Run("nullString", func(t *testing.T) {
		t.Parallel()
		s := nullString(nil)
		require.False(t, s.Valid)

		ts := "dsf"

		s = nullString(&ts)
		require.True(t, s.Valid)
	})

	t.Run("fromNullString", func(t *testing.T) {
		t.Parallel()
		ns := sql.NullString{Valid: false}
		s := fromNullString(ns)
		require.Nil(t, s)

		ns = sql.NullString{Valid: true, String: "ead"}
		s = fromNullString(ns)
		require.NotNil(t, s)
	})

	t.Run("nullInt64", func(t *testing.T) {
		t.Parallel()
		s := nullInt64(nil)
		require.False(t, s.Valid)

		ts := int64(123)

		s = nullInt64(&ts)
		require.True(t, s.Valid)
	})
}
