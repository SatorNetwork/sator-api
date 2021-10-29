package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/SatorNetwork/sator-api/internal/ethereum"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/dmitrymomot/go-env"
	"github.com/zeebo/errs"

	_ "github.com/lib/pq" // init pg driver
)

var (
	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 10)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// Solana
	solanaApiBaseUrl = env.MustString("SOLANA_API_BASE_URL")
	systemProgram    = env.MustString("SOLANA_SYSTEM_PROGRAM")
	sysvarRent       = env.MustString("SOLANA_SYSVAR_RENT")
	sysvarClock      = env.MustString("SOLANA_SYSVAR_CLOCK")
	splToken         = env.MustString("SOLANA_SPL_TOKEN")
	stakeProgramID   = env.MustString("SOLANA_STAKE_PROGRAM_ID")
	rewardProgramID  = env.MustString("SOLANA_REWARD_PROGRAM_ID")
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

	ethereumClient, err := ethereum.NewClient()
	if err != nil {
		log.Fatalf("failed to init eth client: %v", err)
	}

	repo, err := repository.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("walletRepo error: %v", err)
	}
	walletService := wallet.NewService(repo, solana.New(solanaApiBaseUrl, solana.Config{
		SystemProgram:   systemProgram,
		SysvarRent:      sysvarRent,
		SysvarClock:     sysvarClock,
		SplToken:        splToken,
		StakeProgramID:  stakeProgramID,
		RewardProgramID: rewardProgramID,
	}), ethereumClient)

	if err := walletService.Bootstrap(ctx); err != nil {
		log.Fatalf("walletService error: %v", err)
	}
}
