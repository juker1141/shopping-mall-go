package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// RoleTxParams contains the input parameters of the role create
type CreateRoleTxParams struct {
	Name          string  `json:"name"`
	Status        int32   `json:"status"`
	PermissionsID []int64 `json:"permissions_id"`
}

type CreateRoleTxResult struct {
	Role           Role         `json:"role"`
	PermissionList []Permission `json:"permission_list"`
}

// It creates a role, rolePermission, and get all permissions name within a single database trasaction
func (store *SQLStore) CreateRoleTx(ctx context.Context, arg CreateRoleTxParams) (CreateRoleTxResult, error) {
	var result CreateRoleTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var permissionList []Permission

		if len(arg.PermissionsID) <= 0 {
			err = fmt.Errorf("at least one permission is required")
			return err
		}

		result.Role, err = q.CreateRole(ctx, CreateRoleParams{
			Name:   arg.Name,
			Status: arg.Status,
		})
		if err != nil {
			return err
		}

		for _, permissionId := range arg.PermissionsID {
			_, err := q.CreateRolePermission(ctx, CreateRolePermissionParams{
				RoleID: pgtype.Int4{
					Int32: int32(result.Role.ID),
					Valid: true,
				},
				PermissionID: pgtype.Int4{
					Int32: int32(permissionId),
					Valid: true,
				},
			})
			if err != nil {
				return err
			}

			permission, err := q.GetPermission(ctx, permissionId)
			if err != nil {
				return err
			}
			permissionList = append(permissionList, permission)
		}

		result.PermissionList = permissionList

		return nil
	})

	return result, err
}

type UpdateRoleTxParams struct {
	ID            int64   `json:"role_id"`
	Name          string  `json:"name"`
	Status        int32   `json:"status"`
	PermissionsID []int64 `json:"permissions_id"`
}

type UpdateRoleTxResult struct {
	Role           Role         `json:"role"`
	PermissionList []Permission `json:"permission_list"`
}

func (store *SQLStore) UpdateRoleTx(ctx context.Context, arg UpdateRoleTxParams) (UpdateRoleTxResult, error) {
	var result UpdateRoleTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var permissionList []Permission

		result.Role, err = q.UpdateRole(context.Background(), UpdateRoleParams{
			ID: arg.ID,
			Name: pgtype.Text{
				String: arg.Name,
				Valid:  true,
			},
			Status: pgtype.Int4{
				Int32: arg.Status,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		// 刪除角色舊的權限關聯
		err = q.DeleteRolePermissionByRoleId(ctx, pgtype.Int4{
			Int32: int32(result.Role.ID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		// 新增新的權限關聯
		for _, permissionID := range arg.PermissionsID {
			_, err = q.CreateRolePermission(ctx, CreateRolePermissionParams{
				RoleID: pgtype.Int4{
					Int32: int32(result.Role.ID),
					Valid: true,
				},
				PermissionID: pgtype.Int4{
					Int32: int32(permissionID),
					Valid: true,
				},
			})
			if err != nil {
				return err
			}

			permission, err := q.GetPermission(ctx, permissionID)
			if err != nil {
				return err
			}
			permissionList = append(permissionList, permission)
		}

		result.PermissionList = permissionList

		return nil
	})

	return result, err
}

type DeleteRoleTxParams struct {
	ID int64 `json:"role_id"`
}

type DeleteRoleTxResult struct {
	Message string `json:"message"`
}

func (store *SQLStore) DeleteRoleTx(ctx context.Context, arg DeleteRoleTxParams) (DeleteRoleTxResult, error) {
	var result DeleteRoleTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// 在這裡執行刪除角色的操作
		err := q.DeleteRolePermissionByRoleId(ctx, pgtype.Int4{
			Int32: int32(arg.ID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		err = q.DeleteRole(ctx, arg.ID)
		if err != nil {
			return err
		}

		result.Message = "Delete role success."

		return nil
	})

	return result, err
}
