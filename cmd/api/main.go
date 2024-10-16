package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Mensurui/expenseAPI/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	l      *log.Logger
	config config
	db     *sql.DB
	models data.Model
}

func main() {
	logger := log.New(os.Stdout, "Logger", log.LstdFlags)
	var cfg config
	logger.Println("Starting server")
	flag.IntVar(&cfg.port, "port", 9090, "API port number")
	flag.StringVar(&cfg.db.dsn, "database", os.Getenv("EXPENSE_DB_DSN"), "Database Address")

	db, err := openDB(cfg)
	if err != nil {
		log.Println("Error opening the database")
		return
	}
	logger.Println("Connected successfully to the database")

	app := application{
		l:      logger,
		config: cfg,
		db:     db,
		models: data.New(db),
	}

	mux := app.routes()

	srv := http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           mux,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}

	err = srv.ListenAndServe()

	if err != nil {
		logger.Println("error serving")
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Println("Error: ", err)
		return nil, err
	}

	return db, nil
}
