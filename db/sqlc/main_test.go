package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juker1141/shopping-mall-go/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db.", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
	// 開始單元測試，通過 m.Run()
	// m.Run() 會回傳退出代碼，讓 os.Exit() 退出
}
