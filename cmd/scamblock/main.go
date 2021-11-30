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

	dbx "github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/dmitrymomot/go-env"
	"github.com/oklog/run"
	"github.com/zeebo/errs"
)

var (
	dbConnString      = env.MustString("DATABASE_URL")
	dbMaxOpenConns    = env.GetInt("DATABASE_MAX_OPEN_CONNS", 3)
	dbMaxIdleConns    = env.GetInt("DATABASE_IDLE_CONNS", 0)
	interval          = env.GetDuration("EXEC_INTERVAL", time.Minute)
	sanitizerInterval = env.GetDuration("EXEC_SANITIZER_INTERVAL", time.Minute*15)
	bid               = env.MustString("BACKOFFICE_DEVICE_ID")
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

	// Blockers
	{
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
					if err := repo.BlockUsersWithDuplicateEmail(ctx); err != nil {
						log.Printf("[ERROR] BlockUsersWithDuplicateEmail: %v", err)
					}
				}
			}
		}, func(err error) {
			done <- true
		})
	}

	// Sanitizer
	{
		done := make(chan bool)
		defer close(done)

		ticker := time.NewTicker(sanitizerInterval)
		defer ticker.Stop()

		g.Add(func() error {
			for {
				select {
				case <-done:
					return nil
				case <-ticker.C:
					for {
						users, err := repo.GetNotSanitizedUsersListDesc(ctx, repository.GetNotSanitizedUsersListDescParams{
							Limit:  10,
							Offset: 0,
						})
						if err != nil && dbx.IsNotFoundError(err) {
							break
						}

						for _, user := range users {
							if !user.SanitizedEmail.Valid || len(user.SanitizedEmail.String) < 5 {
								// Sanitize email address
								sanitizedEmail, err := utils.SanitizeEmail(user.Email)
								if err != nil {
									log.Printf("could not cleam email: %s: %v", user.Email, err)
									continue
								}

								if err := repo.UpdateUserSanitizedEmail(ctx, repository.UpdateUserSanitizedEmailParams{
									ID:             user.ID,
									SanitizedEmail: sanitizedEmail,
								}); err != nil {
									log.Printf("could not add sanitizesd email for user with id=%s and email=%s: %v", user.ID, user.Email, err)
								}

								log.Printf("Sanitized user with id=%s and email=%s", user.ID, user.Email)
							}
						}
					}
				}
			}
		}, func(err error) {
			done <- true
		})
	}

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
