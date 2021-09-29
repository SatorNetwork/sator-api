package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	dbx "github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/solana"
	userRepository "github.com/SatorNetwork/sator-api/svc/auth/repository"
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
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 10)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// Solana
	solanaApiBaseUrl = env.MustString("SOLANA_API_BASE_URL")
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

	userRepo, err :=  userRepository.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("userRepo error: %v", err)
	}

	usersNumber, err := userRepo.CountAllUsers(ctx)
	if err != nil {
		log.Fatalf("userRepo error: %v", err)
	}

	log.Printf("USERS NUMBER: %d", usersNumber)

	txFn := dbx.Transaction(db)

	for i := 0; int64(i) < usersNumber; i++ {
		// time.Sleep(time.Second*5)

		log.Println("Getting user")
		ul, err := userRepo.GetUsersListDesc(ctx, userRepository.GetUsersListDescParams{
			Limit: 1,
			Offset: int32(i),
		})
		if err != nil {
			if dbx.IsNotFoundError(err) {
				log.Print("No more users")
				break
			}
			continue
		}

		if len(ul) < 1 {
			log.Print("No more users")
			break
		}

		user := ul[0]

		if err := txFn(func (tx dbx.DBTX) error {
			return createSolanaWalletIfNotExists(ctx, repository.New(tx), solana.New(solanaApiBaseUrl), user.ID)
		}); err != nil {
			log.Printf("Create user wallet if not exists: %v", err)
		}
	}
}


func createSolanaWalletIfNotExists(ctx context.Context, repo *repository.Queries, sc *solana.Client, userID uuid.UUID) error {
	log.Println("Getting user SAO wallet")
	userWallet, err := repo.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID: userID,
		WalletType: wallet.WalletTypeSator,
	})
	if err != nil && !dbx.IsNotFoundError(err){
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

	feePayer, err := repo.GetSolanaAccountByType(ctx, wallet.FeePayerAccount.String())
	if err != nil {
		return fmt.Errorf("could not get fee payer account: %w", err)
	}
	issuer, err := repo.GetSolanaAccountByType(ctx, wallet.IssuerAccount.String())
	if err != nil {
		return fmt.Errorf("could not get issuer account: %w", err)
	}
	asset, err := repo.GetSolanaAccountByType(ctx, wallet.AssetAccount.String())
	if err != nil {
		return fmt.Errorf("could not get asset account: %w", err)
	}

	log.Println("Creating user SAO wallet")
	acc := sc.NewAccount()

	txHash, err := sc.InitAccountToUseAsset(
		ctx,
		sc.AccountFromPrivatekey(feePayer.PrivateKey),
		sc.AccountFromPrivatekey(issuer.PrivateKey),
		sc.AccountFromPrivatekey(asset.PrivateKey),
		acc,
	)
	if err != nil {
		return fmt.Errorf("could not init token holder account: %w", err)
	}
	log.Printf("init token holder account transaction: %s", txHash)

	sacc, err := repo.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: wallet.TokenAccount.String(),
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
		UserID: userID,
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

	return nil
}