package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ds "team-task-manager/internal/datastruct"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.RegisterRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
			UserInfo: ds.UserInfo{
				Name: "name",
			},
		}

		ta.authService.EXPECT().RegisterUser(gomock.Any()).Return(&ds.RegisterResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Register(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.RegisterResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.RegisterRequest{
			UserCreds: ds.UserCreds{
				Login:    "",
				Password: "password",
			},
			UserInfo: ds.UserInfo{
				Name: "name",
			},
		}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Register(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.RegisterResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.RegisterRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
			UserInfo: ds.UserInfo{
				Name: "name",
			},
		}

		ta.authService.EXPECT().RegisterUser(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Register(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.RegisterResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		accessToken := ds.AccessToken{AccessToken: "test_access"}
		refreshToken := ds.RefreshToken{RefreshToken: "test_refresh"}
		ta.authService.EXPECT().LoginUser(gomock.Any()).Return(&ds.LoginResponse{
			Status:       ds.Status{Message: ds.StatusSuccess},
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Login(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		respRefreshToken := ""
		for _, c := range w.Result().Cookies() {
			if c.Value == refreshToken.RefreshToken {
				respRefreshToken = c.Value
			}
		}
		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
		require.Equal(t, accessToken.AccessToken, respS.AccessToken.AccessToken)
		require.Equal(t, refreshToken.RefreshToken, respRefreshToken)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "",
				Password: "password",
			},
		}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Login(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedValidatingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.LoginRequest{
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
		}

		ta.authService.EXPECT().LoginUser(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Login(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.LoginResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}

func TestRefresh(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.RefreshRequest{}

		accessToken := ds.AccessToken{AccessToken: "test_access"}
		refreshToken := ds.RefreshToken{RefreshToken: "test_refresh"}
		ta.authService.EXPECT().Refresh(gomock.Any()).Return(&ds.RefreshResponse{
			Status:       ds.Status{Message: ds.StatusSuccess},
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		r.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "old_token",
			Path:     "/",
			Domain:   "",
			MaxAge:   int(ds.DefaultRefreshTokenTTL.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w := httptest.NewRecorder()
		ta.a.Refresh(w, r)

		require.Equal(t, http.StatusOK, w.Code)

		respS := ds.RefreshResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		respRefreshToken := ""
		for _, c := range w.Result().Cookies() {
			if c.Value == refreshToken.RefreshToken {
				respRefreshToken = c.Value
			}
		}
		require.Equal(t, ds.StatusSuccess, respS.Status.Message)
		require.Equal(t, accessToken.AccessToken, respS.AccessToken.AccessToken)
		require.Equal(t, refreshToken.RefreshToken, respRefreshToken)
	})

	t.Run("No required field", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.RefreshRequest{}

		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		w := httptest.NewRecorder()
		ta.a.Refresh(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)

		respS := ds.RefreshResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusFailedExctractingRequest, respS.Status.Message)
	})

	t.Run("No response from service", func(t *testing.T) {
		t.Parallel()

		ta := newTestAPI(t)

		reqS := ds.RefreshRequest{}

		ta.authService.EXPECT().Refresh(gomock.Any()).Return(nil)
		ta.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		body, _ := json.Marshal(reqS)
		r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
		r.AddCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "old_token",
			Path:     "/",
			Domain:   "",
			MaxAge:   int(ds.DefaultRefreshTokenTTL.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w := httptest.NewRecorder()
		ta.a.Refresh(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		respS := ds.RefreshResponse{}
		require.Nil(t, json.Unmarshal(w.Body.Bytes(), &respS))

		require.Equal(t, ds.StatusServiceError, respS.Status.Message)
	})
}
