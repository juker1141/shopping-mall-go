package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateRoleTx(ctx context.Context, arg CreateRoleTxParams) (RoleTxResult, error)
	UpdateRoleTx(ctx context.Context, arg UpdateRoleTxParams) (RoleTxResult, error)
	DeleteRoleTx(ctx context.Context, arg DeleteRoleTxParams) (DeleteRoleTxResult, error)
	CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (OrderTxResult, error)
	UpdateOrderTx(ctx context.Context, arg UpdateOrderTxParams) (OrderTxResult, error)
	DeleteOrderTx(ctx context.Context, arg DeleteOrderTxParams) error
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (UserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
	UpdateCartTx(ctx context.Context, arg UpdateCartTxParams) (CartTxResult, error)
	DeleteCartTx(ctx context.Context, arg DeleteCartTxParams) error
}

// Store provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new Store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
