package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
	cfg    *Config
}

type Config struct {
	address   string
	staticDir string
}

func main() {

	cfg := Config{}
	flag.StringVar(&cfg.address, "addr", ":3000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
		cfg:    &cfg,
	}

	logger.Info("Starting server", slog.String("addr", cfg.address))

	err := http.ListenAndServe(cfg.address, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
