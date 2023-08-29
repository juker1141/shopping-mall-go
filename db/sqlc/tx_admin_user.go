package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/val"
)

type CreateAdminUserTxParams struct {
	Account        string  `json:"account"`
	FullName       string  `json:"full_name"`
	HashedPassword string  `json:"hashed_password"`
	Status         int32   `json:"status"`
	RolesID        []int64 `json:"roles_id"`
}

type AdminUserTxResult struct {
	AdminUser AdminUser `json:"admin_user"`
	RoleList  []Role    `json:"role_list"`
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

		result.RoleList, err = q.ListRolesForAdminUser(ctx, result.AdminUser.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

type UpdateAdminUserTxParams struct {
	ID                int64     `json:"id"`
	FullName          string    `json:"full_name"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	Status            *int32    `json:"status"`
	RolesID           []int64   `json:"roles_id"`
}

func (store *SQLStore) UpdateAdminUserTx(ctx context.Context, arg UpdateAdminUserTxParams) (AdminUserTxResult, error) {
	var result AdminUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		updateAdminUserArg := UpdateAdminUserParams{
			ID: arg.ID,
		}

		if arg.FullName != "" {
			updateAdminUserArg.FullName = pgtype.Text{
				String: arg.FullName,
				Valid:  true,
			}
		}

		if arg.Status != nil && val.IsValidStatus(int(*arg.Status)) {
			updateAdminUserArg.Status = pgtype.Int4{
				Int32: *arg.Status,
				Valid: true,
			}
		}

		if arg.HashedPassword != "" {
			updateAdminUserArg.HashedPassword = pgtype.Text{
				String: arg.HashedPassword,
				Valid:  true,
			}
			updateAdminUserArg.PasswordChangedAt = pgtype.Timestamptz{
				Time:  arg.PasswordChangedAt,
				Valid: true,
			}
		}

		result.AdminUser, err = q.UpdateAdminUser(context.Background(), updateAdminUserArg)
		if err != nil {
			return err
		}

		// 如果有要變更權限才去執行
		if len(arg.RolesID) > 0 || arg.RolesID != nil {
			// 刪除角色舊的權限關聯
			err = q.DeleteAdminUserRoleByAdminUserId(ctx, pgtype.Int4{
				Int32: int32(result.AdminUser.ID),
				Valid: true,
			})
			if err != nil {
				return err
			}

			// 新增新的權限關聯
			for _, roleID := range arg.RolesID {
				_, err = q.CreateAdminUserRole(ctx, CreateAdminUserRoleParams{
					AdminUserID: pgtype.Int4{
						Int32: int32(result.AdminUser.ID),
						Valid: true,
					},
					RoleID: pgtype.Int4{
						Int32: int32(roleID),
						Valid: true,
					},
				})
				if err != nil {
					return err
				}
			}
		}

		roleList, err := q.ListRolesForAdminUser(ctx, result.AdminUser.ID)
		if err != nil {
			return err
		}

		result.RoleList = roleList

		return nil
	})

	return result, err
}

type DeleteAdminUserTxParams struct {
	ID int64 `json:"id"`
}

type DeleteAdminUserTxResult struct {
	Message string `json:"message"`
}

func (store *SQLStore) DeleteAdminUserTx(ctx context.Context, arg DeleteAdminUserTxParams) (DeleteRoleTxResult, error) {
	var result DeleteRoleTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// 在這裡執行刪除角色的操作
		err := q.DeleteAdminUserRoleByAdminUserId(ctx, pgtype.Int4{
			Int32: int32(arg.ID),
			Valid: true,
		})
		if err != nil {
			return err
		}

		err = q.DeleteAdminUser(ctx, arg.ID)
		if err != nil {
			return err
		}

		result.Message = "Delete adminUser success."

		return nil
	})

	return result, err
}
