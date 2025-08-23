package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.chaitanya.observer/internal/models"
)

type application struct {
	logger        *slog.Logger
	cfg           *Config
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

type Config struct {
	address   string
	staticDir string
	dsn       string
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.address, "addr", ":3000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		cfg:           &cfg,
		snippets:      &models.SnippetModel{DB: db}, // instantiate the SnippetModel struct with the positional arguments
		templateCache: templateCache,
	}

	logger.Info("Starting server", slog.String("addr", cfg.address))

	err = http.ListenAndServe(cfg.address, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
