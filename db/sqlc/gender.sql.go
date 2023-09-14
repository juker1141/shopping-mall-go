// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: gender.sql

package db

import (
	"context"
)

const createGender = `-- name: CreateGender :one
INSERT INTO genders (
  name
) VALUES (
  $1
) RETURNING id, name
`

func (q *Queries) CreateGender(ctx context.Context, name string) (Gender, error) {
	row := q.db.QueryRow(ctx, createGender, name)
	var i Gender
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getGender = `-- name: GetGender :one
SELECT id, name FROM genders
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetGender(ctx context.Context, id int64) (Gender, error) {
	row := q.db.QueryRow(ctx, getGender, id)
	var i Gender
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listGenders = `-- name: ListGenders :many
SELECT id, name FROM genders
ORDER BY id
`

func (q *Queries) ListGenders(ctx context.Context) ([]Gender, error) {
	rows, err := q.db.Query(ctx, listGenders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Gender{}
	for rows.Next() {
		var i Gender
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}