// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAdminUser(ctx context.Context, arg CreateAdminUserParams) (AdminUser, error)
	CreateAdminUserRole(ctx context.Context, arg CreateAdminUserRoleParams) (AdminUserRole, error)
	CreatePermission(ctx context.Context, name string) (Permission, error)
	CreateRole(ctx context.Context, name string) (Role, error)
	CreateRolePermission(ctx context.Context, arg CreateRolePermissionParams) (RolePermission, error)
	DeleteAdminUser(ctx context.Context, id int64) error
	DeleteAdminUserRoleByAdminUserId(ctx context.Context, adminUserID pgtype.Int4) error
	DeleteAdminUserRoleByRoleId(ctx context.Context, roleID pgtype.Int4) error
	DeletePermission(ctx context.Context, id int64) error
	DeleteRole(ctx context.Context, id int64) error
	DeleteRolePermissionByPermissionId(ctx context.Context, permissionID pgtype.Int4) error
	DeleteRolePermissionByRoleId(ctx context.Context, roleID pgtype.Int4) error
	GetAdminUser(ctx context.Context, id int64) (AdminUser, error)
	GetAdminUserByAccount(ctx context.Context, account string) (AdminUser, error)
	GetAdminUserRole(ctx context.Context, arg GetAdminUserRoleParams) (AdminUserRole, error)
	GetAdminUsersCount(ctx context.Context) (int64, error)
	GetPermission(ctx context.Context, id int64) (Permission, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	GetRolePermission(ctx context.Context, arg GetRolePermissionParams) (RolePermission, error)
	GetRolesCount(ctx context.Context) (int64, error)
	ListAdminUserRoleByAdminUserId(ctx context.Context, adminUserID pgtype.Int4) ([]AdminUserRole, error)
	ListAdminUserRoleByRoleId(ctx context.Context, roleID pgtype.Int4) ([]AdminUserRole, error)
	ListAdminUsers(ctx context.Context, arg ListAdminUsersParams) ([]AdminUser, error)
	ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]Permission, error)
	ListPermissionsForAdminUser(ctx context.Context, id int64) ([]Permission, error)
	ListPermissionsForRole(ctx context.Context, id int64) ([]Permission, error)
	ListRolePermissionByPermissionId(ctx context.Context, permissionID pgtype.Int4) ([]RolePermission, error)
	ListRolePermissionByRoleId(ctx context.Context, roleID pgtype.Int4) ([]RolePermission, error)
	ListRoles(ctx context.Context, arg ListRolesParams) ([]Role, error)
	ListRolesForAdminUser(ctx context.Context, id int64) ([]Role, error)
	UpdateAdminUser(ctx context.Context, arg UpdateAdminUserParams) (AdminUser, error)
	UpdatePermission(ctx context.Context, arg UpdatePermissionParams) (Permission, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
}

var _ Querier = (*Queries)(nil)
