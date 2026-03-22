package mysql

import (
	"context"
	"database/sql"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddNewUser(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().CreateUserAuth(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

		req := &ds.DBRegisterRequest{}
		_, err := tc.c.AddNewUser(req)
		require.Nil(t, err)
	})

	t.Run("error on CreateUserAuth", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().CreateUserAuth(gomock.Any(), gomock.Any()).Return(tc.sqlRes, errTest)

		req := &ds.DBRegisterRequest{}
		_, err := tc.c.AddNewUser(req)
		require.NotNil(t, err)
	})

	t.Run("error on LastInsertId", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().CreateUserAuth(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(0), errTest)

		req := &ds.DBRegisterRequest{}
		_, err := tc.c.AddNewUser(req)
		require.NotNil(t, err)
	})

	t.Run("error on CreateUser", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().CreateUserAuth(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errTest)

		req := &ds.DBRegisterRequest{}
		_, err := tc.c.AddNewUser(req)
		require.NotNil(t, err)
	})

	t.Run("Ok duplicate", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().CreateUserAuth(gomock.Any(), gomock.Any()).Return(tc.sqlRes, errDuplicate)

		req := &ds.DBRegisterRequest{}
		res, err := tc.c.AddNewUser(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusUserAlreadyExists)
	})
}

func TestGetAuthIdentitiesByUserID(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().
			GetUserIdentitiesById(gomock.Any(), gomock.Any()).
			Return(sqlc.GetUserIdentitiesByIdRow{
				ID:           1,
				Login:        "login",
				PasswordHash: "pass",
				Name:         "name",
				Role:         "role",
			}, nil)

		res, ok, err := tc.c.GetAuthIdentitiesByUserID(1)
		require.Nil(t, err)
		require.True(t, ok)
		require.Equal(t, res.Login, "login")
		require.Equal(t, res.Password, "pass")
		require.Equal(t, res.Name, "name")
		require.Equal(t, res.Role, "role")
		require.Equal(t, res.UserID, int64(1))
	})

	t.Run("error on GetUserIdentitiesById", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().
			GetUserIdentitiesById(gomock.Any(), gomock.Any()).
			Return(sqlc.GetUserIdentitiesByIdRow{}, errTest)

		_, ok, err := tc.c.GetAuthIdentitiesByUserID(1)
		require.NotNil(t, err)
		require.False(t, ok)

	})

	t.Run("Ok no rows", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().
			GetUserIdentitiesById(gomock.Any(), gomock.Any()).
			Return(sqlc.GetUserIdentitiesByIdRow{}, errNoRow)

		_, ok, err := tc.c.GetAuthIdentitiesByUserID(1)
		require.Nil(t, err)
		require.False(t, ok)

	})
}

func TestGetAuthIdentitiesByLogin(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().
			GetUserIdentitiesByLogin(gomock.Any(), gomock.Any()).
			Return(sqlc.GetUserIdentitiesByLoginRow{
				ID:           1,
				Login:        "log",
				PasswordHash: "pass",
				Name:         "name",
				Role:         "role",
			}, nil)

		res, ok, err := tc.c.GetAuthIdentitiesByLogin("login")
		require.Nil(t, err)
		require.True(t, ok)
		require.Equal(t, res.Login, "log")
		require.Equal(t, res.Password, "pass")
		require.Equal(t, res.Name, "name")
		require.Equal(t, res.Role, "role")
		require.Equal(t, res.UserID, int64(1))
	})

	t.Run("error on GetUserIdentitiesByLogin", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().
			GetUserIdentitiesByLogin(gomock.Any(), gomock.Any()).
			Return(sqlc.GetUserIdentitiesByLoginRow{}, errTest)

		_, ok, err := tc.c.GetAuthIdentitiesByLogin("login")
		require.NotNil(t, err)
		require.False(t, ok)

	})

	t.Run("Ok no rows", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().
			GetUserIdentitiesByLogin(gomock.Any(), gomock.Any()).
			Return(sqlc.GetUserIdentitiesByLoginRow{}, errNoRow)

		_, ok, err := tc.c.GetAuthIdentitiesByLogin("login")
		require.Nil(t, err)
		require.False(t, ok)

	})
}

func TestGetRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tt := time.Now()
		tc.querier.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).
			Return(sqlc.RefreshToken{
				Token:     "tok",
				UserID:    1,
				ExpiredAt: tt,
				Revoked:   true,
				Used:      false,
			}, nil)

		res, ok, err := tc.c.GetRefreshToken("tok")
		require.Nil(t, err)
		require.True(t, ok)
		require.Equal(t, res.RefreshToken.RefreshToken, "tok")
		require.Equal(t, res.UserID, int64(1))
		require.Equal(t, res.ExpiresAt, tt)
		require.Equal(t, res.Revoked, true)
		require.Equal(t, res.Used, false)
	})

	t.Run("error on GetRefreshToken", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tt := time.Now()
		tc.querier.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).
			Return(sqlc.RefreshToken{
				Token:     "tok",
				UserID:    1,
				ExpiredAt: tt,
				Revoked:   true,
				Used:      false,
			}, errTest)

		_, ok, err := tc.c.GetRefreshToken("tok")
		require.NotNil(t, err)
		require.False(t, ok)
	})

	t.Run("Ok no rows", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).
			Return(sqlc.RefreshToken{}, errNoRow)

		_, ok, err := tc.c.GetRefreshToken("tok")
		require.Nil(t, err)
		require.False(t, ok)
	})
}

func TestUpdateRefreshToken(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().UpdateRefreshToken(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		req := &ds.DBUpdateRefreshToken{}
		err := tc.c.UpdateRefreshToken(req)
		require.Nil(t, err)
	})

	t.Run("error on UpdateRefreshToken", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().UpdateRefreshToken(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.DBUpdateRefreshToken{}
		err := tc.c.UpdateRefreshToken(req)
		require.NotNil(t, err)

	})

	t.Run("error on RowsAffected", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().UpdateRefreshToken(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), errTest)

		req := &ds.DBUpdateRefreshToken{}
		err := tc.c.UpdateRefreshToken(req)
		require.NotNil(t, err)
	})

	t.Run("error No rows affectet", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().UpdateRefreshToken(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.DBUpdateRefreshToken{}
		err := tc.c.UpdateRefreshToken(req)
		require.NotNil(t, err)
	})
}

func TestCleanupUslessTokens(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().CleanupUselessTokens(gomock.Any()).Return(nil)

		err := tc.c.CleanupUslessTokens()
		require.Nil(t, err)
	})

	t.Run("error on CleanupUslessTokens", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().CleanupUselessTokens(gomock.Any()).Return(errTest)

		err := tc.c.CleanupUslessTokens()
		require.NotNil(t, err)

	})
}

func TestDeleteAllUserSession(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().DeleteUserSessions(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		err := tc.c.DeleteAllUserSession(21)
		require.Nil(t, err)
	})

	t.Run("error on DeleteUserSessions", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().DeleteUserSessions(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		err := tc.c.DeleteAllUserSession(1)
		require.NotNil(t, err)

	})

	t.Run("error on RowsAffected", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().DeleteUserSessions(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), errTest)

		err := tc.c.DeleteAllUserSession(1)
		require.NotNil(t, err)
	})

	t.Run("error No rows affectet", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)
		tc.querier.EXPECT().DeleteUserSessions(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		err := tc.c.DeleteAllUserSession(1)
		require.NotNil(t, err)
	})
}
