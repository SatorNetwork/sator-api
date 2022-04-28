package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/dmitrymomot/go-env"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // init pg driver
	"github.com/zeebo/errs"

	dbx "github.com/SatorNetwork/sator-api/lib/db"
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	userRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
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
	solanaSystemProgram         = env.MustString("SOLANA_SYSTEM_PROGRAM")
	solanaSysvarRent            = env.MustString("SOLANA_SYSVAR_RENT")
	solanaSysvarClock           = env.MustString("SOLANA_SYSVAR_CLOCK")
	solanaSplToken              = env.MustString("SOLANA_SPL_TOKEN")
	solanaStakeProgramID        = env.MustString("SOLANA_STAKE_PROGRAM_ID")

	userEmail = env.MustString("USER_EMAIL")
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

	ur, err := userRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("missed user repo error: %v", err)
	}

	txFn := dbx.Transaction(db)

	user, err := ur.GetUserByEmail(ctx, userEmail)
	if err != nil {
		log.Fatalf("user with email=%s is not found", userEmail)
	}
	if !user.VerifiedAt.Valid {
		log.Fatalf("user with email=%s is not verified yet", userEmail)
	}
	if yes, _ := mwr.IsEmailWhitelisted(ctx, userEmail); !yes {
		log.Fatalf("user with email=%s is blocked", userEmail)
	}

	if err := txFn(func(tx dbx.DBTX) error {
		return createSolanaWalletIfNotExists(
			ctx,
			repository.New(tx),
			solana_client.New(solanaApiBaseUrl, solana_client.Config{
				SystemProgram:  solanaSystemProgram,
				SysvarRent:     solanaSysvarRent,
				SysvarClock:    solanaSysvarClock,
				SplToken:       solanaSplToken,
				StakeProgramID: solanaStakeProgramID,
			}, nil),
			user.ID,
			feePayerPk,
			tokenHolderPk,
		)
	}); err != nil {
		log.Printf("Create user wallet if not exists: %v", err)
	}

	fmt.Printf("finished")
}

func createSolanaWalletIfNotExists(ctx context.Context, repo *repository.Queries, sc lib_solana.Interface, userID uuid.UUID, feePayerPk, tokenHolderPk []byte) error {
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

	feePayer, err := sc.AccountFromPrivateKeyBytes(feePayerPk)
	if err != nil {
		return err
	}
	txHash, err := sc.CreateAccountWithATA(
		ctx,
		solanaAssetAddr,
		acc.PublicKey.ToBase58(),
		feePayer,
	)
	if err != nil {
		return fmt.Errorf("could not init token holder account: %w", err)
	}
	log.Printf("init token holder account transaction: %s", txHash)

	return nil
}
