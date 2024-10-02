package main

import (
	"context"
	"database/sql"
	"github.com/ankodd/todo-server/internal/handler"
	"github.com/ankodd/todo-server/internal/middleware"
	"github.com/ankodd/todo-server/internal/storage"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get port
	PORT := ":8080"
	if len(os.Args) != 1 {
		PORT = ":" + os.Args[1]
	}

	// Initial storage
	const StoragePath = "storage/storage.db"
	db, err := sql.Open("sqlite3", StoragePath)
	if err != nil {
		log.Fatalf("Open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ping: %v", err)
	}

	store, err := storage.New(db)
	if err != nil {
		log.Fatalf("storage.New: %v", err)
	}
	log.Printf("store initialized in %v", StoragePath)

	// Initial server
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	log.Printf("Server Config: %+v\n", s)

	// Initial handlers
	ctx := context.Background()
	Handler := handler.Handler{
		Storage:     store,
		IdleTimeout: s.IdleTimeout,
	}

	mux.Handle("/create", middleware.Middleware(Handler.Create(ctx)))
	mux.Handle("/list", middleware.Middleware(Handler.FetchAll(ctx)))
	mux.Handle("/update/", middleware.Middleware(Handler.Update(ctx)))
	mux.Handle("/delete/", middleware.Middleware(Handler.Delete(ctx)))
	mux.Handle("/count", middleware.Middleware(Handler.CountEntries(ctx)))

	// Start server
	log.Printf("Server listening on %s\n", PORT)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
