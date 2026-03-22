package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddNewTask(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on IsTeamMember", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(false, errTest)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("Not an owner", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(false, nil)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusNotOwner)
	})

	t.Run("error on AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, errTest)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("Long data on AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, errlongData)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusDataTooLong)
	})

	t.Run("Duplicate data on AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, errDuplicate)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResourceAlreadyExists)
	})

	t.Run("Foreign key on AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, errForeignKey)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResurceNotFound)
	})

	t.Run("error on RowsAffected of AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), errTest)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("No RowsAffected of AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("error on LastInsertId of AddNewTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(0), errTest)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("error on AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("Long data on AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errlongData)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusDataTooLong)
	})

	t.Run("Duplicate data on AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errDuplicate)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResourceAlreadyExists)
	})

	t.Run("Foreign key on AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errForeignKey)

		req := &ds.DBCreateTaskRequest{}
		res, err := tc.c.AddNewTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResurceNotFound)
	})

	t.Run("error on RowsAffected of AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), errTest)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})

	t.Run("No RowsAffected of AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().AddNewTask(gomock.Any(), gomock.Any()).Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.sqlRes.EXPECT().LastInsertId().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.DBCreateTaskRequest{}
		_, err := tc.c.AddNewTask(req)
		require.NotNil(t, err)
	})
}

func TestGetTasks(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().GetTasks(gomock.Any(), gomock.Any()).
			Return([]sqlc.GetTasksRow{{TaskID: 1}, {TaskID: 2}}, nil)
		tc.querier.EXPECT().GetTasksComments(gomock.Any(), gomock.Any()).
			Return([]sqlc.TasksComment{{TaskID: 1, Comment: "sesehehe"}}, nil)

		req := &ds.GetTasksRequest{}
		res, err := tc.c.GetTasks(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on IsTeamMember", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(false, errTest)

		req := &ds.GetTasksRequest{}
		_, err := tc.c.GetTasks(req)
		require.NotNil(t, err)
	})

	t.Run("Not a member", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(false, nil)

		req := &ds.GetTasksRequest{}
		res, err := tc.c.GetTasks(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusNotMember)
	})

	t.Run("error on GetTasks", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().GetTasks(gomock.Any(), gomock.Any()).
			Return(nil, errTest)

		req := &ds.GetTasksRequest{}
		_, err := tc.c.GetTasks(req)
		require.NotNil(t, err)

	})

	t.Run("Ok no tasks", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().GetTasks(gomock.Any(), gomock.Any()).
			Return([]sqlc.GetTasksRow{}, nil)

		req := &ds.GetTasksRequest{}
		res, err := tc.c.GetTasks(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on GetTasksComments", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().GetTasks(gomock.Any(), gomock.Any()).
			Return([]sqlc.GetTasksRow{{TaskID: 1}, {TaskID: 2}}, nil)
		tc.querier.EXPECT().GetTasksComments(gomock.Any(), gomock.Any()).
			Return(nil, errTest)

		req := &ds.GetTasksRequest{}
		_, err := tc.c.GetTasks(req)
		require.NotNil(t, err)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("No rown on GetTaskForUpdate", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{}, errNoRow)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResurceNotFound)
	})

	t.Run("Invalid version", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)

		req := &ds.DBUpdateTaskRequest{Version: 10, Status: "completed"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusIvalidVersion)
	})

	t.Run("Conflict version", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 10, Status: "todo"}, nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completed"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusConflict)
	})

	t.Run("No change needed", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "todo"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on IsTeamMember", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(false, errTest)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

	t.Run("Not a member", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(false, nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusNotMember)
	})

	t.Run("Not Admin trying change team ID", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo", TeamID: 2}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes", TeamId: 1}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusForbidden)
	})

	t.Run("error on UpdateTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

	t.Run("Long data on UpdateTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errlongData)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusDataTooLong)
	})

	t.Run("error on RowsAffected of UpdateTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), errTest)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

	t.Run("No RowsAffected of UpdateTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

	t.Run("error on AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

	t.Run("Long data on AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errlongData)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		res, err := tc.c.UpdateTask(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusDataTooLong)
	})

	t.Run("error on RowsAffected of AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), errTest)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

	t.Run("Now RowsAffected of AddChangeToTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTaskForUpdate(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskForUpdateRow{Version: 1, Status: "todo"}, nil)
		tc.querier.EXPECT().IsTeamMember(gomock.Any(), gomock.Any()).Return(true, nil)
		tc.querier.EXPECT().UpdateTask(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)
		tc.querier.EXPECT().AddChangeToTaskHistory(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.DBUpdateTaskRequest{Version: 1, Status: "completes"}
		_, err := tc.c.UpdateTask(req)
		require.NotNil(t, err)
	})

}

func TestGetTaskHistory(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTask(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskRow{}, nil)
		tc.querier.EXPECT().GetTaskHistory(gomock.Any(), gomock.Any()).
			Return([]sqlc.GetTaskHistoryRow{
				{
					ChangedBy: 1,
					Payload:   json.RawMessage(patchTest),
				},
			}, nil)
		tc.querier.EXPECT().GetTaskHistory(gomock.Any(), gomock.Any()).
			Return(nil, nil)

		req := &ds.GetTaskHistoryRequest{TaskId: 1}
		res, err := tc.c.GetTaskHistory(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on GetTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTask(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskRow{}, errTest)

		req := &ds.GetTaskHistoryRequest{TaskId: 1}
		_, err := tc.c.GetTaskHistory(req)
		require.NotNil(t, err)
	})

	t.Run("Now rows on GetTask", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTask(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskRow{}, errNoRow)

		req := &ds.GetTaskHistoryRequest{TaskId: 1}
		res, err := tc.c.GetTaskHistory(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResurceNotFound)
	})

	t.Run("error on GetTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTask(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskRow{}, nil)
		tc.querier.EXPECT().GetTaskHistory(gomock.Any(), gomock.Any()).
			Return(nil, errTest)

		req := &ds.GetTaskHistoryRequest{TaskId: 1}
		_, err := tc.c.GetTaskHistory(req)
		require.NotNil(t, err)
	})

	t.Run("No rown on  GetTaskHistory", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		txExec := func(opts *sql.TxOptions, fn func(ctx context.Context, q IQuerier) error) error {
			return fn(tc.ctx, tc.querier)
		}
		tc.db.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(txExec)

		tc.querier.EXPECT().GetTask(gomock.Any(), gomock.Any()).
			Return(sqlc.GetTaskRow{}, nil)
		tc.querier.EXPECT().GetTaskHistory(gomock.Any(), gomock.Any()).
			Return(nil, errNoRow)

		req := &ds.GetTaskHistoryRequest{TaskId: 1}
		_, err := tc.c.GetTaskHistory(req)
		require.Nil(t, err)
	})
}

