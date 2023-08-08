package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomPermission(t *testing.T) Permission {
	name := util.RandomPermission()

	permission, err := testQueries.CreatePermission(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, permission)

	require.Equal(t, name, permission.Name)
	require.NotZero(t, permission.ID)
	require.NotZero(t, permission.CreatedAt)

	return permission
}

func TestCreatePermission(t *testing.T) {
	createRandomPermission(t)
}

func TestGetPermission(t *testing.T) {
	permission1 := createRandomPermission(t)
	permission2, err := testQueries.GetPermission(context.Background(), permission1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, permission2)

	require.Equal(t, permission1.ID, permission2.ID)
	require.Equal(t, permission1.Name, permission2.Name)

	require.WithinDuration(t, permission1.CreatedAt, permission2.CreatedAt, time.Second)
}

func TestUpdatePermission(t *testing.T) {
	permission1 := createRandomPermission(t)

	arg := UpdatePermissionParams{
		ID:   permission1.ID,
		Name: util.RandomPermission(),
	}

	permission2, err := testQueries.UpdatePermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, permission2)

	require.Equal(t, permission1.ID, permission2.ID)
	require.Equal(t, arg.Name, permission2.Name)
	require.WithinDuration(t, permission1.CreatedAt, permission2.CreatedAt, time.Second)

	require.NotEqual(t, permission1.Name, permission2.Name)
}

func TestDeletePermission(t *testing.T) {
	permission1 := createRandomPermission(t)
	err := testQueries.DeletePermission(context.Background(), permission1.ID)
	require.NoError(t, err)

	permission2, err := testQueries.GetPermission(context.Background(), permission1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, permission2)
}

func TestListPermission(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPermission(t)
	}

	arg := ListPermissionsParams{
		Limit:  5,
		Offset: 5,
	}

	permissions, err := testQueries.ListPermissions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, permissions, 5)

	for _, permission := range permissions {
		require.NotEmpty(t, permission)
	}
}
