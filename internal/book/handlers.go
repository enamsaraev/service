package book

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func AddBookHandler(bm *BookModel) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody Book

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			bm.logger.Errorf("Error while decoding request data: %v", err)
			http.Error(w, "Broken data", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		createdBookId, err := bm.AddBook(ctx, requestBody)
		if err != nil {
			bm.logger.Errorf("Error while adding book: %v", err)
			http.Error(w, "Error while executing a query", http.StatusBadRequest)
			return
		}

		response := map[string]int{
			"created": createdBookId,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})
}

func GetBookHandler(bm *BookModel) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookName, ok := mux.Vars(r)["name"]
		if !ok {
			http.Error(w, "Invalid book name", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		foundBook, err := bm.GetBook(ctx, bookName)
		if err != nil {
			bm.logger.Errorf("Error while finding book: %v", err)
			http.Error(w, "Error while executing a query", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(foundBook)
	})
}

func UpdateBookHandler(bm *BookModel) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody Book

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			bm.logger.Errorf("Error while decoding request data: %v", err)
			http.Error(w, "Broken data", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := bm.UpdateBook(ctx, requestBody); err != nil {
			bm.logger.Errorf("Error while updateting book: %v", err)
			http.Error(w, "Error while executing a query", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
