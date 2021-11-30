package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // init pg driver

	"github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/dmitrymomot/go-env"
	"github.com/oklog/run"
	"github.com/zeebo/errs"
)

var (
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 3)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)
	interval       = env.GetDuration("EXEC_INTERVAL", time.Minute)
	bid            = env.MustString("BACKOFFICE_DEVICE_ID")
)

func main() {
	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("init db connection error: %v", err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	if err := db.Ping(); err != nil {
		log.Fatalf("db pinng error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := repository.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("user repo error: %v", err)
	}

	// runtime group
	var g run.Group

	done := make(chan bool)
	defer close(done)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	g.Add(func() error {
		for {
			select {
			case <-done:
				return nil
			case <-ticker.C:
				if err := repo.BlockUsersOnTheSameDevice(ctx, bid); err != nil {
					log.Printf("[ERROR] BlockUsersOnTheSameDevice: %v", err)
				}
			}
		}
	}, func(err error) {
		done <- true
	})

	g.Add(func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		c := <-sigChan
		return fmt.Errorf("terminated with sig %q", c)
	}, func(err error) {})

	if err := g.Run(); err != nil {
		log.Println("terminated with error:", err)
	}
}
