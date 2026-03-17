package datastruct

import (
	"fmt"
	"strings"
	"time"
)

type TaskStatus string

const (
	Todo       = "todo"
	InProgress = "in progress"
	OnHold     = "on hold"
	Completed  = "completed"
	Canceled   = "canceled"
)

var taskStatuses = []string{
	Todo,
	InProgress,
	OnHold,
	Completed,
	Canceled,
}

func (g *TaskStatus) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	for i := range taskStatuses {
		if s == taskStatuses[i] {
			*g = TaskStatus(s)
			return nil
		}
	}

	return fmt.Errorf("incorrect task status: '%s'", s)
}

func (g *TaskStatus) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", *g), nil
}

type DBCreateTaskRequest struct {
	Role        string
	AssigneeId  int64
	CreatedBy   int64
	TeamId      int64
	Subject     string
	Description string
	Status      TaskStatus
}

type CreateTaskRequest struct {
	JWTCreds
	AssigneeLogin string     `json:"assignee_login" validate:"required" example:"VaKadyk359"`
	Subject       string     `json:"subject" validate:"required" example:"service endpoint"`
	Description   string     `json:"description" validate:"required" example:"add new service endpoint"`
	Status        TaskStatus `json:"status" validate:"required" example:"todo"`
	TeamId        int64      `json:"team_id" validate:"required" example:"1"`
}

type CreateTaskResponse struct {
	Status
}

type GetTasksRequest struct {
	JWTCreds
	AvoidCacheFlag
	TeamId        int64      `schema:"team_id" validate:"required" example:"1"`
	Status        TaskStatus `schema:"status" example:"todo"`
	AssigneeLogin string     `schema:"assignee_login" example:"VaKadyk359"`
	AssigneeId    int64      `schema:"assignee_id" example:"1"`
	Limit         int64      `schema:"limit" validate:"required" example:"10"`
	Offset        int64      `schema:"offset" example:"0"`
}

type TaskOutput struct {
	TaskId      int64      `json:"task_id"`
	Version     int64      `json:"version"`
	AssigneeId  int64      `json:"assignee_id"`
	CreatedBy   int64      `json:"created_by"`
	TeamId      int64      `json:"team_id"`
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

type GetTasksResponse struct {
	Status
	CachedStatus
	Tasks []TaskOutput
}

type TaskUpdatePatch struct {
	TeamId      string `json:"team_id,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	AssigneeId  string `json:"assignee_id,omitempty"`
}

type DBUpdateTaskRequest struct {
	JWTCreds
	TaskId      int64
	AssigneeId  int64
	Subject     string
	Description string
	Status      TaskStatus
	TeamId      int64
	Version     int64
}

type UpdateTaskRequest struct {
	JWTCreds
	TaskId        int64      `json:"task_id" validate:"required" example:"1"`
	AssigneeLogin string     `json:"assignee_login" validate:"required" example:"VaKadyk359"`
	Subject       string     `json:"subject" validate:"required" example:"service endpoint"`
	Description   string     `json:"description" validate:"required" example:"add new service endpoint"`
	Status        TaskStatus `json:"status" validate:"required" example:"todo"`
	TeamId        int64      `json:"team_id" validate:"required" example:"1"`
	Version       int64      `json:"version" validate:"required" example:"1"`
}

type UpdateTaskResponse struct {
	Status
}

type GetTaskHistoryRequest struct {
	JWTCreds
	AvoidCacheFlag
	TaskId int64 `uri:"task_id" validate:"required" example:"1"`
}

type GetTaskHistoryResponse struct {
	Status
	CachedStatus
	TaskHistory []TaskOutput
}
