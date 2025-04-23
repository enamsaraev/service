package main

import (
	"context"
	"os"
	"os/signal"
	"service/internal/book"
	"service/pkg"
	customServer "service/server"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Env struct {
	book *book.BookModel
}

func Start(sigs chan os.Signal) {
	<-sigs
}

func main() {
	godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// logger
	logger := pkg.GetLogger()

	// db connection
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Errorf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	pool, err := dbPool.Acquire(ctx)
	if err != nil {
		logger.Errorf("Error while acquiring connection from the database pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Release()

	// server
	server := customServer.CreateNewServer("localhost", "8080")

	handlers := Env{
		book: book.GetBookModel(pool, logger),
	}

	apiRouter := server.Router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/book/add", book.AddBookHandler(handlers.book)).Methods("POST")
	apiRouter.HandleFunc("/book/update", book.UpdateBookHandler(handlers.book)).Methods("POST")
	apiRouter.HandleFunc("/book/get/{name}", book.GetBookHandler(handlers.book)).Methods("GET")
	server.Run()

	// starting
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	Start(sigs)
}
