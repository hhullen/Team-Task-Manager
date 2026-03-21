package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

const (
	jsonContentType = "application/json"
)

func (s *ServicesTestSuite) Test_Auth_Register() {
	payload := map[string]any{
		"login":    "test_auth_register",
		"name":     "test_auth_register",
		"password": "test_auth_register",
	}
	uri := apiPrefix + "/register"

	s.Run("Ok", func() {
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{})

		s.Equal(http.StatusOK, w.Code)

		contentType := w.Header().Get("Content-Type")
		s.Equal(contentType, jsonContentType)

		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users_auth WHERE login = ?", payload["login"]).Scan(&count)
		s.NoError(err)
		s.Equal(1, count)
	})

	s.Run("Already registered", func() {
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{})

		s.Equal(http.StatusConflict, w.Code)

		contentType := w.Header().Get("Content-Type")
		s.Equal(contentType, jsonContentType)
	})

	s.Run("Without required field", func() {
		payload := map[string]any{
			"name":     "test_auth_register",
			"password": "test_auth_register",
		}
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{})

		s.Equal(http.StatusBadRequest, w.Code)

		contentType := w.Header().Get("Content-Type")
		s.Equal(contentType, jsonContentType)
	})
}

func (s *ServicesTestSuite) Test_Auth_Login() {
	payload := map[string]any{
		"login":    "test_auth_login",
		"name":     "test_auth_login",
		"password": "test_auth_login",
	}

	w := s.JSONBodyRequest(http.MethodPost, payload, apiPrefix+"/register", [][2]string{})
	s.Equal(http.StatusOK, w.Code)

	payload = map[string]any{
		"login":    "test_auth_login",
		"password": "test_auth_login",
	}

	uri := apiPrefix + "/login"

	s.Run("Ok", func() {
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{})
		s.Equal(http.StatusOK, w.Code)

		res := w.Result()

		v := map[string]any{}

		err := json.NewDecoder(w.Body).Decode(&v)
		s.Nil(err)

		at, ok := v["access_token"]
		s.True(ok)
		s.True(at != "")

		var coockie *http.Cookie
		for _, v := range res.Cookies() {
			if v.Name == "refresh_token" {
				coockie = v
			}
		}

		s.NotNil(coockie)
		s.True(coockie.Value != "")

	})

	s.Run("Wrong password", func() {
		payload["password"] = "wrong"
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{})
		s.Equal(http.StatusUnauthorized, w.Code)
	})

	s.Run("Without required field", func() {
		payload = map[string]any{
			"login": "test_auth_login",
		}
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{})
		s.Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ServicesTestSuite) Test_Auth_Refresh() {
	payload := map[string]any{
		"login":    "test_auth_refresh",
		"name":     "test_auth_refresh",
		"password": "test_auth_refresh",
	}

	w := s.JSONBodyRequest(http.MethodPost, payload, apiPrefix+"/register", [][2]string{})
	s.Equal(http.StatusOK, w.Code)

	payload = map[string]any{
		"login":    "test_auth_refresh",
		"password": "test_auth_refresh",
	}

	w = s.JSONBodyRequest(http.MethodPost, payload, apiPrefix+"/login", [][2]string{})
	s.Equal(http.StatusOK, w.Code)

	res := w.Result()
	var coockie *http.Cookie
	for _, v := range res.Cookies() {
		if v.Name == "refresh_token" {
			coockie = v
		}
	}

	uri := apiPrefix + "/refresh"

	s.Run("Ok", func() {
		req := httptest.NewRequest(http.MethodPost, uri, strings.NewReader(""))
		req.Header.Set("Content-Type", jsonContentType)
		req.AddCookie(coockie)

		w := httptest.NewRecorder()

		s.api.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)

		v := map[string]any{}

		err := json.NewDecoder(w.Body).Decode(&v)
		s.Nil(err)

		at, ok := v["access_token"]
		s.True(ok)
		s.True(at != "")

		var coockie *http.Cookie
		for _, v := range res.Cookies() {
			if v.Name == "refresh_token" {
				coockie = v
			}
		}

		s.NotNil(coockie)
		s.True(coockie.Value != "")
	})

	s.Run("Wrong cookie", func() {
		wrongCoockie := *coockie
		wrongCoockie.Value = "Wrong"
		req := httptest.NewRequest(http.MethodPost, uri, strings.NewReader(""))
		req.Header.Set("Content-Type", jsonContentType)
		req.AddCookie(&wrongCoockie)

		w := httptest.NewRecorder()

		s.api.ServeHTTP(w, req)

		s.Equal(http.StatusUnauthorized, w.Code)
	})
}
