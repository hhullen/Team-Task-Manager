package service

import (
	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().AddNewUser(gomock.Any()).
			Return(&ds.RegisterResponse{}, nil)

		req := &ds.RegisterRequest{}
		res := st.s.RegisterUser(req)
		require.NotNil(t, res)
	})

	t.Run("error on AddNewUser", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().AddNewUser(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.RegisterRequest{}
		res := st.s.RegisterUser(req)
		require.Nil(t, res)
	})
}

func TestLoginUser(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)
		userPass := "test"
		hashPass := supports.ArgonHash(userPass)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{
				UserCreds: ds.UserCreds{Password: hashPass},
			}, true, nil)
		st.authStore.EXPECT().AddRefreshToken(gomock.Any()).Return(nil)
		st.secret.EXPECT().ReadSecret(gomock.Any()).Return("sign", nil)

		req := &ds.LoginRequest{UserCreds: ds.UserCreds{
			Password: userPass,
		}}
		res := st.s.LoginUser(req)
		require.NotNil(t, res)
	})

	t.Run("error on GetAuthIdentitiesByLogin", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)
		userPass := "test"

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(nil, false, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.LoginRequest{UserCreds: ds.UserCreds{
			Password: userPass,
		}}
		res := st.s.LoginUser(req)
		require.Nil(t, res)
	})

	t.Run("Not exists", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)
		userPass := "test"

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(nil, false, nil)

		req := &ds.LoginRequest{UserCreds: ds.UserCreds{
			Password: userPass,
		}}
		res := st.s.LoginUser(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusWrongLoginOrPassword)
	})

	t.Run("Wrong password", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)
		userPass := "test"
		hashPass := supports.ArgonHash(userPass)

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{
				UserCreds: ds.UserCreds{Password: hashPass},
			}, true, nil)

		req := &ds.LoginRequest{UserCreds: ds.UserCreds{
			Password: "wrong password",
		}}
		res := st.s.LoginUser(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusWrongLoginOrPassword)
	})
}

func TestRefresh(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(time.Hour * 2),
				Revoked:   false,
				Used:      false,
			}, true, nil)
		st.authStore.EXPECT().GetAuthIdentitiesByUserID(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)
		st.authStore.EXPECT().UpdateRefreshToken(gomock.Any()).
			Return(nil)

		st.authStore.EXPECT().AddRefreshToken(gomock.Any()).Return(nil)
		st.secret.EXPECT().ReadSecret(gomock.Any()).Return("sign", nil)

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.NotNil(t, res)
	})

	t.Run("error on GetRefreshToken", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(nil, false, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.Nil(t, res)
	})

	t.Run("Not exists on GetRefreshToken", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(nil, false, nil)

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusInvalidToken)
	})

	t.Run("Revoked", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(time.Hour * 2),
				Revoked:   true,
				Used:      false,
			}, true, nil)

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusInvalidToken)
	})

	t.Run("Expired", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(-time.Hour * 2),
				Revoked:   false,
				Used:      false,
			}, true, nil)

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusInvalidToken)
	})

	t.Run("error on GetAuthIdentitiesByUserID", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(time.Hour * 2),
				Revoked:   false,
				Used:      false,
			}, true, nil)
		st.authStore.EXPECT().GetAuthIdentitiesByUserID(gomock.Any()).
			Return(nil, true, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.Nil(t, res)
	})

	t.Run("Not exists on GetAuthIdentitiesByUserID", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(time.Hour * 2),
				Revoked:   false,
				Used:      false,
			}, true, nil)
		st.authStore.EXPECT().GetAuthIdentitiesByUserID(gomock.Any()).
			Return(nil, false, nil)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.Nil(t, res)
	})

	t.Run("suspicion of token theft", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(-time.Hour * 2),
				Revoked:   false,
				Used:      true,
			}, true, nil)
		st.authStore.EXPECT().GetAuthIdentitiesByUserID(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)
		st.authStore.EXPECT().DeleteAllUserSession(gomock.Any()).Return(nil)

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusSessionReset)
	})

	t.Run("error on DeleteAllUserSession of suspicion of token theft", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(-time.Hour * 2),
				Revoked:   false,
				Used:      true,
			}, true, nil)
		st.authStore.EXPECT().GetAuthIdentitiesByUserID(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)
		st.authStore.EXPECT().DeleteAllUserSession(gomock.Any()).Return(errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.Nil(t, res)
	})

	t.Run("error on UpdateRefreshToken", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().GetRefreshToken(gomock.Any()).
			Return(&ds.DBRefreshToken{
				ExpiresAt: time.Now().Add(time.Hour * 2),
				Revoked:   false,
				Used:      false,
			}, true, nil)
		st.authStore.EXPECT().GetAuthIdentitiesByUserID(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)
		st.authStore.EXPECT().UpdateRefreshToken(gomock.Any()).
			Return(errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.RefreshRequest{
			RefreshToken: ds.RefreshToken{RefreshToken: "token"},
		}
		res := st.s.Refresh(req)
		require.Nil(t, res)
	})
}

func TestCreateTokensPair(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().AddRefreshToken(gomock.Any()).
			Return(nil)
		st.secret.EXPECT().ReadSecret(gomock.Any()).Return("sign", nil)

		at, rt, err := st.s.createTokensPair(&ds.AuthIdentities{
			JWTCreds: ds.JWTCreds{
				Role:   "user",
				UserID: 12,
			},
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
			UserInfo: ds.UserInfo{
				Name: "name",
			},
		})
		require.Nil(t, err)
		require.NotEqual(t, at, "")
		require.NotEqual(t, rt, "")
	})

	t.Run("error on AddRefreshToken", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().AddRefreshToken(gomock.Any()).
			Return(errTest)

		_, _, err := st.s.createTokensPair(&ds.AuthIdentities{
			JWTCreds: ds.JWTCreds{
				Role:   "user",
				UserID: 12,
			},
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
			UserInfo: ds.UserInfo{
				Name: "name",
			},
		})
		require.NotNil(t, err)
	})

	t.Run("error on ReadSecret", func(t *testing.T) {
		t.Parallel()
		st := newServiceTest(t)

		st.authStore.EXPECT().AddRefreshToken(gomock.Any()).
			Return(nil)
		st.secret.EXPECT().ReadSecret(gomock.Any()).Return("", errTest)

		_, _, err := st.s.createTokensPair(&ds.AuthIdentities{
			JWTCreds: ds.JWTCreds{
				Role:   "user",
				UserID: 12,
			},
			UserCreds: ds.UserCreds{
				Login:    "login",
				Password: "password",
			},
			UserInfo: ds.UserInfo{
				Name: "name",
			},
		})
		require.NotNil(t, err)

	})
}
