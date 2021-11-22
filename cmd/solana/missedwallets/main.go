package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	dbx "github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/dmitrymomot/go-env"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	_ "github.com/lib/pq" // init pg driver
)

var (
	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 3)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// Solana
	solanaApiBaseUrl = env.MustString("SOLANA_API_BASE_URL")
	solanaAssetAddr  = env.MustString("SOLANA_ASSET_ADDR")
	// solanaFeePayerAddr          = env.MustString("SOLANA_FEE_PAYER_ADDR")
	solanaFeePayerPrivateKey = env.MustString("SOLANA_FEE_PAYER_PRIVATE_KEY")
	// solanaTokenHolderAddr       = env.MustString("SOLANA_TOKEN_HOLDER_ADDR")
	solanaTokenHolderPrivateKey = env.MustString("SOLANA_TOKEN_HOLDER_PRIVATE_KEY")
)

func main() {
	feePayerPk, err := base64.StdEncoding.DecodeString(solanaFeePayerPrivateKey)
	if err != nil {
		log.Fatalf("feePayerPk base64 decoding error: %v", err)
	}
	tokenHolderPk, err := base64.StdEncoding.DecodeString(solanaTokenHolderPrivateKey)
	if err != nil {
		log.Fatalf("tokenHolderPk base64 decoding error: %v", err)
	}

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

	usersWithoutWallets, err := mwr.GetUsersWithoutWallet(ctx)
	if err != nil {
		log.Fatalf("userRepo error: %v", err)
	}

	log.Printf("NUMBER OF USERS WITHOUT WALLETS: %d", len(usersWithoutWallets))

	txFn := dbx.Transaction(db)

	counter := 0

	for _, user := range usersWithoutWallets {
		if !user.VerifiedAt.Valid {
			continue
		}
		if skip, _ := mwr.IsEmailBlacklisted(ctx, user.Email); skip {
			continue
		}

		counter++

		if err := txFn(func(tx dbx.DBTX) error {
			return createSolanaWalletIfNotExists(
				ctx,
				repository.New(tx),
				solana.New(solanaApiBaseUrl),
				user.ID,
				feePayerPk,
				tokenHolderPk,
			)
		}); err != nil {
			log.Printf("Create user wallet if not exists: %v", err)
		}
	}

	fmt.Printf("finished: %d", counter)
}

func createSolanaWalletIfNotExists(ctx context.Context, repo *repository.Queries, sc *solana.Client, userID uuid.UUID, feePayerPk, tokenHolderPk []byte) error {
	log.Println("Getting user SAO wallet")
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
		log.Println("Deleting user SAO wallet without solana SPL token account")
		if err := repo.DeleteWalletByID(ctx, userWallet.ID); err != nil {
			log.Printf("Could not delete wallet with id=%s: %v", userWallet.ID.String(), err)
		}
	}

	log.Println("Creating user SAO wallet")
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
		log.Println("Creating user rewards wallet")
		if _, err := repo.CreateWallet(ctx, repository.CreateWalletParams{
			UserID:     userID,
			WalletType: wallet.WalletTypeRewards,
			Sort:       2,
		}); err != nil {
			return fmt.Errorf("could not new rewards wallet for user with id=%s: %w", userID.String(), err)
		}
	}

	txHash, err := sc.CreateAccountWithATA(
		ctx,
		solanaAssetAddr,
		sc.AccountFromPrivateKeyBytes(feePayerPk),
		sc.AccountFromPrivateKeyBytes(tokenHolderPk),
		acc,
	)
	if err != nil {
		return fmt.Errorf("could not init token holder account: %w", err)
	}
	log.Printf("init token holder account transaction: %s", txHash)

	return nil
}
