package mysql

import (
	"context"
	"database/sql"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddNewTeam(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		req := &ds.CreateTeamRequest{Name: "name"}
		res, err := tc.c.AddNewTeam(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on AddNewTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.CreateTeamRequest{Name: "name"}
		_, err := tc.c.AddNewTeam(req)
		require.NotNil(t, err)
	})

	t.Run("Foreign key on AddNewTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errForeignKey)

		req := &ds.CreateTeamRequest{Name: "name"}
		res, err := tc.c.AddNewTeam(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResurceNotFound)
	})

	t.Run("error on LastInsertId of AddNewTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), errTest)

		req := &ds.CreateTeamRequest{Name: "name"}
		_, err := tc.c.AddNewTeam(req)
		require.NotNil(t, err)
	})

	t.Run("error on AddMemberToTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.CreateTeamRequest{Name: "name"}
		_, err := tc.c.AddNewTeam(req)
		require.NotNil(t, err)
	})

	t.Run("error on RowsAffected of AddMemberToTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), errTest)

		req := &ds.CreateTeamRequest{Name: "name"}
		_, err := tc.c.AddNewTeam(req)
		require.NotNil(t, err)
	})

	t.Run("No on RowsAffected of AddMemberToTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().AddNewTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.CreateTeamRequest{Name: "name"}
		_, err := tc.c.AddNewTeam(req)
		require.NotNil(t, err)
	})
}

func TestGetUserTeams(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().GetUserTeams(gomock.Any(), gomock.Any()).
			Return([]sqlc.GetUserTeamsRow{{TeamID: 1, Name: "name"}}, nil)

		res, err := tc.c.GetUserTeams(1)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on GetUserTeams", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().GetUserTeams(gomock.Any(), gomock.Any()).
			Return(nil, errTest)

		_, err := tc.c.GetUserTeams(1)
		require.NotNil(t, err)
	})

	t.Run("No rows on GetUserTeams", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().GetUserTeams(gomock.Any(), gomock.Any()).
			Return(nil, errNoRow)

		res, err := tc.c.GetUserTeams(1)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})
}

func TestAddUserToUserTeam(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		res, err := tc.c.AddUserToUserTeam(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on GetTeamOwner", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(0), errTest)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		_, err := tc.c.AddUserToUserTeam(req)
		require.NotNil(t, err)
	})

	t.Run("No rows on GetTeamOwner", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(0), errNoRow)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		res, err := tc.c.AddUserToUserTeam(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusUserNotFound)
	})

	t.Run("No rows on GetTeamOwner", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(4), nil)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		res, err := tc.c.AddUserToUserTeam(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusNotOwner)
	})

	t.Run("error on AddMemberToTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		_, err := tc.c.AddUserToUserTeam(req)
		require.NotNil(t, err)
	})

	t.Run("Duplicate on AddMemberToTeam", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errDuplicate)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		res, err := tc.c.AddUserToUserTeam(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusUserAlreadyExists)
	})

	t.Run("error on RowsAffected", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), errTest)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		_, err := tc.c.AddUserToUserTeam(req)
		require.NotNil(t, err)
	})

	t.Run("No RowsAffected", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTeamOwner(gomock.Any(), gomock.Any()).
			Return(int64(1), nil)
		tc.querier.EXPECT().AddMemberToTeam(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.DBInviteUserToTeamRequest{InviterId: 1}
		_, err := tc.c.AddUserToUserTeam(req)
		require.NotNil(t, err)
	})
}
