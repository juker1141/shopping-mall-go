// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: admin_user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAdminUser = `-- name: CreateAdminUser :one
INSERT INTO admin_users (
  account,
  full_name,
  hashed_password,
  status,
  role_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, account, full_name, hashed_password, role_id, status, password_changed_at, created_at
`

type CreateAdminUserParams struct {
	Account        string      `json:"account"`
	FullName       string      `json:"full_name"`
	HashedPassword string      `json:"hashed_password"`
	Status         int32       `json:"status"`
	RoleID         pgtype.Int4 `json:"role_id"`
}

func (q *Queries) CreateAdminUser(ctx context.Context, arg CreateAdminUserParams) (AdminUser, error) {
	row := q.db.QueryRow(ctx, createAdminUser,
		arg.Account,
		arg.FullName,
		arg.HashedPassword,
		arg.Status,
		arg.RoleID,
	)
	var i AdminUser
	err := row.Scan(
		&i.ID,
		&i.Account,
		&i.FullName,
		&i.HashedPassword,
		&i.RoleID,
		&i.Status,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAdminUser = `-- name: DeleteAdminUser :exec
DELETE FROM admin_users WHERE id = $1
`

func (q *Queries) DeleteAdminUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAdminUser, id)
	return err
}

const getAdminUser = `-- name: GetAdminUser :one
SELECT id, account, full_name, hashed_password, role_id, status, password_changed_at, created_at FROM admin_users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAdminUser(ctx context.Context, id int64) (AdminUser, error) {
	row := q.db.QueryRow(ctx, getAdminUser, id)
	var i AdminUser
	err := row.Scan(
		&i.ID,
		&i.Account,
		&i.FullName,
		&i.HashedPassword,
		&i.RoleID,
		&i.Status,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getAdminUserByAccount = `-- name: GetAdminUserByAccount :one
SELECT id, account, full_name, hashed_password, role_id, status, password_changed_at, created_at FROM admin_users
WHERE account = $1 LIMIT 1
`

func (q *Queries) GetAdminUserByAccount(ctx context.Context, account string) (AdminUser, error) {
	row := q.db.QueryRow(ctx, getAdminUserByAccount, account)
	var i AdminUser
	err := row.Scan(
		&i.ID,
		&i.Account,
		&i.FullName,
		&i.HashedPassword,
		&i.RoleID,
		&i.Status,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getAdminUsersCount = `-- name: GetAdminUsersCount :one
SELECT COUNT(*) FROM admin_users
`

func (q *Queries) GetAdminUsersCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getAdminUsersCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listAdminUsers = `-- name: ListAdminUsers :many
SELECT id, account, full_name, hashed_password, role_id, status, password_changed_at, created_at FROM admin_users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAdminUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAdminUsers(ctx context.Context, arg ListAdminUsersParams) ([]AdminUser, error) {
	rows, err := q.db.Query(ctx, listAdminUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AdminUser{}
	for rows.Next() {
		var i AdminUser
		if err := rows.Scan(
			&i.ID,
			&i.Account,
			&i.FullName,
			&i.HashedPassword,
			&i.RoleID,
			&i.Status,
			&i.PasswordChangedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPermissionsForAdminUser = `-- name: ListPermissionsForAdminUser :many
SELECT DISTINCT p.id, p.name, p.created_at
FROM admin_users AS au
JOIN roles AS r ON au.role_id = r.id
JOIN role_permissions AS rp ON r.id = rp.role_id
JOIN permissions AS p ON rp.permission_id = p.id
WHERE au.id = $1
`

func (q *Queries) ListPermissionsForAdminUser(ctx context.Context, id int64) ([]Permission, error) {
	rows, err := q.db.Query(ctx, listPermissionsForAdminUser, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Permission{}
	for rows.Next() {
		var i Permission
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPermissionsIDByAccount = `-- name: ListPermissionsIDByAccount :many
SELECT p.id
FROM admin_users AS au
JOIN roles AS r ON au.role_id = r.id
JOIN role_permissions AS rp ON r.id = rp.role_id
JOIN permissions AS p ON rp.permission_id = p.id
WHERE au.account = $1
`

func (q *Queries) ListPermissionsIDByAccount(ctx context.Context, account string) ([]int64, error) {
	rows, err := q.db.Query(ctx, listPermissionsIDByAccount, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAdminUser = `-- name: UpdateAdminUser :one
UPDATE admin_users
SET 
  hashed_password = COALESCE($1, hashed_password),
  password_changed_at = COALESCE($2, password_changed_at),
  full_name = COALESCE($3, full_name),
  status = COALESCE($4, status),
  role_id = COALESCE($5, role_id)
WHERE
  id = $6
RETURNING id, account, full_name, hashed_password, role_id, status, password_changed_at, created_at
`

type UpdateAdminUserParams struct {
	HashedPassword    pgtype.Text        `json:"hashed_password"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	FullName          pgtype.Text        `json:"full_name"`
	Status            pgtype.Int4        `json:"status"`
	RoleID            pgtype.Int4        `json:"role_id"`
	ID                int64              `json:"id"`
}

func (q *Queries) UpdateAdminUser(ctx context.Context, arg UpdateAdminUserParams) (AdminUser, error) {
	row := q.db.QueryRow(ctx, updateAdminUser,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Status,
		arg.RoleID,
		arg.ID,
	)
	var i AdminUser
	err := row.Scan(
		&i.ID,
		&i.Account,
		&i.FullName,
		&i.HashedPassword,
		&i.RoleID,
		&i.Status,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
