package service

import (
	"strconv"
	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
)

func (s *Service) AddNewTask(req *ds.CreateTaskRequest) *ds.CreateTaskResponse {
	const AvoidCache = false
	ident, exists, err := s.getAuthIdentitiesByLogin(req.AssigneeLogin, AvoidCache)

	if err != nil {
		s.logger.ErrorKV("GetTasks.getAuthIdentitiesByLogin", "error", err.Error())
		return nil
	}

	if !exists {
		return &ds.CreateTaskResponse{Status: ds.Status{Message: ds.StatusUserNotFound}}
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
		ident, exists, err := s.getAuthIdentitiesByLogin(req.AssigneeLogin, req.AvoidCache())
		if err != nil {
			s.logger.ErrorKV("GetTasks.getAuthIdentitiesByLogin", "error", err.Error())
			return nil
		}

		if !exists {
			return &ds.GetTasksResponse{Status: ds.Status{Message: ds.StatusUserNotFound}}
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

func (s *Service) UpdateTask(req *ds.UpdateTaskRequest) *ds.UpdateTaskResponse {
	const AvoidCache = false
	ident, exists, err := s.getAuthIdentitiesByLogin(req.AssigneeLogin, AvoidCache)
	if err != nil {
		s.logger.ErrorKV("GetTasks.getAuthIdentitiesByLogin", "error", err.Error())
		return nil
	}

	if !exists {
		return &ds.UpdateTaskResponse{Status: ds.Status{Message: ds.StatusUserNotFound}}
	}

	res, err := s.storageApp.UpdateTask(&ds.DBUpdateTaskRequest{
		JWTCreds:    req.JWTCreds,
		TaskId:      req.TaskId,
		AssigneeId:  ident.UserID,
		Subject:     req.Subject,
		Description: req.Description,
		Status:      req.Status,
		TeamId:      req.TeamId,
		Version:     req.Version,
	})
	if err != nil {
		s.logger.ErrorKV("UpdateTask.UpdateTask", "error", err.Error())
		return nil
	}

	return res
}

func (s *Service) GetTaskHistory(req *ds.GetTaskHistoryRequest) *ds.GetTaskHistoryResponse {
	key := makeCacheKey("TaskHistory", supports.FNV1Hash([]byte(strconv.FormatInt(req.TaskId, 10))))
	res, err := execWithCache(s, key, req.AvoidCache(), func() (*ds.GetTaskHistoryResponse, error) {
		return s.storageApp.GetTaskHistory(req)
	})
	if err != nil {
		s.logger.ErrorKV("GetTaskHistory.GetTaskHistory", "error", err.Error())
		return nil
	}

	return res
}
