package mysql

import (
	"context"
	"fmt"
	"team-task-manager/internal/clients/mysql/sqlc"
	ds "team-task-manager/internal/datastruct"
)

func (c *Client) AddNewUser(req *ds.DBRegisterRequest) (*ds.RegisterResponse, error) {
	var resp *ds.RegisterResponse
	err := c.db.ExecTx(defaultTxOpt, func(ctx context.Context, qtx IQuerier) error {
		res, err := qtx.CreateUserAuth(ctx, sqlc.CreateUserAuthParams{
			Login:        req.Login,
			PasswordHash: req.Password,
		})
		if err != nil {
			if isDuplicate(err) {
				resp = &ds.RegisterResponse{Status: ds.Status{Message: ds.StatusUserAlreadyExists}}
			}
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
			UserID: id,
			Name:   req.Name,
			Role:   req.Role,
		})
		if err != nil {
			return err
		}

		resp = &ds.RegisterResponse{Status: ds.Status{Message: ds.StatusSuccess}}
		return nil
	})

	if resp != nil {
		return resp, nil
	}

	return nil, err
}

func (c *Client) GetAuthIdentitiesByUserID(id int64) (*ds.AuthIdentities, bool, error) {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	res, err := c.db.Querier().GetUserIdentitiesById(ctx, id)
	if err != nil {
		if isNoRows(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &ds.AuthIdentities{
		UserCreds: ds.UserCreds{
			Login:    res.Login,
			Password: res.PasswordHash,
		},
		UserInfo: ds.UserInfo{
			Name: res.Name,
		},
		JWTCreds: ds.JWTCreds{
			Role:   res.Role,
			UserID: res.ID,
		},
	}, true, nil
}

func (c *Client) GetAuthIdentitiesByLogin(login string) (*ds.AuthIdentities, bool, error) {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	res, err := c.db.Querier().GetUserIdentitiesByLogin(ctx, login)
	if err != nil {
		if isNoRows(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &ds.AuthIdentities{
		UserCreds: ds.UserCreds{
			Login:    res.Login,
			Password: res.PasswordHash,
		},
		UserInfo: ds.UserInfo{
			Name: res.Name,
		},
		JWTCreds: ds.JWTCreds{
			Role:   res.Role,
			UserID: res.ID,
		},
	}, true, nil
}

func (c *Client) AddRefreshToken(req *ds.DBRefreshToken) error {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	err := c.db.Querier().AddRefreshToken(ctx, sqlc.AddRefreshTokenParams{
		Token:     req.RefreshToken.RefreshToken,
		UserID:    req.UserID,
		ExpiredAt: req.ExpiresAt,
		Revoked:   req.Revoked,
		Used:      req.Used,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetRefreshToken(token string) (*ds.DBRefreshToken, bool, error) {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	res, err := c.db.Querier().GetRefreshToken(ctx, token)
	if err != nil {
		if isNoRows(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &ds.DBRefreshToken{
		RefreshToken: ds.RefreshToken{RefreshToken: res.Token},
		ExpiresAt:    res.ExpiredAt,
		UserID:       res.UserID,
		Revoked:      res.Revoked,
		Used:         res.Used,
	}, true, nil
}

func (c *Client) UpdateRefreshToken(req *ds.DBUpdateRefreshToken) error {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	res, err := c.db.Querier().UpdateRefreshToken(ctx, sqlc.UpdateRefreshTokenParams{
		ExpiredAt: req.ExpiresAt,
		Used:      req.Used,
		Revoked:   req.Revoked,
		Token:     req.RefreshToken.RefreshToken,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf(ds.StatusResurceNotFound)
	}

	return nil
}

func (c *Client) CleanupUslessTokens() error {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	err := c.db.Querier().CleanupUselessTokens(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteAllUserSession(userId int64) error {
	ctx, cancel := c.db.CtxWithCancel()
	defer cancel()

	res, err := c.db.Querier().DeleteUserSessions(ctx, userId)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf(ds.StatusResurceNotFound)
	}

	return nil
}
