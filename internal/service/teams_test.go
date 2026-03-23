package service

import (
	ds "team-task-manager/internal/datastruct"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTeam(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().AddNewTeam(gomock.Any()).
			Return(&ds.CreateTeamResponse{}, nil)

		req := &ds.CreateTeamRequest{}
		res := st.s.CreateTeam(req)
		require.NotNil(t, res)
	})

	t.Run("error on CreateTeam", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().AddNewTeam(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.Any())

		req := &ds.CreateTeamRequest{}
		res := st.s.CreateTeam(req)
		require.Nil(t, res)
	})
}

func TestListUserTeams(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().GetUserTeams(gomock.Any()).
			Return(&ds.ListUserTeamsResponse{}, nil)

		req := &ds.ListUserTeamsRequest{}
		res := st.s.ListUserTeams(req)
		require.NotNil(t, res)
	})

	t.Run("error on ListUserTeams", func(t *testing.T) {
		st := newServiceTest(t)

		st.appStore.EXPECT().GetUserTeams(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.Any())

		req := &ds.ListUserTeamsRequest{}
		res := st.s.ListUserTeams(req)
		require.Nil(t, res)
	})
}

func TestInviteUserToTeam(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().AddUserToUserTeam(gomock.Any()).
			Return(&ds.InviteUserToTeamResponse{}, nil)

		req := &ds.InviteUserToTeamRequest{}
		res := st.s.InviteUserToTeam(req)
		require.NotNil(t, res)
	})

	t.Run("Not exists", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})

		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(nil, false, nil)

		req := &ds.InviteUserToTeamRequest{}
		res := st.s.InviteUserToTeam(req)
		require.NotNil(t, res)
		require.Equal(t, res.GetStatus(), ds.StatusUserNotFound)
	})

	t.Run("error on InviteUserToTeam", func(t *testing.T) {
		st := newServiceTest(t)

		st.cache.EXPECT().Read(gomock.Any(), gomock.Any()).
			DoAndReturn(func(key string, v any) (bool, error) {
				return false, nil
			})
		st.cache.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
		st.authStore.EXPECT().GetAuthIdentitiesByLogin(gomock.Any()).
			Return(&ds.AuthIdentities{}, true, nil)

		st.appStore.EXPECT().AddUserToUserTeam(gomock.Any()).
			Return(nil, errTest)
		st.logger.EXPECT().ErrorKV(gomock.Any(), gomock.All())

		req := &ds.InviteUserToTeamRequest{}
		res := st.s.InviteUserToTeam(req)
		require.Nil(t, res)
	})
}
