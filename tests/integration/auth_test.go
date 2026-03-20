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

func (s *ServicesTestSuite) TestAuthRegister() {
	payload := map[string]string{
		"login":    "test1",
		"name":     "test1",
		"password": "test1",
	}
	uri := apiPrefix + "/register"

	s.Run("Ok", func() {
		w := s.JSONBodyRequest(http.MethodPost,payload, uri, [][2]string{})

		s.Equal(http.StatusOK, w.Code)

		contentType := w.Header().Get("Content-Type")
		s.Equal(contentType, jsonContentType)

		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users_auth WHERE login = ?", payload["login"]).Scan(&count)
		s.NoError(err)
		s.Equal(1, count)
	})

	s.Run("Already registered", func() {
		w := s.JSONBodyRequest(http.MethodPost,payload, uri, [][2]string{})

		s.Equal(http.StatusConflict, w.Code)

		contentType := w.Header().Get("Content-Type")
		s.Equal(contentType, jsonContentType)
	})

	s.Run("Without required field", func() {
		delete(payload, "login")
		w := s.JSONBodyRequest(http.MethodPost,payload, uri, [][2]string{})

		s.Equal(http.StatusBadRequest, w.Code)

		contentType := w.Header().Get("Content-Type")
		s.Equal(contentType, jsonContentType)
	})
}

func (s *ServicesTestSuite) TestAuthLogin() {
	payload := map[string]string{
		"login":    "test2",
		"name":     "test2",
		"password": "test2",
	}

	w := s.JSONBodyRequest(http.MethodPost,payload, apiPrefix+"/register", [][2]string{})
	s.Equal(http.StatusOK, w.Code)

	payload = map[string]string{
		"login":    "test2",
		"password": "test2",
	}

	uri := apiPrefix + "/login"

	s.Run("Ok", func() {
		w := s.JSONBodyRequest(http.MethodPost,payload, uri, [][2]string{})
		s.Equal(http.StatusOK, w.Code)

		res := w.Result()

		v := map[string]string{}

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
		w := s.JSONBodyRequest(http.MethodPost,payload, uri, [][2]string{})
		s.Equal(http.StatusUnauthorized, w.Code)
	})

	s.Run("Without required field", func() {
		delete(payload, "password")
		w := s.JSONBodyRequest(http.MethodPost,payload, uri, [][2]string{})
		s.Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ServicesTestSuite) TestAuthRefresh() {
	payload := map[string]string{
		"login":    "test3",
		"name":     "test3",
		"password": "test3",
	}

	w := s.JSONBodyRequest(http.MethodPost,payload, apiPrefix+"/register", [][2]string{})
	s.Equal(http.StatusOK, w.Code)

	payload = map[string]string{
		"login":    "test3",
		"password": "test3",
	}

	w = s.JSONBodyRequest(http.MethodPost,payload, apiPrefix+"/login", [][2]string{})
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

		v := map[string]string{}

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
