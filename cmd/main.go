package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const appName = "rest-go-pg"

func main() {
	config, err := parseConfig(appName)
	if err != nil {
		log.Fatalf("parsing config; %v", err)
	}

	connectStr := fmt.Sprintf(config.DB.Tmpl, config.DB.Host, config.DB.Port, config.DB.Name, config.DB.User, config.DB.Password, appName)

	db, err := createDB(connectStr, config.DB.ConnLifetime, config.DB.MaxIdleConns, config.DB.PoolSize)
	if err != nil {
		log.Errorf("opening connection: %v", err)
	}
	log.RegisterExitHandler(func() {
		db.Close()
	})

	handler, err := createHTTPHandler(db)
	if err != nil {
		log.Fatalf("creating http handler: %v", err)
	}

	listenErr := make(chan error, 1)
	server := &http.Server{
		Addr:    config.API.Port,
		Handler: handler,
	}

	go func() {
		log.Println("REST-GO-PG APP STARTED", time.Now().Format(time.RFC3339))
		listenErr <- server.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		log.Fatal(err)
	case <-osSignals:
		server.SetKeepAlivesEnabled(false)
		timeout := time.Second * 5

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("STOP APP")
		log.Exit(0)
	}
}

func createDB(connectStr string, connLife time.Duration, maxIdle, poolSize int) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatalf("opening DB connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("establishing DB connection; %v", err)
	}

	db.SetConnMaxLifetime(connLife)
	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(poolSize)

	return
}
