package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/0xdonnie7/buythrift_API/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	env  string
	port int
	db   struct {
		dsn          string
		maxIdleConns int
		maxOpenConns int
		maxIdleTime  string
	}
	jwt struct {
		secret string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	models data.Models
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "", "JWT secret key")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if cfg.db.dsn == "" {
		log.Fatal("db-dsn is required")
	}

	if cfg.jwt.secret == "" {
		log.Fatal("jwt secret is required")
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	log.Println("database connection pool established")

	app := &application{
		config: cfg,
		models: data.NewModels(db),
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
