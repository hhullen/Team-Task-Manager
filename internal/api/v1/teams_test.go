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

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateTeamRequest{
			Name:        "name",
			Description: "des",
		}

		ta.appService.EXPECT().CreateTeam(gomock.Any()).Return(&ds.CreateTeamResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.CreateTeam(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.CreateTeamResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateTeamRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.CreateTeam(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.CreateTeamResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.CreateTeamRequest{
			Name:        "name",
			Description: "desc",
		}

		ta.appService.EXPECT().CreateTeam(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.CreateTeam(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.CreateTeamResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestListUserTeams(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.ListUserTeamsRequest{}

		ta.appService.EXPECT().ListUserTeams(gomock.Any()).Return(&ds.ListUserTeamsResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.ListUserTeams(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.ListUserTeamsResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.ListUserTeamsRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		// ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.ListUserTeams(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.ListUserTeamsResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedExctractingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.ListUserTeamsRequest{}

		ta.appService.EXPECT().ListUserTeams(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.ListUserTeams(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.ListUserTeamsResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestInviteUserToTeam(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.InviteUserToTeamRequest{
			UserLoginToInvite: "log",
		}

		ta.appService.EXPECT().InviteUserToTeam(gomock.Any()).Return(&ds.InviteUserToTeamResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/teams/23/invite", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.InviteUserToTeam(w, r.WithContext(ctx))

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.InviteUserToTeamResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.InviteUserToTeamRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/teams/23/invite", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.InviteUserToTeam(w, r.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.InviteUserToTeamResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.InviteUserToTeamRequest{
			UserLoginToInvite: "log",
		}

		ta.appService.EXPECT().InviteUserToTeam(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/teams/23/invite", bytes.NewReader(body))
		r.SetPathValue("id", "23")
		ctx := context.WithValue(r.Context(), ds.UserIDKey, int64(222))
		ctx = context.WithValue(ctx, ds.UserRoleKey, "user")
		w := httptest.NewRecorder()
		ta.a.InviteUserToTeam(w, r.WithContext(ctx))

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.InviteUserToTeamResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}
