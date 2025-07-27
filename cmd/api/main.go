package main

import (
	"flag"
	"sync"

	loggerjson "Github.com/Devaraja-Anu/voteblocks/internal/loggerJson"
)

const version = "0.50.0"

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenComms int
		maxIdleComms int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	cfg    config
	logger *loggerjson.Logger
	wg     sync.WaitGroup
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "Port Number", 4000, "API server port")

	// rate limiting
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable Rate Limiter")

	app := &application{
		cfg: cfg,
	}

	app.server()
}
