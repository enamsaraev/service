package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// db connection
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()
}
