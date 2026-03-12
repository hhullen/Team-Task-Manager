package api

import (
	"errors"
	"fmt"
	"net/http"

	ds "team-task-manager/internal/datastruct"
)

const (
	registerPrefix = apiPrefix + "/register"
	loginPrefix    = apiPrefix + "/login"
	refreshPrefix  = apiPrefix + "/refresh"
)

func (a *API) setupAuthHandlers() {
	a.router.Handle(pattern(http.MethodPost, registerPrefix),
		globalRateLimitedMiddleware(a, http.HandlerFunc(a.Register)))
	a.router.Handle(pattern(http.MethodPost, loginPrefix),
		globalRateLimitedMiddleware(a, http.HandlerFunc(a.Login)))
	a.router.Handle(pattern(http.MethodPost, refreshPrefix),
		globalRateLimitedMiddleware(a, http.HandlerFunc(a.Refresh)))
}

// Register register new user
// @Summary      register new user
// @Description  register new user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body      ds.RegisterRequest  true "login and password and info"
// @Success      200   {object}  ds.RegisterResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /register [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.RegisterRequest, ds.RegisterResponse]{
		serviceFunc:      a.authService.RegisterUser,
		responseWriter:   writeJsonResponse[*ds.RegisterResponse],
		requestExtractor: extractJsonBody[*ds.RegisterRequest],
		httpResponse:     &w,
		httpRequest:      r,
		api:              a,
	})
}

// Login login user
// @Summary      login user
// @Description  login user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body      ds.LoginRequest  true "login and password"
// @Success      200   {object}  ds.LoginResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /login [post]
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.LoginRequest, ds.LoginResponse]{
		serviceFunc: a.authService.LoginUser,
		responseWriter: func(w *http.ResponseWriter, v *ds.LoginResponse) error {
			cookie := &http.Cookie{
				Name:     "refresh_token",
				Value:    v.RefreshToken.RefreshToken,
				Path:     "/",
				Domain:   "",
				MaxAge:   int(ds.DefaultRefreshTokenTTL.Seconds()),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}

			http.SetCookie(*w, cookie)

			v.RefreshToken.RefreshToken = ""

			return writeJsonResponse(w, v)

		},
		requestExtractor: extractJsonBody[*ds.LoginRequest],
		httpResponse:     &w,
		httpRequest:      r,
		api:              a,
	})
}

// Refresh refresh tokens
// @Summary      refresh tokens
// @Description  refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200           {object}  ds.RefreshResponse
// @Failure      400           {object}  ds.Status
// @Failure      500           {object}  ds.Status
// @Router       /refresh [post]
func (a *API) Refresh(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.RefreshRequest, ds.RefreshResponse]{
		serviceFunc: a.authService.Refresh,
		responseWriter: func(w *http.ResponseWriter, v *ds.RefreshResponse) error {
			cookie := &http.Cookie{
				Name:     "refresh_token",
				Value:    v.RefreshToken.RefreshToken,
				Path:     "/",
				Domain:   "",
				MaxAge:   int(ds.DefaultRefreshTokenTTL.Seconds()),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}

			http.SetCookie(*w, cookie)

			v.RefreshToken.RefreshToken = ""

			return writeJsonResponse(w, v)

		},
		requestExtractor: func(r *http.Request, v *ds.RefreshRequest) error {
			token, err := r.Cookie("refresh_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return fmt.Errorf("cookie 'refresh_token' was not provided")
				}
				return err
			}
			v.RefreshToken.RefreshToken = token.Value
			return nil
		},
		httpResponse: &w,
		httpRequest:  r,
		api:          a,
	})
}
