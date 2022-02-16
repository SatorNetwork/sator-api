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

	dbx "github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/dmitrymomot/go-env"
	"github.com/google/uuid"
	"github.com/oklog/run"
	"github.com/zeebo/errs"

	_ "github.com/lib/pq" // init pg driver
)

// go build -o ./bin/missedwallets ./cmd/solana/missedwallets/

var (
	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 3)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// Solana
	solanaApiBaseUrl     = env.MustString("SOLANA_API_BASE_URL")
	solanaSystemProgram  = env.MustString("SOLANA_SYSTEM_PROGRAM")
	solanaSysvarRent     = env.MustString("SOLANA_SYSVAR_RENT")
	solanaSysvarClock    = env.MustString("SOLANA_SYSVAR_CLOCK")
	solanaSplToken       = env.MustString("SOLANA_SPL_TOKEN")
	solanaStakeProgramID = env.MustString("SOLANA_STAKE_PROGRAM_ID")

	interval = env.GetDuration("EXEC_INTERVAL", time.Hour)
)

func main() {
	log.Println("Start worker to create missed wallets")

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

	mwr, err := Prepare(ctx, db)
	if err != nil {
		log.Fatalf("missed wallet repo error: %v", err)
	}

	txFn := dbx.Transaction(db)

	// runtime group
	var g run.Group

	stop := false
	g.Add(func() error {
		for {
			usersWithoutWallets, err := mwr.GetUsersWithoutWallet(ctx)
			if err != nil {
				log.Fatalf("[ERROR] userRepo error: %v", err)
			}

			log.Printf("NUMBER OF USERS WITHOUT WALLETS: %d", len(usersWithoutWallets))

			for _, user := range usersWithoutWallets {
				if stop {
					return nil
				}

				if !user.VerifiedAt.Valid {
					continue
				}
				if err := txFn(func(tx dbx.DBTX) error {
					return createSolanaWalletIfNotExists(
						ctx,
						repository.New(tx),
						solana.New(solanaApiBaseUrl, solana.Config{
							SystemProgram:  solanaSystemProgram,
							SysvarRent:     solanaSysvarRent,
							SysvarClock:    solanaSysvarClock,
							SplToken:       solanaSplToken,
							StakeProgramID: solanaStakeProgramID,
						}),
						user.ID,
					)
				}); err != nil {
					log.Printf("[ERROR] Create user wallet if not exists: %v", err)
				}
			}

			time.Sleep(interval)
		}
	}, func(err error) {
		stop = true
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

func createSolanaWalletIfNotExists(ctx context.Context, repo *repository.Queries, sc *solana.Client, userID uuid.UUID) error {
	// log.Println("Getting user SAO wallet")
	userWallet, err := repo.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: wallet.WalletTypeSator,
	})
	if err != nil && !dbx.IsNotFoundError(err) {
		return nil
	}

	if userWallet.SolanaAccountID != uuid.Nil && userWallet.WalletType == wallet.WalletTypeSator {
		return nil
	}

	if userWallet.SolanaAccountID == uuid.Nil && userWallet.WalletType == wallet.WalletTypeSator {
		// log.Println("Deleting user SAO wallet without solana SPL token account")
		if err := repo.DeleteWalletByID(ctx, userWallet.ID); err != nil {
			log.Printf("Could not delete wallet with id=%s: %v", userWallet.ID.String(), err)
		}
	}

	// log.Println("Creating user SAO wallet")
	acc := sc.NewAccount()

	sacc, err := repo.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: wallet.GeneralAccount.String(),
		PublicKey:   acc.PublicKey.ToBase58(),
		PrivateKey:  acc.PrivateKey,
	})
	if err != nil {
		return fmt.Errorf("could not store solana account: %w", err)
	}

	if _, err := repo.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:          userID,
		SolanaAccountID: sacc.ID,
		WalletType:      wallet.WalletTypeSator,
		Sort:            1,
	}); err != nil {
		return fmt.Errorf("could not create new SAO wallet for user with id=%s: %w", userID.String(), err)
	}

	if _, err := repo.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: wallet.WalletTypeRewards,
	}); err != nil && dbx.IsNotFoundError(err) {
		// log.Println("Creating user rewards wallet")
		if _, err := repo.CreateWallet(ctx, repository.CreateWalletParams{
			UserID:     userID,
			WalletType: wallet.WalletTypeRewards,
			Sort:       2,
		}); err != nil {
			return fmt.Errorf("could not new rewards wallet for user with id=%s: %w", userID.String(), err)
		}
	}

	log.Printf("wallets has been created for user=%s", userID.String())

	return nil
}
