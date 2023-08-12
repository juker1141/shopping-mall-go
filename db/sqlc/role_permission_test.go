package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomRolePermission(t *testing.T) RolePermission {
	permission := createRandomPermission(t)
	role := createRandomRole(t)

	rolePermission, err := testStore.CreateRolePermission(context.Background(), CreateRolePermissionParams{
		PermissionID: pgtype.Int4{
			Int32: int32(permission.ID),
			Valid: true,
		},
		RoleID: pgtype.Int4{
			Int32: int32(role.ID),
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, rolePermission)

	require.NotZero(t, rolePermission.PermissionID)
	require.NotZero(t, rolePermission.RoleID)

	return rolePermission
}

func TestCreateRolePermission(t *testing.T) {
	createRandomRolePermission(t)
}

func TestGetRolePermission(t *testing.T) {
	rolePermission1 := createRandomRolePermission(t)

	arg := GetRolePermissionParams{
		RoleID: pgtype.Int4{
			Int32: rolePermission1.RoleID.Int32,
			Valid: true,
		},
		PermissionID: pgtype.Int4{
			Int32: rolePermission1.PermissionID.Int32,
			Valid: true,
		},
	}

	rolePermission2, err := testStore.GetRolePermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rolePermission2)
}

func TestDeleteRolePermissionByPermissionId(t *testing.T) {
	rolePermission1 := createRandomRolePermission(t)

	err := testStore.DeleteRolePermissionByPermissionId(context.Background(), rolePermission1.PermissionID)
	require.NoError(t, err)

	rolePermission2, err := testStore.GetRolePermission(context.Background(), GetRolePermissionParams{
		RoleID: pgtype.Int4{
			Int32: rolePermission1.RoleID.Int32,
			Valid: true,
		},
		PermissionID: pgtype.Int4{
			Int32: rolePermission1.PermissionID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, rolePermission2)
}

func TestDeleteRolePermissionByRoleId(t *testing.T) {
	rolePermission1 := createRandomRolePermission(t)

	err := testStore.DeleteRolePermissionByRoleId(context.Background(), rolePermission1.RoleID)
	require.NoError(t, err)

	rolePermission2, err := testStore.GetRolePermission(context.Background(), GetRolePermissionParams{
		RoleID: pgtype.Int4{
			Int32: rolePermission1.RoleID.Int32,
			Valid: true,
		},
		PermissionID: pgtype.Int4{
			Int32: rolePermission1.PermissionID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, rolePermission2)
}

func TestListRolePermissionByPermissionId(t *testing.T) {
	permission := createRandomPermission(t)
	for i := 0; i < 5; i++ {
		role := createRandomRole(t)
		testStore.CreateRolePermission(context.Background(), CreateRolePermissionParams{
			PermissionID: pgtype.Int4{
				Int32: int32(permission.ID),
				Valid: true,
			},
			RoleID: pgtype.Int4{
				Int32: int32(role.ID),
				Valid: true,
			},
		})
	}

	rolePermissions, err := testStore.ListRolePermissionByPermissionId(context.Background(), pgtype.Int4{
		Int32: int32(permission.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, rolePermissions, 5)

	for _, rolePermission := range rolePermissions {
		require.NotEmpty(t, rolePermission)
	}
}

func TestListRolePermissionByRoleId(t *testing.T) {
	role := createRandomRole(t)
	for i := 0; i < 5; i++ {
		permission := createRandomPermission(t)
		testStore.CreateRolePermission(context.Background(), CreateRolePermissionParams{
			PermissionID: pgtype.Int4{
				Int32: int32(permission.ID),
				Valid: true,
			},
			RoleID: pgtype.Int4{
				Int32: int32(role.ID),
				Valid: true,
			},
		})
	}

	rolePermissions, err := testStore.ListRolePermissionByRoleId(context.Background(), pgtype.Int4{
		Int32: int32(role.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, rolePermissions, 5)

	for _, rolePermission := range rolePermissions {
		require.NotEmpty(t, rolePermission)
	}
}
