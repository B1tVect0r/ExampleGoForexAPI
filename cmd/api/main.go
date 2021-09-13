package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/b1tvect0r/exchangerates/pkg/db"
	"github.com/b1tvect0r/exchangerates/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	dbConnString := os.Getenv("PSQL_CONNECTION_STRING")
	pool, err := pgxpool.Connect(ctx, dbConnString)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to establish connection pool to database at %s", dbConnString))
	}
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database at %s", dbConnString)
	}

	q := db.New(pool)

	s, err := server.New(q, server.WithAESAPIKeyProvider(os.Getenv("AES_KEY")))
	if err != nil {
		log.Fatalf("failed to create server: %s", err.Error())
	}

	if err = s.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
