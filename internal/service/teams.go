package service

import (
	ds "team-task-manager/internal/datastruct"
)

func (s *Service) CreateTeam(req *ds.CreateTeamRequest) *ds.CreateTeamResponse {
	res, err := s.storageApp.AddNewTeam(req)
	if err != nil {
		s.logger.ErrorKV("CreateTeam.AddNewTeam", "error", err.Error())
		return nil
	}

	return res
}

func (s *Service) ListUserTeams(req *ds.ListUserTeamsRequest) *ds.ListUserTeamsResponse {
	res, err := s.storageApp.GetUserTeams(req.UserID)
	if err != nil {
		s.logger.ErrorKV("ListUserTeams.GetUserTeams", "error", err.Error())
		return nil
	}
	return res
}

func (s *Service) InviteUserToTeam(req *ds.InviteUserToTeamRequest) *ds.InviteUserToTeamResponse {
	const AvoidCache = false
	ident, exists, err := s.getAuthIdentitiesByLogin(req.UserLoginToInvite, AvoidCache)
	if err != nil {
		s.logger.ErrorKV("InviteUserToTeam.GetAuthIdentitiesByLogin", "error", err.Error())
		return nil
	}

	if !exists {
		return &ds.InviteUserToTeamResponse{Status: ds.Status{Message: ds.StatusUserNotFound}}
	}

	res, err := s.storageApp.AddUserToUserTeam(&ds.DBInviteUserToTeamRequest{
		InviterId:      req.UserID,
		TeamId:         req.TeamId,
		UserIdToInvite: ident.UserID,
	})
	if err != nil {
		s.logger.ErrorKV("InviteUserToTeam.AddUserToUserTeam", "error", err.Error())
		return nil
	}
	return res
}
