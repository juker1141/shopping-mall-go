package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// RoleTxParams contains the input parameters of the role create
type RoleTxParams struct {
	Name          string  `json:"name"`
	PermissionsID []int64 `json:"permissions_id"`
}

type RoleTxResult struct {
	Role           Role         `json:"role"`
	PermissionList []Permission `json:"permission_list"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database trasaction
func (store *Store) RoleTx(ctx context.Context, arg RoleTxParams) (RoleTxResult, error) {
	var result RoleTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var permissionList []Permission

		if len(arg.PermissionsID) <= 0 {
			err = fmt.Errorf("at least one permission is required")
			return err
		}

		result.Role, err = q.CreateRole(ctx, arg.Name)
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
