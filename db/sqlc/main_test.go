package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	DB_SOURCE = "postgresql://postgres:password@localhost:5432/shopping_mall?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db.", err)
	}

	testQueries = New(connPool)

	os.Exit(m.Run())
	// 開始單元測試，通過 m.Run()
	// m.Run() 會回傳退出代碼，讓 os.Exit() 退出
}
