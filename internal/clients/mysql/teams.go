package mysql

import (
	"context"
	"fmt"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
)

func (c *Client) AddNewTeam(req *ds.CreateTeamRequest) (resp *ds.CreateTeamResponse, err error) {
	err = c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		res, err := qtx.AddNewTeam(ctx, sqlc.AddNewTeamParams{
			OwnerID:     req.UserID,
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			if isForeignKeyErr(err) {
				resp = &ds.CreateTeamResponse{Status: ds.Status{Message: ds.StatusResurceNotFound}}
			}
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		res, err = qtx.AddMemberToTeam(ctx, sqlc.AddMemberToTeamParams{
			UserID: req.UserID,
			TeamID: id,
		})
		if err != nil {
			if isForeignKeyErr(err) {
				resp = &ds.CreateTeamResponse{Status: ds.Status{Message: ds.StatusResurceNotFound}}
			}
			return err
		}

		n, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if n != 1 {
			return fmt.Errorf("no rows affected on AddNewTeam.AddMemberToTeam")
		}

		resp = &ds.CreateTeamResponse{
			Status: ds.Status{Message: ds.StatusSuccess},
		}

		return nil
	})

	if resp != nil {
		return resp, nil
	}

	return nil, err
}

func (c *Client) GetUserTeams(userId int64) (*ds.ListUserTeamsResponse, error) {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	res, err := c.db.Querier().GetUserTeams(ctx, userId)
	if err != nil {
		if isNoRows(err) {
			return &ds.ListUserTeamsResponse{Status: ds.Status{Message: ds.StatusSuccess}}, nil
		}
		return nil, err
	}

	resp := make([]ds.TeamDescription, len(res))
	for i := range res {
		resp[i] = ds.TeamDescription{
			Name:        res[i].Name,
			Description: res[i].Description,
			TeamId:      res[i].TeamID,
		}
	}

	return &ds.ListUserTeamsResponse{
		Status: ds.Status{Message: ds.StatusSuccess},
		List:   resp,
	}, nil
}

func (c *Client) AddUserToUserTeam(req *ds.DBInviteUserToTeamRequest) (*ds.InviteUserToTeamResponse, error) {
	var resp *ds.InviteUserToTeamResponse
	err := c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		id, err := qtx.GetTeamOwner(ctx, req.TeamId)
		if err != nil {
			if isNoRows(err) {
				resp = &ds.InviteUserToTeamResponse{Status: ds.Status{Message: ds.StatusUserNotFound}}
			}
			return err
		}

		if id != req.InviterId {
			resp = &ds.InviteUserToTeamResponse{Status: ds.Status{Message: ds.StatusNotOwner}}
			return interruptTxErr
		}

		res, err := qtx.AddMemberToTeam(ctx, sqlc.AddMemberToTeamParams{
			UserID: req.UserIdToInvite,
			TeamID: req.TeamId,
		})
		if err != nil {
			if isDuplicate(err) {
				resp = &ds.InviteUserToTeamResponse{Status: ds.Status{Message: ds.StatusUserAlreadyExists}}
			}
			return err
		}

		if n, err := res.RowsAffected(); err != nil {
			return err
		} else if n != 1 {
			return fmt.Errorf("no rows affected on AddUserToUserTeam.AddMemberToTeam")
		}

		resp = &ds.InviteUserToTeamResponse{Status: ds.Status{Message: ds.StatusSuccess}}
		return nil
	})

	if resp != nil {
		return resp, nil
	}

	return nil, err
}
