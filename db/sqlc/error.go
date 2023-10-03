package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	CheckDateViolation  = "23514"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrForeignKeyViolation = &pgconn.PgError{
	Code: ForeignKeyViolation,
}

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

var ErrCheckDateFailed = &pgconn.PgError{
	Code: CheckDateViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
