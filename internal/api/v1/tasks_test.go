package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	ds "team-task-manager/internal/datastruct"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateTaskRequest{
			AssigneeLogin: "log",
			Subject:       "sub",
			Description:   "des",
			Status:        "in progress",
			TeamId:        1,
		}

		ta.appService.EXPECT().AddNewTask(gomock.Any()).Return(&ds.CreateTaskResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.CreateTask(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.CreateTaskResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateTaskRequest{

			Status: "in progress",
			TeamId: 1,
		}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.CreateTask(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.CreateTaskResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateTaskRequest{
			AssigneeLogin: "log",
			Subject:       "sub",
			Description:   "des",
			Status:        "in progress",
			TeamId:        1,
		}

		ta.appService.EXPECT().AddNewTask(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.CreateTask(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.CreateTaskResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestGetTasks(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetTasksRequest{}

		ta.appService.EXPECT().GetTasks(gomock.Any()).Return(&ds.GetTasksResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		q := r.URL.Query()
		q.Add("team_id", "999")
		q.Add("status", "todo")
		q.Add("assignee_login", "log")
		q.Add("offset", "0")
		q.Add("limit", "10")
		q.Add("avoid_cache", "true")
		r.URL.RawQuery = q.Encode()
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.GetTasks(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.GetTasksResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetTasksRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		q := r.URL.Query()
		q.Add("team_id", "999")
		// q.Add("status", "todo")
		// q.Add("assignee_login", "log")
		q.Add("offset", "0")
		q.Add("limit", "10")
		q.Add("avoid_cache", "true")
		r.URL.RawQuery = q.Encode()
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.GetTasks(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.GetTasksResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetTasksRequest{}

		ta.appService.EXPECT().GetTasks(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		q := r.URL.Query()
		q.Add("team_id", "999")
		q.Add("status", "todo")
		q.Add("assignee_login", "log")
		q.Add("offset", "0")
		q.Add("limit", "10")
		q.Add("avoid_cache", "true")
		r.URL.RawQuery = q.Encode()
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.GetTasks(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.GetTasksResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.UpdateTaskRequest{
			AssigneeLogin: "log",
			Subject:       "sub",
			Description:   "des",
			Status:        "in progress",
			TeamId:        1,
			Version:       2,
		}

		ta.appService.EXPECT().UpdateTask(gomock.Any()).Return(&ds.UpdateTaskResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPut, "/tasks/{id}", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.UpdateTask(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.UpdateTaskResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.UpdateTaskRequest{

			Status: "in progress",
			TeamId: 1,
		}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.UpdateTask(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.UpdateTaskResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.UpdateTaskRequest{
			AssigneeLogin: "log",
			Subject:       "sub",
			Description:   "des",
			Status:        "in progress",
			TeamId:        1,
			Version:       2,
		}

		ta.appService.EXPECT().UpdateTask(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.UpdateTask(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.UpdateTaskResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestGetTaskHistory(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetTaskHistoryRequest{}

		ta.appService.EXPECT().GetTaskHistory(gomock.Any()).Return(&ds.GetTaskHistoryResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPut, "/tasks/23/history", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.GetTaskHistory(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.GetTaskHistoryResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetTaskHistoryRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23/history", bytes.NewReader(body))
		// r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.GetTaskHistory(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.GetTaskHistoryResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedExctractingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.GetTaskHistoryRequest{}

		ta.appService.EXPECT().GetTaskHistory(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23/history", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.GetTaskHistory(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.GetTaskHistoryResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestAddTaskComment(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.AddTaskCommentRequest{
			Text: "xfvfx",
		}

		ta.appService.EXPECT().AddTaskComment(gomock.Any()).Return(&ds.AddTaskCommentResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23/comment", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.AddTaskComment(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.AddTaskCommentResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.AddTaskCommentRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23/comment", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.AddTaskComment(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.AddTaskCommentResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.AddTaskCommentRequest{
			Text: "sdf",
		}

		ta.appService.EXPECT().AddTaskComment(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/tasks/23/comment", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.AddTaskComment(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.AddTaskCommentResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}
