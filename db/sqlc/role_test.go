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

func createRandomRole(t *testing.T) Role {
	arg := CreateRoleParams{
		Name:   util.RandomRole(),
		Status: 1,
	}

	role, err := testQueries.CreateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role)

	require.Equal(t, arg.Name, role.Name)
	require.Equal(t, arg.Status, role.Status)
	require.NotZero(t, role.ID)
	require.NotZero(t, role.CreatedAt)

	return role
}

func TestCreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestGetRole(t *testing.T) {
	role1 := createRandomRole(t)
	role2, err := testQueries.GetRole(context.Background(), role1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, role1.Name, role2.Name)
	require.Equal(t, role1.Status, role2.Status)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
}

func TestUpdateRole(t *testing.T) {
	role1 := createRandomRole(t)

	newRole := util.RandomRole()
	newStatus := int32(0)
	arg := UpdateRoleParams{
		ID: role1.ID,
		Name: pgtype.Text{
			String: newRole,
			Valid:  true,
		},
		Status: pgtype.Int4{
			Int32: newStatus,
			Valid: true,
		},
	}

	role2, err := testQueries.UpdateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, newRole, role2.Name)
	require.Equal(t, newStatus, role2.Status)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
}

func TestDeleteRole(t *testing.T) {
	role1 := createRandomRole(t)
	err := testQueries.DeleteRole(context.Background(), role1.ID)
	require.NoError(t, err)

	role2, err := testQueries.GetRole(context.Background(), role1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
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

	roles, err := testQueries.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, 5)

	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}
