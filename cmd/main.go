package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	customServer "service/server"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Start(sigs chan os.Signal) {
	<-sigs
}

func main() {
	// db connection
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// server
	server := customServer.CreateNewServer("localhost", "8080")

	apiRouter := server.Router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	server.Run()

	// starting
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	Start(sigs)
}
