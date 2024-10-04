package main

import (
	"context"
	"database/sql"
	"github.com/ankodd/todo-server/internal/handler"
	"github.com/ankodd/todo-server/internal/metrics"
	"github.com/ankodd/todo-server/internal/middleware"
	"github.com/ankodd/todo-server/internal/storage/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get address for service
	serviceAddr := ":8080"
	if len(os.Args) > 2 {
		serviceAddr = ":" + os.Args[1]
	}

	// Get address for metrics
	metricsAddr := ":8082"
	if len(os.Args) == 3 {
		metricsAddr = ":" + os.Args[2]
	}

	// Initial storage
	const StoragePath = "storage/storage.db"
	db, err := sql.Open("sqlite3", StoragePath)
	if err != nil {
		log.Fatalf("Open: %v", err)
	}
	defer db.Close()

	// Test connection to storage
	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging db: %v", err)
	}

	store, err := sqlite.New(db)
	if err != nil {
		log.Fatalf("error initializing store: %v", err)
	}
	log.Printf("store initialized in %v", StoragePath)

	// Initial server
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         serviceAddr,
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
		Metrics:     metrics.NewMetrics(),
	}

	mux.Handle("/create", middleware.All(Handler.Create(ctx)))
	mux.Handle("/list", middleware.All(Handler.FetchAll(ctx)))
	mux.Handle("/update/", middleware.All(Handler.Update(ctx)))
	mux.Handle("/delete/", middleware.All(Handler.Delete(ctx)))
	mux.Handle("/count", middleware.All(Handler.CountEntries(ctx)))

	// Start metrics
	go func() {
		log.Printf("Starting metrics on %s\n", metricsAddr)

		if err := metrics.Listen(metricsAddr); err != nil {
			log.Fatalf("error starting metrics: %v", err)
		}
	}()

	// Start server
	log.Printf("Server listening on %s\n", serviceAddr)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
