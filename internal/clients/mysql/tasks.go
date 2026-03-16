package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
)

func (c *Client) AddNewTask(req *ds.DBCreateTaskRequest) (*ds.CreateTaskResponse, error) {
	var resp *ds.CreateTaskResponse
	err := c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		if req.Role != ds.RoleAdmin {
			isOwner, err := qtx.IsTeamMember(ctx, sqlc.IsTeamMemberParams{
				UserID: req.CreatedBy,
				TeamID: req.TeamId,
			})
			if err != nil {
				return err
			}

			if !isOwner {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusNotOwner}}
				return interruptTxErr
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
			}
			if isDuplicate(err) {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusAlreadyExists}}
			}
			return err
		}

		if n, err := res.RowsAffected(); err != nil {
			return err
		} else if n < 1 {
			return fmt.Errorf("no rows affected on AddNewTask.AddNewTask")
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		v := ds.TaskUpdatePatch{
			TeamId:      supports.MakePatchFromTexts("", strconv.FormatInt(req.TeamId, 10)),
			AssigneeId:  supports.MakePatchFromTexts("", strconv.FormatInt(req.AssigneeId, 10)),
			Status:      supports.MakePatchFromTexts("", string(req.Status)),
			Subject:     supports.MakePatchFromTexts("", req.Subject),
			Description: supports.MakePatchFromTexts("", req.Description),
		}

		payload, err := json.Marshal(v)
		if err != nil {
			return err
		}

		res, err = qtx.AddChangeToTaskHistory(ctx, sqlc.AddChangeToTaskHistoryParams{
			TaskID:    id,
			ChangedBy: req.CreatedBy,
			Payload:   json.RawMessage(payload),
		})
		if err != nil {
			if isLongData(err) {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusDataTooLong}}
			}
			if isDuplicate(err) {
				resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusAlreadyExists}}
			}
			return err
		}

		if n, err := res.RowsAffected(); err != nil {
			return err
		} else if n < 1 {
			return fmt.Errorf("no rows affected on AddNewTask.AddChangeToTaskHistory")
		}

		resp = &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusSuccess}}
		return nil
	})

	if resp != nil {
		return resp, nil
	}

	return nil, err
}

func (c *Client) GetTasks(req *ds.GetTasksRequest) (*ds.GetTasksResponse, error) {
	var resp *ds.GetTasksResponse
	err := c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		if req.Role != ds.RoleAdmin {
			member, err := qtx.IsTeamMember(ctx, sqlc.IsTeamMemberParams{
				UserID: req.UserID,
				TeamID: req.TeamId,
			})
			if err != nil {
				return err
			}
			if !member {
				resp = &ds.GetTasksResponse{Status: ds.Status{Message: ds.StatusNotMember}}
				return interruptTxErr
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
			Tasks: make([]ds.TaskOutput, 0, len(tasks)),
		}
		for i := range tasks {
			resp.Tasks = append(resp.Tasks, ds.TaskOutput{
				TaskId:      tasks[i].TaskID,
				AssigneeId:  tasks[i].AssigneeID,
				CreatedBy:   tasks[i].CreatedBy,
				TeamId:      tasks[i].TeamID,
				Subject:     tasks[i].Subject,
				Description: tasks[i].Description,
				Status:      ds.TaskStatus(tasks[i].Status),
				CreatedAt:   tasks[i].CreatedAt.Time,
				Version:     tasks[i].Version,
			})
		}
		resp.Status = ds.Status{Message: ds.StatusSuccess}
		return nil
	})

	if resp != nil {
		return resp, nil
	}

	return nil, err
}

func (c *Client) UpdateTask(req *ds.DBUpdateTaskRequest) (*ds.UpdateTaskResponse, error) {
	var resp *ds.UpdateTaskResponse
	err := c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		oldTaks, err := qtx.GetTaskForUpdate(ctx, req.TaskId)
		if req.AssigneeId == oldTaks.AssigneeID &&
			req.Description == oldTaks.Description &&
			req.Status == ds.TaskStatus(oldTaks.Status) &&
			req.Subject == oldTaks.Subject &&
			req.TeamId == oldTaks.TeamID {
			resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusSuccess}}
			return nil
		}

		if err != nil {
			if isNoRows(err) {
				resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusResurceNotFound}}
			}
			return err
		}

		if req.Role != ds.RoleAdmin {
			member, err := qtx.IsTeamMember(ctx, sqlc.IsTeamMemberParams{
				UserID: req.UserID,
				TeamID: req.TeamId,
			})
			if err != nil {
				return err
			}

			if !member {
				resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusNotMember}}
				return interruptTxErr
			}

			if oldTaks.TeamID != req.TeamId {
				resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusForbidden}}
				return interruptTxErr
			}
		}

		if req.Version > oldTaks.Version {
			resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusIvalidVersion}}
			return interruptTxErr
		} else if req.Version < oldTaks.Version {
			resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusConflict}}
			return interruptTxErr
		} else {
			res, err := qtx.UpdateTask(ctx, sqlc.UpdateTaskParams{
				AssigneeID:  req.AssigneeId,
				TaskID:      req.TaskId,
				TeamID:      req.TeamId,
				Subject:     req.Subject,
				Status:      string(req.Status),
				Description: req.Description,
			})
			if err != nil {
				return err
			}

			if n, err := res.RowsAffected(); err != nil {
				return err
			} else if n < 1 {
				return fmt.Errorf("no rows affected on UpdateTask.UpdateTask")
			}

			v := ds.TaskUpdatePatch{
				TeamId: supports.MakePatchFromTexts(
					strconv.FormatInt(oldTaks.TeamID, 10),
					strconv.FormatInt(req.TeamId, 10)),
				AssigneeId: supports.MakePatchFromTexts(
					strconv.FormatInt(oldTaks.AssigneeID, 10),
					strconv.FormatInt(req.AssigneeId, 10)),
				Status: supports.MakePatchFromTexts(
					oldTaks.Status, string(req.Status)),
				Subject: supports.MakePatchFromTexts(
					oldTaks.Subject, req.Subject),
				Description: supports.MakePatchFromTexts(
					oldTaks.Description, req.Description),
			}

			payload, err := json.Marshal(v)
			if err != nil {
				return err
			}

			res, err = qtx.AddChangeToTaskHistory(ctx, sqlc.AddChangeToTaskHistoryParams{
				TaskID:    req.TaskId,
				ChangedBy: req.UserID,
				Payload:   json.RawMessage(payload),
			})
			if err != nil {
				if isLongData(err) {
					resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusDataTooLong}}
					return interruptTxErr
				}
				return err
			}

			if n, err := res.RowsAffected(); err != nil {
				return err
			} else if n < 1 {
				return fmt.Errorf("no rows affected on  AddNewTask.AddChangeToTaskHistory")
			}

			resp = &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusSuccess}}
			return nil
		}
	})

	if resp != nil {
		return resp, nil
	}

	return nil, err
}
