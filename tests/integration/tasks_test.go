package integration

import (
	"encoding/json"
	"net/http"
	"strconv"
	ds "team-task-manager/internal/datastruct"
)

func (s *ServicesTestSuite) Test_Tasks_CreateTask() {
	name := "test_tasks_create_task"
	teamName := "test_team_create_task"
	_ = s.register(name)
	at, _ := s.login(name)
	_ = s.createTeam(teamName, at)

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", teamName).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	uri := apiPrefix + "/tasks"

	s.Run("Ok", func() {
		payload := map[string]any{
			"assignee_login": name,
			"subject":        name,
			"description":    name,
			"status":         "todo",
			"team_id":        teamId,
		}

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusOK, w.Code)
	})

	s.Run("Withou required field", func() {
		payload := map[string]any{
			// "assignee_login": name,
			"subject":     name,
			"description": name,
			"status":      "todo",
			"team_id":     teamId,
		}

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("With wrong status", func() {
		payload := map[string]any{
			"assignee_login": name,
			"subject":        name,
			"description":    name,
			"status":         "wrong",
			"team_id":        teamId,
		}

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("With wrong team id", func() {
		payload := map[string]any{
			"assignee_login": name,
			"subject":        name,
			"description":    name,
			"status":         "todo",
			"team_id":        int64(9999),
		}

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusForbidden, w.Code)
	})
}

func (s *ServicesTestSuite) Test_Tasks_GetTasks() {
	name := "test_tasks_get_tasks"
	teamName := "test_team_get_tasks"
	_ = s.register(name)
	at, _ := s.login(name)
	_ = s.createTeam(teamName, at)
	_ = s.createTaks(name, teamName, at)

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", teamName).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	uri := apiPrefix + "/tasks"

	s.Run("Ok", func() {
		query := map[string]string{
			"team_id":        strconv.Itoa(teamId),
			"status":         "todo",
			"assignee_login": name,
			"offset":         "0",
			"limit":          "100",
		}

		w := s.QueryRequest(http.MethodGet, query, nil, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusOK, w.Code)
	})

	s.Run("Without required field", func() {
		query := map[string]string{
			"team_id":        strconv.Itoa(teamId),
			"status":         "todo",
			"assignee_login": name,
		}

		w := s.QueryRequest(http.MethodGet, query, nil, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("With wrong team id", func() {
		query := map[string]string{
			"team_id":        "43556787",
			"status":         "todo",
			"assignee_login": name,
			"offset":         "0",
			"limit":          "100",
		}

		w := s.QueryRequest(http.MethodGet, query, nil, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})

		s.Equal(http.StatusForbidden, w.Code)
	})
}

func (s *ServicesTestSuite) Test_Tasks_UpdateTask() {
	name := "test_tasks_update_task"
	teamName := "test_team_update_task"
	_ = s.register(name)
	at, _ := s.login(name)
	_ = s.createTeam(teamName, at)
	_ = s.createTaks(name, teamName, at)

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", teamName).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	var taskId int
	err = s.db.QueryRow("SELECT task_id FROM tasks WHERE team_id = ?", teamId).Scan(&taskId)
	s.NoError(err)
	s.True(taskId > 0)

	uri := apiPrefix + "/tasks/" + strconv.Itoa(taskId)

	s.Run("Ok", func() {
		payload := map[string]any{
			"assignee_login": name,
			"subject":        "new one",
			"description":    "new one",
			"status":         "completed",
			"team_id":        teamId,
			"version":        1,
		}

		w := s.JSONBodyRequest(http.MethodPut, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusOK, w.Code)
	})

	s.Run("With no required field", func() {
		payload := map[string]any{
			"subject":     "new one",
			"description": "new one",
			"status":      "completed",
			"team_id":     teamId,
			"version":     2,
		}

		w := s.JSONBodyRequest(http.MethodPut, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusBadRequest, w.Code)
	})

	s.Run("With wrong team id", func() {
		payload := map[string]any{
			"assignee_login": name,
			"subject":        "new one",
			"description":    "new one",
			"status":         "completed",
			"team_id":        int64(99999),
			"version":        2,
		}

		w := s.JSONBodyRequest(http.MethodPut, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusForbidden, w.Code)
	})

	s.Run("with vrong version", func() {
		payload := map[string]any{
			"assignee_login": name,
			"subject":        "new one",
			"description":    "new one",
			"status":         "completed",
			"team_id":        teamId,
			"version":        1, // was updated and should be 2 now
		}

		w := s.JSONBodyRequest(http.MethodPut, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusConflict, w.Code)
	})
}

func (s *ServicesTestSuite) Test_Tasks_GetTaskHistory() {
	name := "test_tasks_get_task_history"
	teamName := "test_team_get_task_history"
	_ = s.register(name)
	at, _ := s.login(name)
	_ = s.createTeam(teamName, at)
	_ = s.createTaks(name, teamName, at)

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", teamName).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	var taskId int
	err = s.db.QueryRow("SELECT task_id FROM tasks WHERE team_id = ?", teamId).Scan(&taskId)
	s.NoError(err)
	s.True(taskId > 0)

	uri := apiPrefix + "/tasks/" + strconv.Itoa(taskId)

	payload := map[string]any{
		"assignee_login": name,
		"subject":        "new one",
		"description":    "new one",
		"status":         "completed",
		"team_id":        teamId,
		"version":        1,
	}

	w := s.JSONBodyRequest(http.MethodPut, payload, uri, [][2]string{
		{"Authorization", "Bearer " + at},
	})
	s.Equal(http.StatusOK, w.Code)

	uri = apiPrefix + "/tasks/" + strconv.Itoa(taskId) + "/history"

	s.Run("Ok", func() {
		query := map[string]string{
			"avoid_cache": "true",
		}

		w := s.QueryRequest(http.MethodGet, query, nil, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusOK, w.Code)

		v := ds.GetTaskHistoryResponse{}

		err := json.Unmarshal(w.Body.Bytes(), &v)
		s.Nil(err)
		s.False(v.CachedStatus.Cached)

		s.True(len(v.TaskHistory) > 1)

		s.Equal(v.TaskHistory[0].Status, ds.TaskStatus("todo"))
		s.Equal(v.TaskHistory[1].Status, ds.TaskStatus("completed"))
	})

	s.Run("With wrong taks id", func() {
		uri := apiPrefix + "/tasks/9999/history"
		query := map[string]string{
			"avoid_cache": "true",
		}

		w := s.QueryRequest(http.MethodGet, query, nil, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusNotFound, w.Code)
	})
}

func (s *ServicesTestSuite) Test_Tasks_AddTaskComment() {
	name := "test_tasks_get_task_history"
	teamName := "test_team_get_task_history"
	_ = s.register(name)
	at, _ := s.login(name)
	_ = s.createTeam(teamName, at)
	_ = s.createTaks(name, teamName, at)

	var teamId int
	err := s.db.QueryRow("SELECT team_id FROM teams WHERE name = ?", teamName).Scan(&teamId)
	s.NoError(err)
	s.True(teamId > 0)

	var taskId int
	err = s.db.QueryRow("SELECT task_id FROM tasks WHERE team_id = ?", teamId).Scan(&taskId)
	s.NoError(err)
	s.True(taskId > 0)

	s.Run("Ok", func() {
		uri := apiPrefix + "/tasks/" + strconv.Itoa(taskId) + "/comment"
		payload := map[string]any{
			"text": "Abracadabra, amor-ooh-na-na, morta-ooh-ga-ga, abra-ooh-na-na",
		}
		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusOK, w.Code)
	})

	s.Run("With wrong task id", func() {
		payload := map[string]any{
			"text": "Abracadabra, amor-ooh-na-na, morta-ooh-ga-ga, abra-ooh-na-na",
		}

		uri := apiPrefix + "/tasks/9999/comment"

		w := s.JSONBodyRequest(http.MethodPost, payload, uri, [][2]string{
			{"Authorization", "Bearer " + at},
		})
		s.Equal(http.StatusNotFound, w.Code)
	})
}
