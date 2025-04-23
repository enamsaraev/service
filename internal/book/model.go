package book

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"service/pkg"
)

const (
	createNewBook = "INSERT INTO book (name, count) VALUES ($1, $2) RETURNING id;"

	getBookByName        = "SELECT * FROM book WHERE name=$1;"
	getBookByIdForUpdate = "SELECT * FROM book WHERE id=$1 FOR UPDATE;"

	updateBookById = "UPDATE book SET name = $1, count = $2 WHERE id = $3"
)

type Book struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func UpdateFields(oldBook Book, newBook *Book) {
	reflectOldBook := reflect.ValueOf(oldBook)
	reflectNewBook := reflect.ValueOf(newBook).Elem()

	for i := 0; i < reflectNewBook.NumField(); i++ {
		oldField := reflectOldBook.Field(i)
		newField := reflectNewBook.Field(i)

		if newField.IsZero() {
			newField.Set(oldField)
		}
	}
}

type BookModel struct {
	conn   *pgxpool.Conn
	logger *pkg.Logger
}

func GetBookModel(conn *pgxpool.Conn, logger *pkg.Logger) *BookModel {
	return &BookModel{
		conn:   conn,
		logger: logger,
	}
}

func (bm BookModel) AddBook(ctx context.Context, book Book) (int, error) {
	rows, err := bm.conn.Query(ctx, createNewBook, book.Name, book.Count)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var bookId int

	if rows.Next() {
		err = rows.Scan(&bookId)
		if err != nil {
			return 0, err
		}
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}

	return bookId, nil
}

func (bm BookModel) GetBook(ctx context.Context, name string) (Book, error) {
	var book Book

	rows, err := bm.conn.Query(ctx, getBookByName, name)
	if err != nil {
		return book, err
	}

	defer rows.Close()

	if rows.Next() {
		// TODO: проверка количества строк, пока только UNIQUE на поле name
		err = rows.Scan(&book.Id, &book.Name, &book.Count)
		if err != nil {
			return book, err
		}
	} else {
		if err = rows.Err(); err != nil {
			return book, err
		}

		return book, errors.New("book not found")
	}
	return book, nil
}

func (bm BookModel) UpdateBook(ctx context.Context, book Book) (err error) {
	tx, err := bm.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}(tx, ctx)

	bookRow, err := tx.Query(ctx, getBookByIdForUpdate, book.Id)
	if err != nil {
		return err
	}

	if bookRow.Next() {
		var oldBook Book

		err = bookRow.Scan(&oldBook.Id, &oldBook.Name, &oldBook.Count)
		if err != nil {
			return err
		}

		bookRow.Close()
		UpdateFields(oldBook, &book)
	} else {
		return errors.New("book not found")
	}

	_, err = tx.Exec(ctx, updateBookById, book.Name, book.Count, book.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
