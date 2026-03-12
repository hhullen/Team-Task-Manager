package api

import (
	"net/http"
)

const (
	teamsPrefix       = apiPrefix + "/teams"
	teamsInvitePrefix = teamsPrefix + "/{id}/invite"
)

func (a *API) setupTeamsHandlers() {
	a.router.Handle(pattern(http.MethodPost, teamsPrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.CreateTeam)))
}

// CreateTeam create new team
// @Summary      create new team
// @Description  create new team.
// @Tags         Teams
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input body      ds.RegisterRequest  true "login and password and info"
// @Success      200   {object}  ds.RegisterResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /teams [post]
func (a *API) CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Exec(ExecArgs[ds.RegisterRequest, ds.RegisterResponse]{
	// 	serviceFunc:      a.authService.RegisterUser,
	// 	responseWriter:   writeJsonResponse[*ds.RegisterResponse],
	// 	requestExtractor: extractJsonBody[*ds.RegisterRequest],
	// 	httpResponse:     &w,
	// 	httpRequest:      r,
	// 	api:              a,
	// })
}
