package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAdminUserTxParams struct {
	Account        string  `json:"account"`
	FullName       string  `json:"full_name"`
	HashedPassword string  `json:"hashed_password"`
	Status         int32   `json:"status"`
	RolesID        []int64 `json:"roles_id"`
}

type AdminUserTxResult struct {
	AdminUser      AdminUser    `json:"admin_user"`
	PermissionList []Permission `json:"permission_list"`
}

func (store *SQLStore) CreateAdminUserTx(ctx context.Context, arg CreateAdminUserTxParams) (AdminUserTxResult, error) {
	var result AdminUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		if len(arg.RolesID) <= 0 {
			err = fmt.Errorf("at least one role is required")
			return err
		}

		adArg := CreateAdminUserParams{
			Account:        arg.Account,
			FullName:       arg.FullName,
			HashedPassword: arg.HashedPassword,
			Status:         arg.Status,
		}

		result.AdminUser, err = q.CreateAdminUser(ctx, adArg)
		if err != nil {
			return err
		}

		for _, roleId := range arg.RolesID {
			_, err := q.CreateAdminUserRole(ctx, CreateAdminUserRoleParams{
				RoleID: pgtype.Int4{
					Int32: int32(roleId),
					Valid: true,
				},
				AdminUserID: pgtype.Int4{
					Int32: int32(result.AdminUser.ID),
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
		}

		result.PermissionList, err = q.ListPermissionForAdminUser(ctx, result.AdminUser.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
