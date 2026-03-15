package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
)

func (c *Client) AddNewTask(req *ds.DBCreateTaskRequest) (resp *ds.CreateTaskResponse, err error) {
	err = c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		if req.Role != ds.RoleAdmin {
			isOwner, err := qtx.IsTeamMember(ctx, sqlc.IsTeamMemberParams{
				UserID: req.CreatedBy,
				TeamID: req.TeamId,
			})
			if err != nil {
				return err
			}

			if !isOwner {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StaturNotOwner}}
				return nil
			}
		}

		res, err := qtx.AddNewTask(ctx, sqlc.AddNewTaskParams{
			AssigneeID:  req.AssigneeId,
			CreatedBy:   req.CreatedBy,
			TeamID:      req.TeamId,
			Subject:     req.Subject,
			Description: req.Description,
			Status:      string(req.Status),
		})
		if err != nil {
			if isLongData(err) {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusDataTooLong}}
				return nil
			}
			if isDuplicate(err) {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusAlreadyExists}}
				return nil
			}
			return err
		}

		n, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if n < 1 {
			return fmt.Errorf("no rows affected on AddNewTask.AddNewTask")
		}

		resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusSuccess}}
		return nil
	})

	return
}

func (c *Client) GetTasks(req *ds.GetTasksRequest) (resp *ds.GetTasksResponse, err error) {
	c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		if req.Role != ds.RoleAdmin {
			member, err := qtx.IsTeamMember(ctx, sqlc.IsTeamMemberParams{
				UserID: req.UserID,
				TeamID: req.TeamId,
			})
			if err != nil {
				return err
			}
			if !member {
				resp = &ds.GetTasksResponse{Status: ds.Status{Message: ds.StaturNotMember}}
			}
		}

		var status sql.NullString
		if req.Status != "" {
			status = nullString((*string)(&req.Status))
		}
		var assignee sql.NullInt64
		if req.AssigneeId > 0 {
			assignee = nullInt64(&req.AssigneeId)
		}
		tasks, err := qtx.GetTasks(ctx, sqlc.GetTasksParams{
			TeamID:     req.TeamId,
			Status:     status,
			AssigneeID: assignee,
			Limit:      int32(req.Limit),
			Offset:     int32(req.Offset),
		})
		if err != nil {
			return err
		}

		resp = &ds.GetTasksResponse{
			Tasks: make([]ds.Task, 0, len(tasks)),
		}
		for i := range tasks {
			resp.Tasks = append(resp.Tasks, ds.Task{
				TaskId:      tasks[i].TaskID.Int64,
				AssigneeId:  tasks[i].AssigneeID,
				CreatedBy:   tasks[i].CreatedBy,
				TeamId:      tasks[i].TeamID,
				Subject:     tasks[i].Subject,
				Description: tasks[i].Description,
				Status:      ds.TaskStatus(tasks[i].Status),
				CreatedAt:   tasks[i].CreatedAt.Time,
			})
		}
		resp.Status = ds.Status{Message: ds.StatusSuccess}
		return nil
	})

	return
}
