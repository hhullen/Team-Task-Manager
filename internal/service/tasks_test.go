package service

import (
	ds "team-task-manager/internal/datastruct"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddNewTask(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().AddNewTask(gomock.Any()).
			Return(&ds.CreateTaskResponse{}, nil)

		req := &ds.CreateTaskRequest{}
		res := st.s.AddNewTask(req)
		require.NotNil(t, res)
	})

	t.Run("error on AddNewTask", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().AddNewTask(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.CreateTaskRequest{}
		res := st.s.AddNewTask(req)
		require.Nil(t, res)
	})
}

func TestGetTasks(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().GetTasks(gomock.Any()).
			Return(&ds.GetTasksResponse{}, nil)

		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		req := &ds.GetTasksRequest{AvoidCacheFlag: ds.AvoidCacheFlag{Flag: true}}
		res := st.s.GetTasks(req)
		require.NotNil(t, res)
	})

	t.Run("error on GetTasks", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().GetTasks(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.Any())

		req := &ds.GetTasksRequest{AvoidCacheFlag: ds.AvoidCacheFlag{Flag: true}}
		res := st.s.GetTasks(req)
		require.Nil(t, res)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.appStore.EXPECT().UpdateTask(gomock.Any()).
			Return(&ds.UpdateTaskResponse{}, nil)

		req := &ds.UpdateTaskRequest{}
		res := st.s.UpdateTask(req)
		require.NotNil(t, res)
	})

	t.Run("error on UpdateTask", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().UpdateTask(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.Any())

		req := &ds.UpdateTaskRequest{}
		res := st.s.UpdateTask(req)
		require.Nil(t, res)
	})
}

func TestGetTaskHistory(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)

		st.appStore.EXPECT().GetTaskHistory(gomock.Any()).
			Return(&ds.GetTaskHistoryResponse{}, nil)

		req := &ds.GetTaskHistoryRequest{AvoidCacheFlag: ds.AvoidCacheFlag{Flag: true}}
		res := st.s.GetTaskHistory(req)
		require.NotNil(t, res)
	})

	t.Run("error on GetTaskHistory", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().GetTaskHistory(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.Any())

		req := &ds.GetTaskHistoryRequest{AvoidCacheFlag: ds.AvoidCacheFlag{Flag: true}}
		res := st.s.GetTaskHistory(req)
		require.Nil(t, res)
	})
}

func TestAddTaskComment(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().AddTaskComment(gomock.Any()).
			Return(&ds.AddTaskCommentResponse{}, nil)

		req := &ds.AddTaskCommentRequest{}
		res := st.s.AddTaskComment(req)
		require.NotNil(t, res)
	})

	t.Run("error on AddTaskComment", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().AddTaskComment(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.Any())

		req := &ds.AddTaskCommentRequest{}
		res := st.s.AddTaskComment(req)
		require.Nil(t, res)
	})
}
