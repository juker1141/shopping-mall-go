package db

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAdminUserTx(t *testing.T) AdminUserTxResult {
	account := util.RandomAccount()
	fullName := util.RandomName()
	status := int32(1)
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	var rolesID []int64

	n := 5
	for i := 0; i < n; i++ {
		roleTxResult := createRandomRoleTx(t)
		rolesID = append(rolesID, roleTxResult.Role.ID)
	}

	adArg := CreateAdminUserTxParams{
		Account:        account,
		FullName:       fullName,
		HashedPassword: hashedPassword,
		Status:         status,
		RolesID:        rolesID,
	}

	result, err := testStore.CreateAdminUserTx(context.Background(), adArg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	createdAdminUser := result.AdminUser
	require.NotEmpty(t, createdAdminUser)
	require.Equal(t, account, createdAdminUser.Account)
	require.Equal(t, fullName, createdAdminUser.FullName)
	require.Equal(t, hashedPassword, createdAdminUser.HashedPassword)
	require.Equal(t, status, createdAdminUser.Status)
	require.NotZero(t, createdAdminUser.CreatedAt)

	permissions := result.PermissionList

	targetLen := int(math.Pow(float64(n), 2))
	require.Len(t, permissions, targetLen)

	return result
}

func TestCreateAdminUserTx(t *testing.T) {
	createRandomAdminUserTx(t)
}

func TestUpdateAdminUserTx(t *testing.T) {
	newFullName := util.RandomName()
	newStatus := int32(0)
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	newPasswordChangeAt := time.Now()

	adminUserTx := createRandomAdminUserTx(t)
	n := 5
	rolesID := createRandomRoleList(t, n)

	result, err := testStore.UpdateAdminUserTx(context.Background(), UpdateAdminUserTxParams{
		ID:                adminUserTx.AdminUser.ID,
		FullName:          newFullName,
		HashedPassword:    newHashedPassword,
		PasswordChangedAt: newPasswordChangeAt,
		Status:            newStatus,
		RolesID:           rolesID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	updatedAdminUser := result.AdminUser
	require.NotEmpty(t, updatedAdminUser)
	require.NotEqual(t, adminUserTx.AdminUser.FullName, updatedAdminUser.FullName)
	require.NotEqual(t, adminUserTx.AdminUser.Status, updatedAdminUser.Status)
	require.NotEqual(t, adminUserTx.AdminUser.HashedPassword, updatedAdminUser.HashedPassword)

	require.NotZero(t, updatedAdminUser.PasswordChangedAt)
	require.NotZero(t, updatedAdminUser.CreatedAt)

	require.WithinDuration(t, newPasswordChangeAt, updatedAdminUser.PasswordChangedAt, time.Second)
}

func TestDeleteAdminUserTx(t *testing.T) {
	adminUserTx := createRandomAdminUserTx(t)

	result, err := testStore.DeleteAdminUserTx(context.Background(), DeleteAdminUserTxParams{
		ID: adminUserTx.AdminUser.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	adminUser, err := testStore.GetAdminUser(context.Background(), adminUserTx.AdminUser.ID)
	require.Error(t, err)
	require.Empty(t, adminUser)
}

func createRandomRoleList(t *testing.T, size int) []int64 {
	var rolesID []int64
	for i := 0; i < size; i++ {
		name := util.RandomRole()
		role, err := testStore.CreateRole(context.Background(), name)
		require.NoError(t, err)
		rolesID = append(rolesID, role.ID)
	}
	require.Len(t, rolesID, size)
	return rolesID
}
