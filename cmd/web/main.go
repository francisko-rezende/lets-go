package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

func main() {
	mux := http.NewServeMux()
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static/", "path to the directory used for serving static files")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	logger.Info("starting server", slog.String("addr", cfg.addr))
	error := http.ListenAndServe(cfg.addr, mux)
	logger.Error(error.Error())
	os.Exit(1)
}
