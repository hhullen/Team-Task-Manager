package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	ds "team-task-manager/internal/datastruct"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestExtractJsonBody(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		body, _ := json.Marshal(v)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))

		res := ds.LoginRequest{}
		err := extractJsonBody(r, &res)
		require.Nil(t, err)

		require.Equal(t, v.Login, res.Login)
		require.Equal(t, v.Password, res.Password)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{asdfae}"))

		res := ds.LoginRequest{}
		err := extractJsonBody(r, &res)
		require.NotNil(t, err)

	})
}

func TestWriteJsonResponse(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		v := ds.AddTaskCommentResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		}

		w := httptest.NewRecorder()
		err := writeJsonResponse(w, &v)
		require.Nil(t, err)

		res := ds.AddTaskCommentResponse{}
		err = json.Unmarshal(w.Body.Bytes(), &res)
		require.Nil(t, err)

		require.Equal(t, v.GetStatus(), res.GetStatus())
	})
}

func TestExtractSchemaQuery(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(""))
		q := r.URL.Query()
		q.Add("team_id", "999")
		q.Add("status", "todo")
		r.URL.RawQuery = q.Encode()
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")

		v := ds.GetTasksRequest{}
		err := extractSchemaQuery(r.WithContext(ctx), &v)
		require.Nil(t, err)

		require.Equal(t, v.TeamId, int64(999))
		require.Equal(t, v.Status, ds.TaskStatus("todo"))
		require.Equal(t, v.UserID, int64(222))
		require.Equal(t, v.Role, "user")

	})
}

func TestStructValidator(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		err := structValidator(v)
		require.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login: "login",
			},
		}

		err := structValidator(v)
		require.NotNil(t, err)
	})
}

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		Exec(ExecArgs[ds.LoginRequest, ds.LoginResponse]{
			api: ta.a,
			serviceFunc: func(lr *ds.LoginRequest) *ds.LoginResponse {
				return &ds.LoginResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: extractJsonBody[*ds.LoginRequest],
			responseWriter:   writeJsonResponse[*ds.LoginResponse],
			httpRequest:      r,
			httpResponse:     w,
		})

		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusSuccess)
	})

	t.Run("extracting error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[ds.LoginRequest, ds.LoginResponse]{
			api: ta.a,
			serviceFunc: func(lr *ds.LoginRequest) *ds.LoginResponse {
				return &ds.LoginResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: func(r *http.Request, v *ds.LoginRequest) error {
				return errTest
			},
			responseWriter: writeJsonResponse[*ds.LoginResponse],
			httpRequest:    r,
			httpResponse:   w,
		})

		require.Equal(t, w.Code, http.StatusBadRequest)
		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusFailedExctractingRequest)
	})

	t.Run("validating error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[ds.LoginRequest, ds.LoginResponse]{
			api: ta.a,
			serviceFunc: func(lr *ds.LoginRequest) *ds.LoginResponse {
				return &ds.LoginResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: extractJsonBody[*ds.LoginRequest],
			responseWriter:   writeJsonResponse[*ds.LoginResponse],
			httpRequest:      r,
			httpResponse:     w,
			validator: func(s *ds.LoginRequest) error {
				return errTest
			},
		})

		require.Equal(t, w.Code, http.StatusBadRequest)
		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusFailedValidatingRequest)
	})

	t.Run("service func error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[ds.LoginRequest, ds.LoginResponse]{
			api: ta.a,
			serviceFunc: func(lr *ds.LoginRequest) *ds.LoginResponse {
				return nil
			},
			requestExtractor: extractJsonBody[*ds.LoginRequest],
			responseWriter:   writeJsonResponse[*ds.LoginResponse],
			httpRequest:      r,
			httpResponse:     w,
		})

		require.Equal(t, w.Code, http.StatusInternalServerError)
		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))
		require.Equal(t, respS.GetStatus(), ds.StatusServiceError)
	})

	t.Run("response writer error", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		v := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		body, _ := json.Marshal(v)

		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		Exec(ExecArgs[ds.LoginRequest, ds.LoginResponse]{
			api: ta.a,
			serviceFunc: func(lr *ds.LoginRequest) *ds.LoginResponse {
				return &ds.LoginResponse{
					Status: ds.Status{Message: ds.StatusSuccess},
				}
			},
			requestExtractor: extractJsonBody[*ds.LoginRequest],
			responseWriter: func(w http.ResponseWriter, v *ds.LoginResponse) error {
				return errTest
			},
			httpRequest:  r,
			httpResponse: w,
		})

		require.Equal(t, w.Code, http.StatusInternalServerError)
		
	})
}
