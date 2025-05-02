package db

import (
	"context"
	"testing"
	"time"

	"github.com/shafi21064/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createTestUsers(t *testing.T) User {
	arg := CreateUsersParams{
		UserName:       util.RandomOwnerName(),
		HassedPassword: "secret",
		FullName:       util.RandomOwnerName(),
		Email:          util.RandomEmail(),
	}

	user, err := testQuaries.CreateUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.HassedPassword, user.HassedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUsers(t)
}

func TestGetUser(t *testing.T) {
	user1 := createTestUsers(t)

	user2, err := testQuaries.GetUsers(context.Background(), user1.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt.Time, user2.PasswordChangedAt.Time, time.Second)
}
