package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/dmitrymomot/go-env"
	_ "github.com/lib/pq" // init pg driver
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/types"
	"github.com/zeebo/errs"
)

var (
	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 10)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// Wallets
	feePayerAccountPrivateKey = env.MustString("FEE_PAYER_ACCOUNT")
	assetAccountPrivateKey    = env.MustString("ASSET_ACCOUNT")
	issuerAccountPrivateKey   = env.MustString("ISSUER_ACCOUNT")
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
		log.Fatalf("walletRepo error: %v", err)
	}

	feePayer := types.AccountFromPrivateKeyBytes([]byte(feePayerAccountPrivateKey))
	asset := types.AccountFromPrivateKeyBytes([]byte(assetAccountPrivateKey))
	issuer := types.AccountFromPrivateKeyBytes([]byte(issuerAccountPrivateKey))
	walletService := wallet.NewService(repo, *solana.New(client.DevnetRPCEndpoint, feePayer, asset, issuer))

	if err := walletService.Bootstrap(ctx); err != nil {
		log.Fatalf("walletService error: %v", err)
	}
}
