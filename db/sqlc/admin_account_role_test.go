package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomAdminUserRole(t *testing.T) AdminUserRole {
	role := createRandomRole(t)
	adminUser := createRandomAdminUser(t)

	arg := CreateAdminUserRoleParams{
		RoleID: pgtype.Int4{
			Int32: int32(role.ID),
			Valid: true,
		},
		AdminUserID: pgtype.Int4{
			Int32: int32(adminUser.ID),
			Valid: true,
		},
	}

	adminUserRole, err := testStore.CreateAdminUserRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, adminUserRole)

	require.NotZero(t, adminUserRole.RoleID)
	require.NotZero(t, adminUserRole.AdminUserID)

	return adminUserRole
}

func TestCreateAdminUserRole(t *testing.T) {
	createRandomAdminUserRole(t)
}

func TestGetAdminUserRole(t *testing.T) {
	adminUserRole1 := createRandomAdminUserRole(t)

	arg := GetAdminUserRoleParams{
		RoleID: pgtype.Int4{
			Int32: adminUserRole1.RoleID.Int32,
			Valid: true,
		},
		AdminUserID: pgtype.Int4{
			Int32: adminUserRole1.AdminUserID.Int32,
			Valid: true,
		},
	}

	adminUserRole2, err := testStore.GetAdminUserRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, adminUserRole2)
}

func TestDeleteAdminUserRoleByRoleId(t *testing.T) {
	adminUserRole1 := createRandomAdminUserRole(t)

	err := testStore.DeleteAdminUserRoleByRoleId(context.Background(), adminUserRole1.RoleID)
	require.NoError(t, err)

	adminUserRole2, err := testStore.GetAdminUserRole(context.Background(), GetAdminUserRoleParams{
		RoleID: pgtype.Int4{
			Int32: adminUserRole1.RoleID.Int32,
			Valid: true,
		},
		AdminUserID: pgtype.Int4{
			Int32: adminUserRole1.AdminUserID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, adminUserRole2)
}

func TestDeleteAdminUserRoleByAdminUserId(t *testing.T) {
	adminUserRole1 := createRandomAdminUserRole(t)

	err := testStore.DeleteAdminUserRoleByAdminUserId(context.Background(), adminUserRole1.AdminUserID)
	require.NoError(t, err)

	adminUserRole2, err := testStore.GetAdminUserRole(context.Background(), GetAdminUserRoleParams{
		RoleID: pgtype.Int4{
			Int32: adminUserRole1.RoleID.Int32,
			Valid: true,
		},
		AdminUserID: pgtype.Int4{
			Int32: adminUserRole1.AdminUserID.Int32,
			Valid: true,
		},
	})
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, adminUserRole2)
}

func TestListAdminUserRoleByRoleId(t *testing.T) {
	role := createRandomRole(t)
	for i := 0; i < 5; i++ {
		adminUser := createRandomAdminUser(t)
		testStore.CreateAdminUserRole(context.Background(), CreateAdminUserRoleParams{
			RoleID: pgtype.Int4{
				Int32: int32(role.ID),
				Valid: true,
			},
			AdminUserID: pgtype.Int4{
				Int32: int32(adminUser.ID),
				Valid: true,
			},
		})
	}

	adminUserRoles, err := testStore.ListAdminUserRoleByRoleId(context.Background(), pgtype.Int4{
		Int32: int32(role.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, adminUserRoles, 5)

	for _, adminUserRole := range adminUserRoles {
		require.NotEmpty(t, adminUserRole)
	}
}

func TestListAdminUserRoleByAdminUserId(t *testing.T) {
	adminUser := createRandomAdminUser(t)
	for i := 0; i < 5; i++ {
		role := createRandomRole(t)
		testStore.CreateAdminUserRole(context.Background(), CreateAdminUserRoleParams{
			RoleID: pgtype.Int4{
				Int32: int32(role.ID),
				Valid: true,
			},
			AdminUserID: pgtype.Int4{
				Int32: int32(adminUser.ID),
				Valid: true,
			},
		})
	}

	adminUserRoles, err := testStore.ListAdminUserRoleByAdminUserId(context.Background(), pgtype.Int4{
		Int32: int32(adminUser.ID),
		Valid: true,
	})
	require.NoError(t, err)
	require.Len(t, adminUserRoles, 5)

	for _, adminUserRole := range adminUserRoles {
		require.NotEmpty(t, adminUserRole)
	}
}
