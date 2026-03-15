package service

import (
	"strconv"
	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
)

func (s *Service) AddNewTask(req *ds.CreateTaskRequest) *ds.CreateTaskResponse {
	ident, exist, err := s.storageAuth.GetAuthIdentitiesByLogin(req.AssigneeLogin)
	if err != nil {
		s.logger.ErrorKV("AddNewTask.GetAuthIdentitiesByLogin", "error", err.Error())
		return nil
	}

	if !exist {
		return &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusNotFound}}
	}

	res, err := s.storageApp.AddNewTask(&ds.DBCreateTaskRequest{
		Role:        req.Role,
		AssigneeId:  ident.UserID,
		CreatedBy:   req.UserID,
		TeamId:      req.TeamId,
		Subject:     req.Subject,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		s.logger.ErrorKV("AddNewTask.AddNewTask", "error", err.Error())
		return nil
	}

	return res
}

func (s *Service) GetTasks(req *ds.GetTasksRequest) *ds.GetTasksResponse {
	if req.AssigneeId < 1 {
		ident, exist, err := s.storageAuth.GetAuthIdentitiesByLogin(req.AssigneeLogin)
		if err != nil {
			s.logger.ErrorKV("GetTasks.GetAuthIdentitiesByLogin", "error", err.Error())
			return nil
		}
		if !exist {
			return &ds.GetTasksResponse{Status: ds.Status{Message: ds.StatusNotFound}}
		}
		req.AssigneeId = ident.UserID
	}

	reqKey := supports.Concat(req.Role, strconv.FormatInt(req.UserID, 10),
		string(req.Status), req.AssigneeLogin,
		strconv.FormatInt(req.AssigneeId, 10),
		strconv.FormatInt(req.Limit, 10),
		strconv.FormatInt(req.Offset, 10),
		strconv.FormatInt(req.TeamId, 10))
	key := makeCacheKey("GetTasks", supports.FNV1Hash([]byte(reqKey)))

	res, err := execWithCache(s, key, req.AvoidCache(), func() (*ds.GetTasksResponse, error) {
		return s.storageApp.GetTasks(req)
	})
	if err != nil {
		s.logger.ErrorKV("GetTasks.GetTasks", "error", err.Error())
		return nil
	}

	return res
}
