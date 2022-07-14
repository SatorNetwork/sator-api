package main

import (
	"database/sql"
	"encoding/base64"
	"time"

	log "github.com/sirupsen/logrus"

	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	authRepo "github.com/SatorNetwork/sator-api/svc/auth/repository"
	tx_watcher_svc "github.com/SatorNetwork/sator-api/svc/tx_watcher"
	tx_watcher_repository "github.com/SatorNetwork/sator-api/svc/tx_watcher/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	walletRepo "github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/dmitrymomot/go-env"
	"github.com/portto/solana-go-sdk/types"
	"github.com/zeebo/errs"
	"golang.org/x/net/context"
)

var (
	// DB
	dbConnString = env.MustString("DATABASE_URL")
	//dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 3)
	//dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 0)

	// Solana
	solanaApiBaseUrl            = env.MustString("SOLANA_API_BASE_URL")
	solanaAssetAddr             = env.MustString("SOLANA_ASSET_ADDR")
	solanaFeePayerAddr          = env.MustString("SOLANA_FEE_PAYER_ADDR")
	solanaFeePayerPrivateKey    = env.MustString("SOLANA_FEE_PAYER_PRIVATE_KEY")
	solanaTokenHolderAddr       = env.MustString("SOLANA_TOKEN_HOLDER_ADDR")
	solanaTokenHolderPrivateKey = env.MustString("SOLANA_TOKEN_HOLDER_PRIVATE_KEY")
	solanaSystemProgram         = env.MustString("SOLANA_SYSTEM_PROGRAM")
	solanaSysvarRent            = env.MustString("SOLANA_SYSVAR_RENT")
	solanaSysvarClock           = env.MustString("SOLANA_SYSVAR_CLOCK")
	solanaSplToken              = env.MustString("SOLANA_SPL_TOKEN")
	solanaStakeProgramID        = env.MustString("SOLANA_STAKE_PROGRAM_ID")

	minAmountToTransfer            = env.GetFloat("MIN_AMOUNT_TO_TRANSFER", 0)
	solanaStakePoolAddr            = env.MustString("SOLANA_STAKE_POOL_ADDR")
	fraudDetectionMode             = env.GetBool("FRAUD_DETECTION_MODE", false)
	tokenTransferPercent           = env.GetFloat("TOKEN_TRANSFER_PERCENT", 0.75)
	claimRewardsPercent            = env.GetFloat("CLAIM_REWARDS_PERCENT", 0.75)
	enableResourceIntensiveQueries = env.GetBool("ENABLE_RESOURCE_INTENSIVE_QUERIES", false)
)

func main() {
	solanaClient := solana_client.New(solanaApiBaseUrl, solana_client.Config{
		SystemProgram:  solanaSystemProgram,
		SysvarRent:     solanaSysvarRent,
		SysvarClock:    solanaSysvarClock,
		SplToken:       solanaSplToken,
		StakeProgramID: solanaStakeProgramID,
	}, nil)

	var feePayer types.Account
	{
		feePayerPk, err := base64.StdEncoding.DecodeString(solanaFeePayerPrivateKey)
		if err != nil {
			log.Fatalf("feePayerPk base64 decoding error: %v", err)
		}
		if err := solanaClient.CheckPrivateKey(solanaFeePayerAddr, feePayerPk); err != nil {
			log.Fatalf("solanaClient.CheckPrivateKey: fee payer: %v", err)
		}
		feePayer, err = types.AccountFromBytes(feePayerPk)
		if err != nil {
			log.Fatalf("can't get fee payer account from bytes")
		}
	}

	var tokenHolder types.Account
	{
		tokenHolderPk, err := base64.StdEncoding.DecodeString(solanaTokenHolderPrivateKey)
		if err != nil {
			log.Fatalf("tokenHolderPk base64 decoding error: %v", err)
		}
		if err := solanaClient.CheckPrivateKey(solanaTokenHolderAddr, tokenHolderPk); err != nil {
			log.Fatalf("solanaClient.CheckPrivateKey: token holder: %v", err)
		}
		tokenHolder, err = types.AccountFromBytes(tokenHolderPk)
		if err != nil {
			log.Fatalf("can't get token holder account from bytes")
		}
	}

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("init db connection error: %v", err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	walletRepository, err := walletRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("can't prepare wallet repository: %v", err)
	}

	txWatcherRepository, err := tx_watcher_repository.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("can't prepare tx watcher repository: %v", err)
	}

	txWatcherSvc := tx_watcher_svc.NewService(
		txWatcherRepository,
		solanaClient,
		feePayer,
		tokenHolder,
	)

	walletService := wallet.NewService(
		walletRepository,
		solanaClient,
		nil,
		txWatcherSvc,
		wallet.WithAssetSolanaAddress(solanaAssetAddr),
		wallet.WithSolanaFeePayer(solanaFeePayerAddr, feePayer.PrivateKey),
		wallet.WithSolanaTokenHolder(solanaTokenHolderAddr, tokenHolder.PrivateKey),
		wallet.WithMinAmountToTransfer(minAmountToTransfer),
		wallet.WithStakePoolSolanaAddress(solanaStakePoolAddr),
		wallet.WithFraudDetectionMode(fraudDetectionMode),
		wallet.WithTokenTransferPercent(tokenTransferPercent),
		wallet.WithClaimRewardsPercent(claimRewardsPercent),
		wallet.WithResourceIntensiveQueries(enableResourceIntensiveQueries),
	)

	// auth repo
	authRepository, err := authRepo.Prepare(ctx, db)
	if err != nil {
		log.Fatalf("authRepo error: %v", err)
	}

	stakes, err := walletRepository.GetAllStakes(ctx)
	if err != nil {
		log.Fatalf("can't get wallets: %v", err)
	}

	total := len(stakes)
	log.Infof("total wallets: %d", total)

	for i, s := range stakes {
		if i > 0 {
			time.Sleep(time.Second * 30)
		}

		user, err := authRepository.GetUserByID(ctx, s.UserID)
		if err != nil {
			log.Errorf("can't get user user_id=%s err: %v", s.UserID, err)
			continue
		}

		solAcc, err := walletRepository.GetSolanaAccountByUserIDAndType(ctx, walletRepo.GetSolanaAccountByUserIDAndTypeParams{
			UserID:     s.UserID,
			WalletType: wallet.WalletTypeSator,
		})
		if err != nil {
			log.Errorf("could not get solana account by user id=%s and type=%s", s.UserID, wallet.WalletTypeSator)
			continue
		}

		log.Infof("%d/%d: %s %.9f wallet: %s", i+1, total, user.Email, s.StakeAmount, solAcc.PublicKey)

		// if user.Disabled {
		// 	continue
		// }

		if err = walletService.Unstake(ctx, s.UserID, s.WalletID, true); err != nil {
			log.Errorf("can't unstake user_id = %s, wallet_id = %s, err: %v", s.UserID, s.ID, err)
			continue
		}

		log.Infof("Unstake for user_id = %s done. %v/%v", s.UserID, i+1, total)
	}
}