func TestAddTaskComment(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().
			AddTaskComment(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), nil)

		req := &ds.AddTaskCommentRequest{TaskId: 1}
		res, err := tc.c.AddTaskComment(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusSuccess)
	})

	t.Run("error on AddTaskComment", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().
			AddTaskComment(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errTest)

		req := &ds.AddTaskCommentRequest{TaskId: 1}
		_, err := tc.c.AddTaskComment(req)
		require.NotNil(t, err)
	})

	t.Run("error on RowsAffected", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().
			AddTaskComment(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(1), errTest)

		req := &ds.AddTaskCommentRequest{TaskId: 1}
		_, err := tc.c.AddTaskComment(req)
		require.NotNil(t, err)
	})

	t.Run("No RowsAffected", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().
			AddTaskComment(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, nil)
		tc.sqlRes.EXPECT().RowsAffected().Return(int64(0), nil)

		req := &ds.AddTaskCommentRequest{TaskId: 1}
		_, err := tc.c.AddTaskComment(req)
		require.NotNil(t, err)
	})

	t.Run("Foreign key on AddTaskComment", func(t *testing.T) {
		t.Parallel()

		tc := newTestClient(t)

		tc.db.EXPECT().CtxWithCancel().Return(context.Background(), func() {})
		tc.db.EXPECT().Querier().Return(tc.querier)

		tc.querier.EXPECT().
			AddTaskComment(gomock.Any(), gomock.Any()).
			Return(tc.sqlRes, errForeignKey)

		req := &ds.AddTaskCommentRequest{TaskId: 1}
		res, err := tc.c.AddTaskComment(req)
		require.Nil(t, err)
		require.Equal(t, res.GetStatus(), ds.StatusResurceNotFound)
	})
}
