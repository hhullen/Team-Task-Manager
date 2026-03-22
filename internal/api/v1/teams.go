package api

import (
	"fmt"
	"net/http"
	"strconv"
	ds "team-task-manager/internal/datastruct"
)

const (
	teamsPrefix       = apiPrefix + "/teams"
	teamsInvitePrefix = teamsPrefix + "/{id}/invite"
)

func (a *API) setupTeamsHandlers() {
	a.router.Handle(pattern(http.MethodPost, teamsPrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.CreateTeam)))
	a.router.Handle(pattern(http.MethodGet, teamsPrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.ListUserTeams)))
	a.router.Handle(pattern(http.MethodPost, teamsInvitePrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.InviteUserToTeam)))
}

// CreateTeam create new team
// @Summary      create new team
// @Description  create new team.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input body      ds.CreateTeamRequest  true "team description"
// @Success      200   {object}  ds.CreateTeamResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /teams [post]
func (a *API) CreateTeam(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.CreateTeamRequest, ds.CreateTeamResponse]{
		serviceFunc:      a.appService.CreateTeam,
		responseWriter:   writeJsonResponse[*ds.CreateTeamResponse],
		requestExtractor: extractJsonBody[*ds.CreateTeamRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// CreateTeam list teams the user member of
// @Summary      list teams the user member of
// @Description  list teams the user member of .
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200   {object}  ds.ListUserTeamsResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /teams [get]
func (a *API) ListUserTeams(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.ListUserTeamsRequest, ds.ListUserTeamsResponse]{
		serviceFunc:      a.appService.ListUserTeams,
		responseWriter:   writeJsonResponse[*ds.ListUserTeamsResponse],
		requestExtractor: extractJWTCredsOnly[*ds.ListUserTeamsRequest],
		httpResponse:     w,
		httpRequest:      r,
		api:              a,
	})
}

// InviteUserToTeam invite user to team
// @Summary      invite user to team
// @Description  invite user to team .
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      integer  true  "team id"
// @Param        input body      ds.InviteUserToTeamRequest  true "team id and login of invited user"
// @Success      200   {object}  ds.InviteUserToTeamResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /teams/{id}/invite [post]
func (a *API) InviteUserToTeam(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.InviteUserToTeamRequest, ds.InviteUserToTeamResponse]{
		serviceFunc:    a.appService.InviteUserToTeam,
		responseWriter: writeJsonResponse[*ds.InviteUserToTeamResponse],
		requestExtractor: func(r *http.Request, v *ds.InviteUserToTeamRequest) error {
			idRaw := r.PathValue("id")
			if idRaw == "" {
				return fmt.Errorf("team id was not provided")
			}

			id, err := strconv.ParseInt(idRaw, 10, 64)
			if err != nil {
				return err
			}

			err = extractJsonBody(r, v)
			if err != nil {
				return err
			}

			v.TeamId = id

			return nil
		},
		httpResponse: w,
		httpRequest:  r,
		api:          a,
	})
}
