package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAdminUser(t *testing.T) AdminUser {
	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateAdminUserParams{
		Account:        util.RandomAccount(),
		FullName:       util.RandomName(),
		HashedPassword: hashedPassword,
		Status:         1,
	}

	adminUser, err := testStore.CreateAdminUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, adminUser)

	require.Equal(t, arg.Account, adminUser.Account)
	require.Equal(t, arg.FullName, adminUser.FullName)
	require.Equal(t, arg.HashedPassword, adminUser.HashedPassword)

	require.True(t, adminUser.PasswordChangedAt.IsZero())
	require.NotZero(t, adminUser.CreatedAt)

	return adminUser
}

func TestCreateAdminUser(t *testing.T) {
	createRandomAdminUser(t)
}

func TestGetAdminUser(t *testing.T) {
	adminUser1 := createRandomAdminUser(t)

	adminUser2, err := testStore.GetAdminUser(context.Background(), adminUser1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, adminUser2)

	require.Equal(t, adminUser1.ID, adminUser2.ID)
	require.Equal(t, adminUser1.Account, adminUser2.Account)
	require.Equal(t, adminUser1.FullName, adminUser2.FullName)
	require.Equal(t, adminUser1.HashedPassword, adminUser2.HashedPassword)
	require.Equal(t, adminUser1.Status, adminUser2.Status)

	require.WithinDuration(t, adminUser1.PasswordChangedAt, adminUser2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, adminUser1.CreatedAt, adminUser2.CreatedAt, time.Second)
}

func TestUpdateAdminUserOnlyFullName(t *testing.T) {
	oldAdminUser := createRandomAdminUser(t)

	newFullName := util.RandomName()
	updateAdminUser, err := testStore.UpdateAdminUser(context.Background(), UpdateAdminUserParams{
		ID: oldAdminUser.ID,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updateAdminUser)

	require.NotEqual(t, oldAdminUser.FullName, updateAdminUser.FullName)
	require.Equal(t, newFullName, updateAdminUser.FullName)
	require.Equal(t, oldAdminUser.Account, updateAdminUser.Account)
	require.Equal(t, oldAdminUser.HashedPassword, updateAdminUser.HashedPassword)
	require.Equal(t, oldAdminUser.Status, updateAdminUser.Status)
}

func TestUpdateAdminUserOnlyStatus(t *testing.T) {
	oldAdminUser := createRandomAdminUser(t)

	newStatus := int32(0)
	updateAdminUser, err := testStore.UpdateAdminUser(context.Background(), UpdateAdminUserParams{
		ID: oldAdminUser.ID,
		Status: pgtype.Int4{
			Int32: newStatus,
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateAdminUser)

	require.NotEqual(t, oldAdminUser.Status, updateAdminUser.Status)
	require.Equal(t, newStatus, updateAdminUser.Status)
	require.Equal(t, oldAdminUser.Account, updateAdminUser.Account)
	require.Equal(t, oldAdminUser.FullName, updateAdminUser.FullName)
	require.Equal(t, oldAdminUser.HashedPassword, updateAdminUser.HashedPassword)
}

func TestUpdateAdminUserOnlyPassword(t *testing.T) {
	oldAdminUser := createRandomAdminUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashedPassword(newPassword)
	require.NoError(t, err)

	updateAdminUser, err := testStore.UpdateAdminUser(context.Background(), UpdateAdminUserParams{
		ID: oldAdminUser.ID,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateAdminUser)

	require.NotEqual(t, oldAdminUser.HashedPassword, updateAdminUser.HashedPassword)
	require.Equal(t, newHashedPassword, updateAdminUser.HashedPassword)
	require.Equal(t, oldAdminUser.Account, updateAdminUser.Account)
	require.Equal(t, oldAdminUser.FullName, updateAdminUser.FullName)
	require.Equal(t, oldAdminUser.Status, updateAdminUser.Status)
}

func TestUpdateAdminUserAllFields(t *testing.T) {
	oldAdminUser := createRandomAdminUser(t)

	newFullName := util.RandomName()
	newStatus := int32(0)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashedPassword(newPassword)
	require.NoError(t, err)

	updateAdminUser, err := testStore.UpdateAdminUser(context.Background(), UpdateAdminUserParams{
		ID: oldAdminUser.ID,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Status: pgtype.Int4{
			Int32: newStatus,
			Valid: true,
		},
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateAdminUser)

	require.Equal(t, oldAdminUser.Account, updateAdminUser.Account)
	require.NotEqual(t, oldAdminUser.FullName, updateAdminUser.FullName)
	require.NotEqual(t, oldAdminUser.Status, updateAdminUser.Status)
	require.NotEqual(t, oldAdminUser.HashedPassword, updateAdminUser.HashedPassword)

	require.Equal(t, newFullName, updateAdminUser.FullName)
	require.Equal(t, newStatus, updateAdminUser.Status)
	require.Equal(t, newHashedPassword, updateAdminUser.HashedPassword)
}

func TestDeleteAdminUser(t *testing.T) {
	adminUser1 := createRandomAdminUser(t)

	err := testStore.DeleteAdminUser(context.Background(), adminUser1.ID)
	require.NoError(t, err)

	adminUser2, err := testStore.GetAdminUser(context.Background(), adminUser1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, adminUser2)
}

func TestListAdminUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAdminUser(t)
	}

	arg := ListAdminUsersParams{
		Limit:  5,
		Offset: 5,
	}

	adminUsers, err := testStore.ListAdminUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, adminUsers, 5)

	for _, adminUser := range adminUsers {
		require.NotEmpty(t, adminUser)
	}
}