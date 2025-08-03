package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"Github.com/Devaraja-Anu/voteblocks/internal/db"
	loggerjson "Github.com/Devaraja-Anu/voteblocks/internal/loggerJson"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "0.50.0"

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int32
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	cfg     config
	logger  *loggerjson.Logger
	wg      sync.WaitGroup
	queries *db.Queries
}

func main() {

	var cfg config

	var maxOpenConns int

	flag.IntVar(&cfg.port, "Port", 4000, "API server port")

	flag.StringVar(&cfg.db.dsn, "db-dsn",
		"postgres://greenlight:pa55word@db:5432/greenlight?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&maxOpenConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// rate limiting
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable Rate Limiter")

	cfg.db.maxOpenConns = int32(maxOpenConns)

	logger := loggerjson.New(os.Stdout, loggerjson.LevelInfo)

	conn, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
	}

	queries := db.New(conn)

	app := &application{
		cfg:     cfg,
		logger:  logger,
		queries: queries,
	}

	err = app.server()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	// Parse max idle time to Duration
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.db.maxOpenConns
	poolConfig.MaxConnIdleTime = duration

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
