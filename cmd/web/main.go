package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr string
}

type application struct {
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))
	error := http.ListenAndServe(cfg.addr, app.routes())
	logger.Error(error.Error())
	os.Exit(1)
}
