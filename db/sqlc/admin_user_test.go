package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAdminUser(t *testing.T) AdminUser {
	role := createRandomRole(t)
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateAdminUserParams{
		Account:        util.RandomAccount(),
		FullName:       util.RandomName(),
		HashedPassword: hashedPassword,
		Status:         1,
		RoleID: pgtype.Int4{
			Int32: int32(role.ID),
			Valid: true,
		},
	}

	adminUser, err := testStore.CreateAdminUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, adminUser)

	require.Equal(t, arg.Account, adminUser.Account)
	require.Equal(t, arg.FullName, adminUser.FullName)
	require.Equal(t, arg.HashedPassword, adminUser.HashedPassword)
	require.Equal(t, arg.RoleID.Int32, adminUser.RoleID.Int32)

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
	require.Equal(t, adminUser1.RoleID.Int32, adminUser2.RoleID.Int32)
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
	require.Equal(t, oldAdminUser.RoleID.Int32, oldAdminUser.RoleID.Int32)
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
	require.Equal(t, oldAdminUser.RoleID.Int32, oldAdminUser.RoleID.Int32)
}

func TestUpdateAdminUserOnlyPassword(t *testing.T) {
	oldAdminUser := createRandomAdminUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
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
	require.Equal(t, oldAdminUser.RoleID.Int32, oldAdminUser.RoleID.Int32)
}

func TestUpdateAdminUserOnlyRoleId(t *testing.T) {
	role := createRandomRole(t)
	oldAdminUser := createRandomAdminUser(t)

	updateAdminUser, err := testStore.UpdateAdminUser(context.Background(), UpdateAdminUserParams{
		ID: oldAdminUser.ID,
		RoleID: pgtype.Int4{
			Int32: int32(role.ID),
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateAdminUser)

	require.NotEqual(t, oldAdminUser.RoleID.Int32, updateAdminUser.RoleID.Int32)
	require.Equal(t, int32(role.ID), updateAdminUser.RoleID.Int32)
	require.Equal(t, oldAdminUser.Account, updateAdminUser.Account)
	require.Equal(t, oldAdminUser.FullName, updateAdminUser.FullName)
	require.Equal(t, oldAdminUser.Status, updateAdminUser.Status)
	require.Equal(t, oldAdminUser.HashedPassword, oldAdminUser.HashedPassword)
}

func TestUpdateAdminUserAllFields(t *testing.T) {
	role := createRandomRole(t)
	oldAdminUser := createRandomAdminUser(t)

	newFullName := util.RandomName()
	newStatus := int32(0)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
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
		RoleID: pgtype.Int4{
			Int32: int32(role.ID),
			Valid: true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateAdminUser)

	require.Equal(t, oldAdminUser.Account, updateAdminUser.Account)
	require.NotEqual(t, oldAdminUser.FullName, updateAdminUser.FullName)
	require.NotEqual(t, oldAdminUser.Status, updateAdminUser.Status)
	require.NotEqual(t, oldAdminUser.HashedPassword, updateAdminUser.HashedPassword)
	require.NotEqual(t, oldAdminUser.RoleID.Int32, updateAdminUser.RoleID.Int32)

	require.Equal(t, newFullName, updateAdminUser.FullName)
	require.Equal(t, newStatus, updateAdminUser.Status)
	require.Equal(t, newHashedPassword, updateAdminUser.HashedPassword)
	require.Equal(t, int32(role.ID), updateAdminUser.RoleID.Int32)
}

func TestDeleteAdminUser(t *testing.T) {
	adminUser1 := createRandomAdminUser(t)

	err := testStore.DeleteAdminUser(context.Background(), adminUser1.ID)
	require.NoError(t, err)

	adminUser2, err := testStore.GetAdminUser(context.Background(), adminUser1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, adminUser2)
}

func TestListAdminUsers(t *testing.T) {
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

func TestGetAdminUsersCount(t *testing.T) {
	createRandomAdminUser(t)

	count, err := testStore.GetAdminUsersCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)
}

func TestGetAdminUserByAccount(t *testing.T) {
	adminUser1 := createRandomAdminUser(t)

	adminUser2, err := testStore.GetAdminUserByAccount(context.Background(), adminUser1.Account)
	require.NoError(t, err)
	require.NotEmpty(t, adminUser2)

	require.Equal(t, adminUser1.ID, adminUser2.ID)
	require.Equal(t, adminUser1.Account, adminUser2.Account)
	require.Equal(t, adminUser1.FullName, adminUser2.FullName)
	require.Equal(t, adminUser1.HashedPassword, adminUser2.HashedPassword)
	require.Equal(t, adminUser1.Status, adminUser2.Status)
	require.Equal(t, adminUser1.RoleID.Int32, adminUser2.RoleID.Int32)

	require.WithinDuration(t, adminUser1.PasswordChangedAt, adminUser2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, adminUser1.CreatedAt, adminUser2.CreatedAt, time.Second)
}
