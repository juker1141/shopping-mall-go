package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juker1141/shopping-mall-go/api"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	DB_SOURCE     = "postgresql://postgres:password@localhost:5432/shopping_mall?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db.", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
