package api

import (
	"fmt"
	"net/http"
	"strconv"
	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
)

const (
	tasksPrefix    = apiPrefix + "/tasks"
	taskByIdPrefix = tasksPrefix + "/{id}"
	taskHistory    = taskByIdPrefix + "/history"
	taskComment    = taskByIdPrefix + "/comment"
)

func (a *API) setupTasksHandlers() {
	a.router.Handle(pattern(http.MethodPost, tasksPrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.CreateTask)))
	a.router.Handle(pattern(http.MethodGet, tasksPrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.GetTasks)))
	a.router.Handle(pattern(http.MethodPut, taskByIdPrefix),
		jwtBasedMiddleware(a, http.HandlerFunc(a.UpdateTask)))
	a.router.Handle(pattern(http.MethodGet, taskHistory),
		jwtBasedMiddleware(a, http.HandlerFunc(a.GetTaskHistory)))
	a.router.Handle(pattern(http.MethodPost, taskComment),
		jwtBasedMiddleware(a, http.HandlerFunc(a.AddTaskComment)))
}

// CreateTask create new task
// @Summary      create new task
// @Description  create new task.
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input body      ds.CreateTaskRequest  true "task"
// @Success      200   {object}  ds.CreateTaskResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /tasks [post]
func (a *API) CreateTask(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.CreateTaskRequest, ds.CreateTaskResponse]{
		serviceFunc:      a.appService.AddNewTask,
		responseWriter:   writeJsonResponse[*ds.CreateTaskResponse],
		requestExtractor: extractJsonBody[*ds.CreateTaskRequest],
		httpResponse:     &w,
		httpRequest:      r,
		api:              a,
	})
}

// GetTasks get tasks
// @Summary      get tasks
// @Description  get tasks.
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        team_id        query     string              true  "team_id"        example(1)
// @Param        status         query     string              false  "status"         example(todo)
// @Param        assignee_id    query     string              false "assignee_id"    example(1)
// @Param        assignee_login query     string              false "assignee_login" example(VaKadyk359)
// @Param        offset         query     string              true  "offset"         example(0)
// @Param        limit          query     string              true  "limit"          example(10)
// @Param        avoid_cache    query     string              false "avoid_cache"    example(true)
// @Success      200            {object}  ds.GetTasksResponse
// @Failure      400            {object}  ds.Status
// @Failure      500            {object}  ds.Status
// @Router       /tasks  [get]
func (a *API) GetTasks(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.GetTasksRequest, ds.GetTasksResponse]{
		serviceFunc:      a.appService.GetTasks,
		responseWriter:   writeJsonResponse[*ds.GetTasksResponse],
		requestExtractor: extractSchemaQuery[*ds.GetTasksRequest],
		validator: func(v *ds.GetTasksRequest) error {
			if v.AssigneeId == 0 && v.AssigneeLogin == "" {
				return fmt.Errorf("assignee_id and assignee_login empty but expected at least one")
			}
			return supports.StructValidator().Struct(v)
		},
		httpResponse: &w,
		httpRequest:  r,
		api:          a,
	})
}

// UpdateTask update task
// @Summary      update task
// @Description  update task.
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      integer  true  "task id"
// @Param        input body      ds.UpdateTaskRequest  true "task"
// @Success      200   {object}  ds.UpdateTaskResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /tasks/{id} [put]
func (a *API) UpdateTask(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.UpdateTaskRequest, ds.UpdateTaskResponse]{
		serviceFunc:    a.appService.UpdateTask,
		responseWriter: writeJsonResponse[*ds.UpdateTaskResponse],
		requestExtractor: func(r *http.Request, v *ds.UpdateTaskRequest) error {
			idRaw := r.PathValue("id")
			if idRaw == "" {
				return fmt.Errorf("task id was not provided")
			}

			id, err := strconv.ParseInt(idRaw, 10, 64)
			if err != nil {
				return err
			}

			v.TaskId = id
			return extractJsonBody(r, v)
		},
		httpResponse: &w,
		httpRequest:  r,
		api:          a,
	})
}

// GetTaskHistory get tasks history
// @Summary      get tasks history
// @Description  get tasks history.
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id             path      integer                   true  "task_id"     example(1)
// @Param        avoid_cache    query     string                    false "avoid_cache" example(true)
// @Success      200            {object}  ds.GetTaskHistoryResponse
// @Failure      400            {object}  ds.Status
// @Failure      500            {object}  ds.Status
// @Router       /tasks/{id}/history  [get]
func (a *API) GetTaskHistory(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.GetTaskHistoryRequest, ds.GetTaskHistoryResponse]{
		serviceFunc:    a.appService.GetTaskHistory,
		responseWriter: writeJsonResponse[*ds.GetTaskHistoryResponse],
		requestExtractor: func(r *http.Request, v *ds.GetTaskHistoryRequest) error {
			idRaw := r.PathValue("id")
			if idRaw == "" {
				return fmt.Errorf("task id was not provided")
			}

			id, err := strconv.ParseInt(idRaw, 10, 64)
			if err != nil {
				return err
			}

			v.TaskId = id
			return extractSchemaQuery(r, v)
		},
		httpResponse: &w,
		httpRequest:  r,
		api:          a,
	})
}

// AddTaskComment add task comment
// @Summary      add task comment
// @Description  add task comment.
// @Tags         Tasks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      integer  true  "task_id"
// @Param        input body      ds.AddTaskCommentRequest  true "comment"
// @Success      200   {object}  ds.AddTaskCommentResponse
// @Failure      400   {object}  ds.Status
// @Failure      500   {object}  ds.Status
// @Router       /tasks/{id}/comment [post]
func (a *API) AddTaskComment(w http.ResponseWriter, r *http.Request) {
	Exec(ExecArgs[ds.AddTaskCommentRequest, ds.AddTaskCommentResponse]{
		serviceFunc:    a.appService.AddTaskComment,
		responseWriter: writeJsonResponse[*ds.AddTaskCommentResponse],
		requestExtractor: func(r *http.Request, v *ds.AddTaskCommentRequest) error {
			idRaw := r.PathValue("id")
			if idRaw == "" {
				return fmt.Errorf("task id was not provided")
			}

			id, err := strconv.ParseInt(idRaw, 10, 64)
			if err != nil {
				return err
			}

			v.TaskId = id
			return extractJsonBody(r, v)
		},
		httpResponse: &w,
		httpRequest:  r,
		api:          a,
	})
}
