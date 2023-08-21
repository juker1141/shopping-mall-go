package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomRole(t *testing.T) Role {
	name := util.RandomRole()

	role, err := testStore.CreateRole(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, name, role.Name)
	require.NotZero(t, role.ID)
	require.NotZero(t, role.CreatedAt)

	return role
}

func TestCreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestGetRole(t *testing.T) {
	role1 := createRandomRole(t)
	role2, err := testStore.GetRole(context.Background(), role1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, role1.Name, role2.Name)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
}

func TestUpdateRole(t *testing.T) {
	role1 := createRandomRole(t)

	newRoleName := util.RandomRole()
	arg := UpdateRoleParams{
		ID: role1.ID,
		Name: pgtype.Text{
			String: newRoleName,
			Valid:  true,
		},
	}

	role2, err := testStore.UpdateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, newRoleName, role2.Name)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
}

func TestDeleteRole(t *testing.T) {
	role1 := createRandomRole(t)
	err := testStore.DeleteRole(context.Background(), role1.ID)
	require.NoError(t, err)

	role2, err := testStore.GetRole(context.Background(), role1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, role2)
}

func TestListRole(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRole(t)
	}

	arg := ListRolesParams{
		Limit:  5,
		Offset: 5,
	}

	roles, err := testStore.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, 5)

	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}

func TestGetRolesCount(t *testing.T) {
	createRandomRole(t)

	count, err := testStore.GetRolesCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)
}
