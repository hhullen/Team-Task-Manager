package integration_test

import (
	"encoding/json"
	"net/http"
	"strconv"

	ds "team-task-manager/internal/datastruct"
)

func (s *ServicesTestSuite) TestTeamsCreateTeam() {
	_ = s.register("test1")
	at, _ := s.login("test1")

	uri := apiPrefix + "/teams"

	s.Run("Ok", func() {
		payload := map[string]string{
			"name":        "name",
			"description": "desc",
		}

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusOK, w.Code)

		w = s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusOK, w.Code)
	})

	s.Run("Without required field", func() {
		payload := map[string]string{
			"description": "desc",
		}

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusBadRequest, w.Code)
	})
}

func (s *ServicesTestSuite) TestTeamsListUserTeams() {
	_ = s.register("test1")
	at, _ := s.login("test1")

	teamName := "name"
	_ = s.createTeam(teamName, at)

	uri := apiPrefix + "/teams"

	s.Run("Ok", func() {
		w := s.JSONBodyRequest(http.MethodGet, map[string]string{}, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusOK, w.Code)

		v := ds.ListUserTeamsResponse{}

		err := json.Unmarshal(w.Body.Bytes(), &v)
		s.Nil(err)

		s.True(len(v.List) > 0)

		var name, desc string
		for _, v := range v.List {
			if v.Name == teamName && v.Description == teamName {
				name = v.Name
				desc = v.Description
			}
		}
		s.Equal(name, teamName)
		s.Equal(desc, teamName)

	})
}

func (s *ServicesTestSuite) TestTeamsInviteUserToTeam() {
	s.register("test1")
	s.register("test2")
	at, _ := s.login("test1")
	newTeam := "New team"
	_ = s.createTeam(newTeam, at)

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", newTeam).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	uri := apiPrefix + "/teams/" + strconv.Itoa(teamId) + "/invite"

	s.Run("Ok", func() {
		w := s.JSONBodyRequest(http.MethodPost, map[string]string{"login": "test2"}, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusOK, w.Code)
	})

	s.Run("Wrong Login", func() {
		w := s.JSONBodyRequest(http.MethodPost, map[string]string{"login": "thisIsWrongMadafakas"}, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusNotFound, w.Code)
	})
}
