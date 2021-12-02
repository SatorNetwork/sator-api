package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // init pg driver

	dbx "github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/dmitrymomot/go-env"
	"github.com/oklog/run"
	"github.com/zeebo/errs"
)

var (
	dbConnString          = env.MustString("DATABASE_URL")
	dbMaxOpenConns        = env.GetInt("DATABASE_MAX_OPEN_CONNS", 3)
	dbMaxIdleConns        = env.GetInt("DATABASE_IDLE_CONNS", 0)
	bid                   = env.MustString("BACKOFFICE_DEVICE_ID")
	interval              = env.GetDuration("EXEC_INTERVAL", time.Minute)
	sanitizerInterval     = env.GetDuration("EXEC_SANITIZER_INTERVAL", time.Minute*15)
	rewardsInterval       = env.GetDuration("REWARDS_INTERVAL", time.Minute)
	rewardsEarnPeriod     = env.GetString("REWARDS_EARN_PERIOD", "00:05")
	rewardsWithdrawPeriod = env.GetString("REWARDS_WITHDRAW_PERIOD", "01:00")
	rewardsMinOperations  = env.GetInt("REWARDS_MIN_OPERATIONS", 10)
	rewardsCheckForPeriod = env.GetDuration("REWARDS_CHECK_FOR_PERIOD", time.Hour*24*1)
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
			log.Println("Start block with duplicate emails and users on the same device")
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

		g.Add(func() error {
			log.Println("Start user emails sanitizer")

			ticker := time.NewTicker(sanitizerInterval)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return nil
				case <-ticker.C:
					sanitizeUserEmails(ctx, repo)
				}
			}
		}, func(err error) {
			done <- true
		})
	}

	// Rewards scam
	{
		done := make(chan bool)
		defer close(done)

		g.Add(func() error {
			log.Println("Start rewards scam determining")

			if err := blockUsersWithFrequentTransactions(
				ctx, db, repo, rewardsEarnPeriod,
				rewardsWithdrawPeriod, rewardsMinOperations,
				rewardsCheckForPeriod,
			); err != nil {
				log.Printf("[ERROR] blockUsersWithFrequentTransactions: %v", err)
			}

			ticker := time.NewTicker(rewardsInterval)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return nil
				case <-ticker.C:
					if err := blockUsersWithFrequentTransactions(
						ctx, db, repo, rewardsEarnPeriod,
						rewardsWithdrawPeriod, rewardsMinOperations,
						rewardsCheckForPeriod,
					); err != nil {
						log.Printf("[ERROR] blockUsersWithFrequentTransactions: %v", err)
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

func sanitizeUserEmails(ctx context.Context, repo *repository.Queries) {
	log.Println("Start: sanitizeUserEmails")
	defer func() {
		log.Println("Finish: sanitizeUserEmails")
	}()

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
					log.Printf("could not clean email: %s: %v", user.Email, err)

					if err := repo.UpdateUserStatus(ctx, repository.UpdateUserStatusParams{
						ID:          user.ID,
						Disabled:    true,
						BlockReason: sql.NullString{String: "invalid email address", Valid: true},
					}); err != nil {
						log.Printf("could not block user with id=%s: %v", user.ID, err)
					}

					if !errors.Is(err, utils.ErrInvalidIcanSuffix) {
						sanitizedEmail = "n/a"
					}
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

func determineScamTransactions(ctx context.Context, db *sql.DB, uid string, trType int, period string) (bool, error) {
	query := `
SELECT
	COUNT(*) > 0 AS scam
FROM (
	SELECT
		user_id,
		(created_at - lag(created_at, 1) OVER (ORDER BY created_at)) AS diff
	FROM
		rewards
	WHERE
		user_id = $1
		AND transaction_type = $2 
    ORDER BY created_at DESC
) AS rewards_tx
WHERE
	diff < $3;
	`

	row := db.QueryRowContext(ctx, query, uid, trType, period)
	var result bool
	err := row.Scan(&result)

	return result, err
}

func getRewardedUserIDs(ctx context.Context, db *sql.DB, limit, offset int32, minOperations int, forPeriod time.Duration) ([]string, error) {
	query := `
SELECT
	user_id
FROM (
	SELECT
		user_id
	FROM
		rewards
	WHERE
		user_id NOT IN(
			SELECT
				id FROM users
			WHERE
				disabled = TRUE)
		AND created_at >= $4
	ORDER BY
		created_at DESC) AS rew_usr
GROUP BY
	user_id
HAVING
	count(*) > $3
LIMIT $1
OFFSET $2;
	`
	rows, err := db.QueryContext(ctx, query, limit, offset,
		minOperations, time.Now().Add(-forPeriod).Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []string
	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			return nil, err
		}
		items = append(items, row)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func blockUsersWithFrequentTransactions(ctx context.Context, db *sql.DB, repo *repository.Queries, earnPeriod, withdrawPeriod string, minOperations int, checkForPeriod time.Duration) error {
	log.Println("Start: blockUsersWithFrequentTransactions")
	defer func() {
		log.Println("Finish: blockUsersWithFrequentTransactions")
	}()

	var (
		limit  int32 = 10
		offset int32 = 0
	)
	for {
		userIDs, err := getRewardedUserIDs(ctx, db, limit, offset, minOperations, checkForPeriod)
		if err != nil {
			return fmt.Errorf("getRewardedUserIDs: %w", err)
		}

		offset += limit

		for _, id := range userIDs {
			// log.Printf("userID: %s", id)
			// Earning
			{
				isScam, err := determineScamTransactions(ctx, db, id, 1, earnPeriod)
				if err != nil {
					continue
				}
				if isScam {
					userID := uuid.MustParse(id)
					if err := repo.UpdateUserStatus(ctx, repository.UpdateUserStatusParams{
						ID:          userID,
						Disabled:    true,
						BlockReason: sql.NullString{String: "frequent rewards earning", Valid: true},
					}); err != nil {
						log.Printf("could not block user with id=%s: %v", userID, err)
					}
					log.Printf("blocked user with id=%s by reason: frequent rewards earning", id)
					continue
				}
			}

			// Withdrawing
			{
				isScam, err := determineScamTransactions(ctx, db, id, 2, withdrawPeriod)
				if err != nil {
					continue
				}
				if isScam {
					userID := uuid.MustParse(id)
					if err := repo.UpdateUserStatus(ctx, repository.UpdateUserStatusParams{
						ID:          userID,
						Disabled:    true,
						BlockReason: sql.NullString{String: "frequent rewards withdrawn", Valid: true},
					}); err != nil {
						log.Printf("could not block user with id=%s: %v", userID, err)
					}
					log.Printf("blocked user with id=%s by reason: frequent rewards withdrawn", id)
					continue
				}
			}
		}

	}
}
