package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

// config struct
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

// app status struct
type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

// main function
func main() {
	var cfg config

	password := os.Getenv("DB_PASSWORD")
	connection_string := fmt.Sprintf("postgres://arturofiliovilla:%spassword@localhost/movies?sslmode=disable", password)
	// setting cfg.port
	flag.IntVar(&cfg.port, "port", 8080, "Server running on port...")
	// setting cfg.env
	flag.StringVar(&cfg.env, "env", "development", "Application Environment (dev | prod)")
	flag.StringVar(&cfg.db.dsn, "dsn", connection_string, "Postgress connection string")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "kjlasdflkj1432lkadf089asdfljk23408asdfljk32408", "secret")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// databes connection
	db, err := openDb(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// app config
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}
	// routes
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// start listening to clients on server
	logger.Println("starting server on port", cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDb(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
