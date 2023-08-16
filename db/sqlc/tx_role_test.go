package db

import (
	"context"
	"testing"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomRoleTx(t *testing.T) CreateRoleTxResult {
	roleName := util.RandomName()

	n := 5
	permissionsID := CreateRandomPermissionList(t, n)

	result, err := testStore.CreateRoleTx(context.Background(), CreateRoleTxParams{
		Name:          roleName,
		PermissionsID: permissionsID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	role := result.Role
	require.NotEmpty(t, role)
	require.Equal(t, roleName, role.Name)
	require.NotZero(t, role.CreatedAt)
	permissions := result.PermissionList
	require.Len(t, permissions, n)
	return result
}

func TestCreateRoleTx(t *testing.T) {
	createRandomRoleTx(t)
}

func TestUpdateRoleTx(t *testing.T) {
	roleTx := createRandomRoleTx(t)
	n := 5
	permissionsID := CreateRandomPermissionList(t, n)

	result, err := testStore.UpdateRoleTx(context.Background(), UpdateRoleTxParams{
		ID:            roleTx.Role.ID,
		Name:          util.RandomName(),
		PermissionsID: permissionsID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	updatedRole := result.Role
	require.NotEmpty(t, updatedRole)
	require.NotEqual(t, roleTx.Role.Name, updatedRole.Name)
	require.NotZero(t, updatedRole.CreatedAt)
	permissions := result.PermissionList
	require.Len(t, permissions, n)
}

func TestDeleteRoleTx(t *testing.T) {
	roleTx := createRandomRoleTx(t)

	result, err := testStore.DeleteRoleTx(context.Background(), DeleteRoleTxParams{
		ID: roleTx.Role.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	role, err := testStore.GetRole(context.Background(), roleTx.Role.ID)
	require.Error(t, err)
	require.Empty(t, role)
}

func CreateRandomPermissionList(t *testing.T, size int) []int64 {
	var permissionsID []int64
	for i := 0; i < size; i++ {
		name := util.RandomPermission()
		permission, err := testStore.CreatePermission(context.Background(), name)
		require.NoError(t, err)
		permissionsID = append(permissionsID, permission.ID)
	}
	require.Len(t, permissionsID, size)
	return permissionsID
}
