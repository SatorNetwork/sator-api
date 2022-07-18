package wallet

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
	log "github.com/sirupsen/logrus"

	"github.com/SatorNetwork/sator-api/lib/db"
	lib_errors "github.com/SatorNetwork/sator-api/lib/errors"
	"github.com/SatorNetwork/sator-api/lib/ethereum"
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	tx_watcher_alias "github.com/SatorNetwork/sator-api/svc/tx_watcher/alias"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

type (
	// Service struct
	Service struct {
		wr walletRepository
		sc solanaClient
		ec ethereumClient
		// rw rewardsService

		txWatcher txWatcher

		satorAssetName  string
		solanaAssetName string

		satorAssetSolanaAddr        string
		feePayerSolanaAddr          string
		feePayerSolanaPrivateKey    []byte
		stakePoolSolanaPublicKey    string
		tokenHolderSolanaAddr       string
		tokenHolderSolanaPrivateKey []byte

		walletDetailsURL        string // url template to get SOL & SAO wallet types details
		walletTransactionsURL   string // url template to get SOL & SAO wallet types transactions list
		rewardsWalletDetailsURL string // url template to get rewards wallet type details
		rewardsTransactionsURL  string // url template to get rewards wallet type transactions list

		minAmountToTransfer float64 // minimum amount to transfer request

		fraudDetectionMode bool // fraud detection mode

		tokenTransferPercent float64
		claimRewardsPercent  float64

		enableResourceIntensiveQueries bool
	}

	// ServiceOption function
	// interface to extend service via options
	ServiceOption func(*Service)

	walletRepository interface {
		CreateWallet(ctx context.Context, arg repository.CreateWalletParams) (repository.Wallet, error)
		GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error)
		GetWalletBySolanaAccountID(ctx context.Context, solanaAccountID uuid.UUID) (repository.Wallet, error)
		GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error)
		GetWalletByUserIDAndType(ctx context.Context, arg repository.GetWalletByUserIDAndTypeParams) (repository.Wallet, error)

		AddSolanaAccount(ctx context.Context, arg repository.AddSolanaAccountParams) (repository.SolanaAccount, error)
		GetSolanaAccountByID(ctx context.Context, id uuid.UUID) (repository.SolanaAccount, error)
		GetSolanaAccountByType(ctx context.Context, accountType string) (repository.SolanaAccount, error)
		GetSolanaAccountTypeByPublicKey(ctx context.Context, publicKey string) (string, error)
		GetSolanaAccountByUserIDAndType(ctx context.Context, arg repository.GetSolanaAccountByUserIDAndTypeParams) (repository.SolanaAccount, error)

		AddEthereumAccount(ctx context.Context, arg repository.AddEthereumAccountParams) (repository.EthereumAccount, error)
		GetEthereumAccountByID(ctx context.Context, id uuid.UUID) (repository.EthereumAccount, error)
		GetEthereumAccountByUserIDAndType(ctx context.Context, arg repository.GetEthereumAccountByUserIDAndTypeParams) (repository.EthereumAccount, error)

		AddStake(ctx context.Context, arg repository.AddStakeParams) (repository.Stake, error)
		DeleteStakeByUserID(ctx context.Context, userID uuid.UUID) error
		GetStakeByUserID(ctx context.Context, userID uuid.UUID) (repository.Stake, error)
		GetTotalStake(ctx context.Context) (float64, error)
		UpdateStake(ctx context.Context, arg repository.UpdateStakeParams) error

		GetAllStakeLevels(ctx context.Context) ([]repository.StakeLevel, error)
		GetAllEnabledStakeLevels(ctx context.Context) ([]repository.StakeLevel, error)
		GetStakeLevelByAmount(ctx context.Context, amount float64) (repository.GetStakeLevelByAmountRow, error)
		GetMinimalStakeLevel(ctx context.Context) (repository.StakeLevel, error)

		AddTokenTransfer(ctx context.Context, arg repository.AddTokenTransferParams) (repository.TokenTransfer, error)
		UpdateTokenTransfer(ctx context.Context, arg repository.UpdateTokenTransferParams) error

		CheckRecipientAddress(ctx context.Context, arg repository.CheckRecipientAddressParams) (int64, error)
		DoesUserHaveFraudulentTransfers(ctx context.Context, userID uuid.UUID) (bool, error)
		DoesUserMakeTransferForLastMinute(ctx context.Context, userID uuid.UUID) (bool, error)
	}

	solanaClient interface {
		GetAccountBalanceSOL(ctx context.Context, accPubKey string) (float64, error)
		GetTokenAccountBalanceWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (float64, error)
		NewAccount() types.Account
		AccountFromPrivateKeyBytes(pk []byte) (types.Account, error)
		FeeAccumulatorAddress() string
		GiveAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, issuer types.Account, recipientAddr string, amount float64) (string, error)
		PrepareSendAssetsTx(
			ctx context.Context,
			assetAddr string,
			feePayer types.Account,
			source types.Account,
			recipientAddr string,
			amount float64,
			cfg *lib_solana.SendAssetsConfig,
		) (*lib_solana.PrepareTxResponse, error)
		SendAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, source types.Account, recipientAddr string, amount float64, cfg *lib_solana.SendAssetsConfig) (string, error)
		GetTransactionsWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) ([]lib_solana.ConfirmedTransactionResponse, error)

		InitializeStakePool(ctx context.Context, feePayer, issuer types.Account, asset common.PublicKey) (txHast string, stakePool types.Account, err error)
		Stake(ctx context.Context, feePayer, userWallet types.Account, pool, asset common.PublicKey, duration int64, amount float64) (string, error)
		Unstake(ctx context.Context, feePayer, userWallet types.Account, stakePool, asset common.PublicKey) (string, error)
		IsTransactionSuccessful(ctx context.Context, txhash string) (bool, error)
	}

	ethereumClient interface {
		CreateAccount() (ethereum.Wallet, error)
	}

	txWatcher interface {
		SendAndWatchTx(ctx context.Context, message types.Message, accountAliases []tx_watcher_alias.Alias) (string, error)
	}

	// PreparedTransaction ...
	PreparedTransaction struct {
		AssetName       string  `json:"asset_name,omitempty"`
		Amount          float64 `json:"amount,omitempty"`
		RecipientAddr   string  `json:"recipient_address,omitempty"`
		Fee             float64 `json:"fee,omitempty"`
		TransactionHash string  `json:"tx_hash,omitempty"`
		SenderWalletID  string  `json:"sender_wallet_id,omitempty"`
	}

	StakeLevel struct {
		ID             string  `json:"id"`
		MinAmount      float64 `json:"min_amount"`
		MinDaysToStake int     `json:"min_days_to_stake"`
		Title          string  `json:"title"`
		SubTitle       string  `json:"sub_title"`
		Rewards        string  `json:"rewards"`
		IsCurrent      bool    `json:"is_current"`
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(wr walletRepository, sc solanaClient, ec ethereumClient, txWatcher txWatcher, opt ...ServiceOption) *Service {
	s := &Service{
		wr: wr,
		sc: sc,
		ec: ec,
		// rw: rw,
		txWatcher: txWatcher,

		solanaAssetName: "SOL",
		satorAssetName:  "SAO",

		walletDetailsURL:        "wallets/%s",
		walletTransactionsURL:   "wallets/%s/transactions",
		rewardsWalletDetailsURL: "rewards/wallet/%s",
		rewardsTransactionsURL:  "rewards/wallet/%s/transactions",

		minAmountToTransfer: 0,
	}

	for _, o := range opt {
		o(s)
	}

	return s
}

type UserWallet struct {
	WalletAddress    string `json:"wallet_public_key"`
	WalletPrivateKey string `json:"wallet_private_key"`
	TokenSymbol      string `json:"token_symbol"`
	TokenMintAddress string `json:"token_mint_address"`
}

// GetSaoWalletByUserID returns a wallet private key in base58 format or error
func (s *Service) GetSaoWalletByUserID(ctx context.Context, userID uuid.UUID) (UserWallet, error) {
	wallet, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return UserWallet{}, fmt.Errorf("failed to get wallet by userID: %w", err)
	}

	solAcc, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return UserWallet{}, fmt.Errorf("failed to get solana account by ID: %w", err)
	}

	return UserWallet{
		WalletAddress:    solAcc.PublicKey,
		WalletPrivateKey: base58.Encode(solAcc.PrivateKey),
		TokenSymbol:      s.satorAssetName,
		TokenMintAddress: s.satorAssetSolanaAddr,
	}, nil
}

// GetWallets returns current user's wallets list with balance
func (s *Service) GetWallets(ctx context.Context, uid uuid.UUID) (Wallets, error) {
	wallets, err := s.wr.GetWalletsByUserID(ctx, uid)
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get wallets list: %w", err)
	}

	result := make(Wallets, 0, len(wallets))
	for _, w := range wallets {
		if w.WalletType == WalletTypeEthereum {
			// disable etherium wallet
			continue
		}

		wli := WalletsListItem{
			ID:    w.ID.String(),
			Type:  w.WalletType,
			Order: w.Sort,
		}

		switch w.WalletType {
		case WalletTypeSolana, WalletTypeSator, WalletTypeEthereum:
			wli.GetDetailsURL = fmt.Sprintf(s.walletDetailsURL, w.ID.String())
			wli.GetTransactionsURL = fmt.Sprintf(s.walletTransactionsURL, w.ID.String())
		case WalletTypeRewards:
			wli.GetDetailsURL = fmt.Sprintf(s.rewardsWalletDetailsURL, w.ID.String())
			wli.GetTransactionsURL = fmt.Sprintf(s.rewardsTransactionsURL, w.ID.String())
		}

		result = append(result, wli)
	}

	return result, nil
}

// GetWalletByID returns wallet details by wallet id
func (s *Service) GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (Wallet, error) {
	w, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return Wallet{}, fmt.Errorf("%w wallet", ErrNotFound)
		}
		return Wallet{}, fmt.Errorf("could not get wallet: %w", err)
	}

	if w.UserID != userID {
		return Wallet{}, fmt.Errorf("%w: you have no permissions to get this wallet", ErrForbidden)
	}

	// if w.EthereumAccountID.UUID != uuid.Nil {
	// 	ea, err := s.wr.GetEthereumAccountByID(ctx, w.EthereumAccountID.UUID)
	// 	if err != nil {
	// 		if db.IsNotFoundError(err) {
	// 			return Wallet{}, fmt.Errorf("%w ethereum account for this wallet", ErrNotFound)
	// 		}
	// 		return Wallet{}, fmt.Errorf("could not get ethereum account for this wallet: %w", err)
	// 	}
	// 	return Wallet{
	// 		ID:                     w.ID.String(),
	// 		Order:                  w.Sort,
	// 		EthereumAccountAddress: ea.Address,
	// 		Actions: []Action{
	// 			{
	// 				Type: ActionSendTokens.String(),
	// 				Name: ActionSendTokens.Name(),
	// 				URL:  "",
	// 			},
	// 			{
	// 				Type: ActionReceiveTokens.String(),
	// 				Name: ActionReceiveTokens.Name(),
	// 				URL:  "",
	// 			},
	// 		},
	// 		Balance: []Balance{
	// 			{
	// 				Currency: "SAOE",
	// 				Amount:   0,
	// 			},
	// 			// {
	// 			// 	Currency: "USD",
	// 			// 	Amount:   0,
	// 			// },
	// 		},
	// 	}, nil
	// }

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return Wallet{}, fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return Wallet{}, fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	var balance []Balance

	switch sa.AccountType {
	case GeneralAccount.String():
		if bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey); err == nil {
			balance = []Balance{
				{
					Currency: s.satorAssetName,
					Amount:   bal,
				},
				// {
				// 	Currency: "USD",
				// 	Amount:   bal * 0.04, // FIXME: setup currency rate
				// },
			}
		}
		// case GeneralAccount.String():
		// 	if bal, err := s.sc.GetAccountBalanceSOL(ctx, sa.PublicKey); err == nil {
		// 		balance = []Balance{
		// 			{
		// 				Currency: s.solanaAssetName,
		// 				Amount:   bal,
		// 			},
		// 			// {
		// 			// 	Currency: "USD",
		// 			// 	Amount:   bal * 70, // FIXME: setup currency rate
		// 			// },
		// 		}
		// 	}
	}

	return Wallet{
		ID:                   w.ID.String(),
		Type:                 w.WalletType,
		Order:                w.Sort,
		SolanaAccountAddress: sa.PublicKey,
		Actions: []Action{
			{
				Type: ActionSendTokens.String(),
				Name: ActionSendTokens.Name(),
				URL:  "",
			},
			{
				Type: ActionReceiveTokens.String(),
				Name: ActionReceiveTokens.Name(),
				URL:  "",
			},
			// {
			// 	Type: ActionStakeTokens.String(),
			// 	Name: ActionStakeTokens.Name(),
			// 	URL:  "",
			// },
		},
		Balance: balance,
	}, nil
}

// CreateWallet creates wallet for user with specified id.
func (s *Service) CreateWallet(ctx context.Context, userID uuid.UUID) error {
	acc := s.sc.NewAccount()
	sacc, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: GeneralAccount.String(),
		PublicKey:   acc.PublicKey.ToBase58(),
		PrivateKey:  acc.PrivateKey,
	})
	if err != nil {
		return fmt.Errorf("could not store solana account: %w", err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:          userID,
		SolanaAccountID: sacc.ID,
		WalletType:      WalletTypeSator,
		Sort:            1,
	}); err != nil {
		return fmt.Errorf("could not create new SAO wallet for user with id=%s: %w", userID.String(), err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:     userID,
		WalletType: WalletTypeRewards,
		Sort:       2,
	}); err != nil {
		return fmt.Errorf("could not new rewards wallet for user with id=%s: %w", userID.String(), err)
	}

	/** Disabled untill the next release **

	ethAccount, err := s.ec.CreateAccount()
	if err != nil {
		return fmt.Errorf("could not create new eth account for user with id=%s: %w", userID.String(), err)
	}

	ethPrivateBytes, err := json.Marshal(ethAccount.PrivateKey)
	if err != nil {
		return fmt.Errorf("could not marshal eth private key for user=%s: %w", userID.String(), err)
	}

	ethPublicBytes, err := json.Marshal(ethAccount.PublicKey)
	if err != nil {
		return fmt.Errorf("could not marshal eth private key for user=%s: %w", userID.String(), err)
	}

	eacc, err := s.wr.AddEthereumAccount(ctx, repository.AddEthereumAccountParams{
		PublicKey:  ethPublicBytes,
		PrivateKey: ethPrivateBytes,
		Address:    ethAccount.Address,
	})
	if err != nil {
		return fmt.Errorf("could not store ethereum account: %w", err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:            userID,
		WalletType:        WalletTypeEthereum,
		EthereumAccountID: eacc.ID,
	}); err != nil {
		return fmt.Errorf("could not create new ethereum wallet for user with id=%s: %w", userID.String(), err)
	}
	**/

	return nil
}

// WithdrawRewards convert rewards into sator tokens
func (s *Service) WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (txhash string, err error) {
	user, err := s.wr.GetSolanaAccountByUserIDAndType(ctx, repository.GetSolanaAccountByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return "", fmt.Errorf("could not get user token account: %w", err)
	}

	if thbalance, err := s.sc.GetTokenAccountBalanceWithAutoDerive(
		ctx,
		s.satorAssetSolanaAddr,
		s.tokenHolderSolanaAddr,
	); err != nil || thbalance < amount {
		return "", ErrTokenHolderBalance
	}

	feePayer, err := s.sc.AccountFromPrivateKeyBytes(s.feePayerSolanaPrivateKey)
	if err != nil {
		return "", err
	}
	tokenHolder, err := s.sc.AccountFromPrivateKeyBytes(s.tokenHolderSolanaPrivateKey)
	if err != nil {
		return "", err
	}

	prepareTxResp, err := s.sc.PrepareSendAssetsTx(
		ctx,
		s.satorAssetSolanaAddr,
		feePayer,
		tokenHolder,
		user.PublicKey,
		amount,
		&lib_solana.SendAssetsConfig{
			PercentToCharge:           s.claimRewardsPercent,
			ChargeSolanaFeeFromSender: true,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "can't prepare send assets tx")
	}

	// sends token
	for i := 0; i < 5; i++ {
		if txhash, err = s.txWatcher.SendAndWatchTx(
			ctx,
			prepareTxResp.Tx.Message,
			[]tx_watcher_alias.Alias{tx_watcher_alias.FeePayerAlias, tx_watcher_alias.TokenHolderAlias},
		); err != nil {
			if i < 4 {
				log.Println(err)
			} else {
				return "", fmt.Errorf("could not claim rewards: %w", err)
			}
			time.Sleep(time.Second * 10)
		} else {
			log.Printf("user %s: successful transaction: rewards withdraw: %s", userID.String(), txhash)
			break
		}
	}

	return txhash, nil
}

// GetListTransactionsByWalletID returns list of all transactions of specific wallet.
func (s *Service) GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ Transactions, err error) {
	if !s.enableResourceIntensiveQueries {
		return Transactions{}, nil
	}

	transactions, err := s.getListTransactionsByWalletID(ctx, userID, walletID)
	if err != nil {
		return Transactions{}, err
	}

	var skippedTx uint
	notEmptyTransactions := make(Transactions, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.Amount != 0 {
			notEmptyTransactions = append(notEmptyTransactions, transaction)
		} else {
			skippedTx++
		}
	}

	fmt.Printf("going to return %v transactions, skipped tx: %v\n", len(notEmptyTransactions), skippedTx)

	// pagination for solana trannsactions is disabled
	// transactionsPaginated := paginateTransactions(notEmptyTransactions, int(offset), int(limit))
	// return transactionsPaginated, nil
	return notEmptyTransactions, nil
}

// getListTransactionsByWalletID returns list of all transactions of specific wallet.
func (s *Service) getListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID) (Transactions, error) {
	wallet, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != userID {
		return nil, ErrForbidden
	}

	if wallet.SolanaAccountID == uuid.Nil {
		return Transactions{}, nil
	}

	solanaAcc, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.sc.GetTransactionsWithAutoDerive(ctx, s.satorAssetSolanaAddr, solanaAcc.PublicKey)
	if err != nil {
		return nil, err
	}

	txList := make(Transactions, 0, len(transactions))
	for _, tx := range transactions {
		txList = append(txList, castSolanaTxToTransaction(tx, walletID))
	}

	return txList, nil
}

// CreateTransfer crates transaction from one account to another.
func (s *Service) CreateTransfer(ctx context.Context, walletID uuid.UUID, recipientPK, asset string, amount float64) (tx PreparedTransferTransaction, err error) {
	w, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return PreparedTransferTransaction{}, fmt.Errorf("could not find wallet: %w", err)
	}

	if tooManyRequests, _ := s.wr.DoesUserMakeTransferForLastMinute(ctx, w.UserID); tooManyRequests {
		return PreparedTransferTransaction{}, ErrTooManyRequests
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return PreparedTransferTransaction{}, fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return PreparedTransferTransaction{}, fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey)
	if err != nil {
		return PreparedTransferTransaction{}, fmt.Errorf("could not get wallet balance")
	}

	if bal < s.minAmountToTransfer {
		return PreparedTransferTransaction{}, fmt.Errorf("%w: %.2f", ErrMinimalAmountToSend, s.minAmountToTransfer)
	}

	if bal < amount {
		return PreparedTransferTransaction{}, fmt.Errorf("balance is lower then requested amount: %.2f", bal)
	}

	var toEncode struct {
		Amount        float64
		Asset         string
		RecipientAddr string
	}

	toEncode.Asset = asset
	toEncode.Amount = amount
	toEncode.RecipientAddr = recipientPK

	// log.Printf("toEncode: %+v", toEncode)

	encodedData, err := json.Marshal(toEncode)
	if err != nil {
		return PreparedTransferTransaction{}, fmt.Errorf("could not marshal amount and recipient pk: %w", err)
	}

	// log.Printf("CreateTransfer: toEncode: encodedData: %s", string(encodedData))

	var feeInSAO float64
	{
		feePayer, err := s.sc.AccountFromPrivateKeyBytes(s.feePayerSolanaPrivateKey)
		if err != nil {
			return PreparedTransferTransaction{}, err
		}

		wallet, err := s.wr.GetWalletByID(ctx, walletID)
		if err != nil {
			return PreparedTransferTransaction{}, fmt.Errorf("could not get solana account: %w", err)
		}
		solanaAcc, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
		if err != nil {
			return PreparedTransferTransaction{}, fmt.Errorf("could not get solana account: %w", err)
		}
		source, err := s.sc.AccountFromPrivateKeyBytes(solanaAcc.PrivateKey)
		if err != nil {
			return PreparedTransferTransaction{}, err
		}

		resp, err := s.sc.PrepareSendAssetsTx(
			ctx,
			s.satorAssetSolanaAddr,
			feePayer,
			source,
			recipientPK,
			amount,
			&lib_solana.SendAssetsConfig{
				PercentToCharge:           s.tokenTransferPercent,
				ChargeSolanaFeeFromSender: true,
				AllowFallbackToDefaultFee: true,
				DefaultFee:                1,
			},
		)
		if err != nil {
			return PreparedTransferTransaction{}, err
		}

		feeInSAO = resp.FeeInSAO
	}

	return PreparedTransferTransaction{
		AssetName:       asset,
		Amount:          amount,
		RecipientAddr:   recipientPK,
		Fee:             feeInSAO,
		TransactionHash: base58.Encode(encodedData),
		SenderWalletID:  walletID.String(),
	}, nil
}

func (s *Service) ConfirmTransfer(ctx context.Context, uid, walletID uuid.UUID, encodedData string) error {
	decoded, err := base58.Decode(encodedData)
	if err != nil {
		return fmt.Errorf("could not decode from base58: %w", err)
	}

	var toDecode struct {
		Amount        float64
		RecipientAddr string
	}
	if err := json.Unmarshal(decoded, &toDecode); err != nil {
		return fmt.Errorf("could not unmarshal: %w", err)
	}

	solAcc, err := s.wr.GetSolanaAccountByUserIDAndType(ctx, repository.GetSolanaAccountByUserIDAndTypeParams{
		UserID:     uid,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return fmt.Errorf("could not get user token account: %w", err)
	}

	tr, err := s.wr.AddTokenTransfer(ctx, repository.AddTokenTransferParams{
		UserID:           uid,
		SenderAddress:    solAcc.PublicKey,
		RecipientAddress: toDecode.RecipientAddr,
		Amount:           toDecode.Amount,
		Status:           TokenTransferStatusPending,
	})
	if err != nil {
		log.Printf("could not add token transfer: %v", err)
	}

	if s.fraudDetectionMode {
		fraudDetected, _ := s.wr.DoesUserHaveFraudulentTransfers(ctx, uid)
		i, _ := s.wr.CheckRecipientAddress(ctx, repository.CheckRecipientAddressParams{
			RecipientAddress: toDecode.RecipientAddr,
			UserID:           uid,
		})
		if i > 0 || fraudDetected {
			log.Printf("%v: transfer %s", ErrFraudDetection, tr.ID.String())

			if tr.ID != uuid.Nil {
				if err := s.wr.UpdateTokenTransfer(ctx, repository.UpdateTokenTransferParams{
					ID:     tr.ID,
					Status: TokenTransferStatusFraud,
				}); err != nil {
					log.Printf("could not update token transfer: %v", err)
				}
			}

			return ErrTransactionFailed
		}
	}

	tx, err := s.execTransfer(ctx, walletID, toDecode.RecipientAddr, toDecode.Amount, &lib_solana.SendAssetsConfig{
		PercentToCharge:           s.tokenTransferPercent,
		ChargeSolanaFeeFromSender: true,
		AllowFallbackToDefaultFee: true,
		DefaultFee:                1,
	})
	if err != nil {
		if tr.ID != uuid.Nil {
			if err := s.wr.UpdateTokenTransfer(ctx, repository.UpdateTokenTransferParams{
				ID:     tr.ID,
				Status: TokenTransferStatusFailed,
			}); err != nil {
				log.Printf("could not update token transfer: %v", err)
			}
		}

		return fmt.Errorf("could not confirm transfer: %w", err)
	}

	if tr.ID != uuid.Nil {
		if err := s.wr.UpdateTokenTransfer(ctx, repository.UpdateTokenTransferParams{
			ID:     tr.ID,
			Status: TokenTransferStatusSuccess,
			TxHash: sql.NullString{String: tx, Valid: true},
		}); err != nil {
			log.Printf("could not update token transfer: %v", err)
		}
	}

	return nil
}

func (s *Service) execTransfer(ctx context.Context, walletID uuid.UUID, recipientAddr string, amount float64, cfg *lib_solana.SendAssetsConfig) (string, error) {
	wallet, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return "", fmt.Errorf("could not get solana account: %w", err)
	}

	solanaAcc, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return "", fmt.Errorf("could not get solana account: %w", err)
	}

	balance, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, solanaAcc.PublicKey)
	if err != nil {
		return "", fmt.Errorf("could not get current balance: %w", err)
	}

	if balance < amount {
		return "", ErrNotEnoughBalance
	}

	feePayer, err := s.sc.AccountFromPrivateKeyBytes(s.feePayerSolanaPrivateKey)
	if err != nil {
		return "", err
	}
	source, err := s.sc.AccountFromPrivateKeyBytes(solanaAcc.PrivateKey)
	if err != nil {
		return "", err
	}

	for i := 0; i < 5; i++ {
		if tx, err := s.sc.SendAssetsWithAutoDerive(
			ctx,
			s.satorAssetSolanaAddr,
			feePayer,
			source,
			recipientAddr,
			amount,
			cfg,
		); err != nil {
			if i < 4 {
				log.Printf("can't send assets with auto derive: %v, fee accumulator address: %v\n", err, s.sc.FeeAccumulatorAddress())
			} else {
				log.Error(errors.Wrap(err, "can't send transaction"))
				return "", lib_errors.ErrCantSendSolanaTransaction
			}
			time.Sleep(time.Second * 5)
		} else {
			log.Printf("successful transaction: %s", tx)
			return tx, nil
		}
	}

	return "", ErrTransactionFailed
}

// GetStake Mocked method for stake
func (s *Service) GetStake(ctx context.Context, userID uuid.UUID) (Stake, error) {
	var stake repository.Stake
	stake, err := s.wr.GetStakeByUserID(ctx, userID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Stake{}, fmt.Errorf("could not get locked tokens by user id: %w", err)
		}

		stake.StakeAmount = 0
	}

	totalStake, err := s.wr.GetTotalStake(ctx)
	if err != nil {
		return Stake{}, fmt.Errorf("could not get total locked tokens: %w", err)
	}

	multiplier, err := s.GetMultiplier(ctx, userID)
	if err != nil {
		return Stake{}, fmt.Errorf("could not get miltiplier: %w", err)
	}

	w, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return Stake{}, fmt.Errorf("could not get sao wallet: %w", err)
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return Stake{}, fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return Stake{}, fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey)
	if err != nil {
		return Stake{}, fmt.Errorf("could not get sao balance: %w", err)
	}

	return Stake{
		TotalLocked:       totalStake,
		LockedByYou:       stake.StakeAmount,
		CurrentMultiplier: multiplier,
		AvailableToLock:   bal,
	}, nil
}

// SetStake method for set stake
// duration is in days
func (s *Service) SetStake(ctx context.Context, userID, walletID uuid.UUID, duration int64, amount float64) (bool, error) {
	feePayer, err := s.sc.AccountFromPrivateKeyBytes(s.feePayerSolanaPrivateKey)
	if err != nil {
		return false, err
	}
	stakePool := common.PublicKeyFromString(s.stakePoolSolanaPublicKey)
	asset := common.PublicKeyFromString(s.satorAssetSolanaAddr)

	if duration < 1 {
		duration = 30
	}

	lockDuration := time.Hour * 24 * time.Duration(duration)

	wallet, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return false, fmt.Errorf("could not get wallet: %w", err)
	}

	solanaAccount, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return false, fmt.Errorf("could not get wallet: %w", err)
	}

	bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, solanaAccount.PublicKey)
	if err != nil {
		return false, fmt.Errorf("could not get wallet balance")
	}

	level, err := s.wr.GetMinimalStakeLevel(ctx)
	if err != nil {
		return false, fmt.Errorf("could not get minimal tokens lock level")
	}

	if bal < amount || bal < level.MinStakeAmount.Float64 {
		return false, fmt.Errorf("insufficient balance amount. You have %.5f SAO", bal)
	}

	userWallet, err := types.AccountFromBytes(solanaAccount.PrivateKey)
	if err != nil {
		return false, err
	}

	var tx string
	for i := 0; i < 5; i++ {
		newCtx, cancel := context.WithCancel(context.Background())
		if tx, err = s.sc.Stake(newCtx, feePayer, userWallet, stakePool, asset, duration, amount); err != nil {
			if i < 4 {
				log.Println(err)
			} else {
				cancel()
				return false, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
			}
			cancel()
			time.Sleep(time.Second * 10)
		} else {
			cancel()
			log.Printf("successful transaction: %s", tx)
			break
		}
	}

	for i := 0; i < 3; i++ {
		ok, err := s.sc.IsTransactionSuccessful(ctx, tx)
		if err != nil {
			return false, fmt.Errorf("%w: %v", ErrTransactionFailed, err)
		}
		if ok {
			break
		}
		time.Sleep(time.Second * 10)
	}

	// Store stake data in our db.
	staked, err := s.wr.GetStakeByUserID(ctx, userID)
	if err != nil {
		if db.IsNotFoundError(err) {
			_, err := s.wr.AddStake(ctx, repository.AddStakeParams{
				UserID:      userID,
				WalletID:    walletID,
				StakeAmount: amount,
				StakeDuration: sql.NullInt32{
					Int32: int32(duration),
					Valid: true,
				},
				UnstakeDate:      time.Now().Add(lockDuration),
				UnstakeTimestamp: time.Now().Add(lockDuration).Unix(),
			})
			if err != nil {
				log.Printf("could not add stake to db for user= %v, %v", userID, err)
				return false, ErrCouldNotLock
			}

			return true, nil
		}

		log.Printf("could not add stake to db for user= %v, %v", userID, err)
		return false, ErrCouldNotLock
	}

	err = s.wr.UpdateStake(ctx, repository.UpdateStakeParams{
		UserID:      userID,
		StakeAmount: staked.StakeAmount + amount,
		StakeDuration: sql.NullInt32{
			Int32: int32(duration),
			Valid: true,
		},
		UnstakeDate:      time.Now().Add(lockDuration),
		UnstakeTimestamp: time.Now().Add(lockDuration).Unix(),
	})
	if err != nil {
		log.Printf("could not update stake in db for user= %v, %v", userID, err)
		return true, ErrCouldNotLock
	}

	return true, nil
}

// Unstake method for unstake
func (s *Service) Unstake(ctx context.Context, userID, walletID uuid.UUID) error {
	feePayer, err := s.sc.AccountFromPrivateKeyBytes(s.feePayerSolanaPrivateKey)
	if err != nil {
		return err
	}
	stakePool := common.PublicKeyFromString(s.stakePoolSolanaPublicKey)
	asset := common.PublicKeyFromString(s.satorAssetSolanaAddr)

	wallet, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return fmt.Errorf("could not get wallet: %w", err)
	}

	solanaAccount, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return fmt.Errorf("could not get wallet: %w", err)
	}

	userWallet, err := types.AccountFromBytes(solanaAccount.PrivateKey)
	if err != nil {
		return err
	}

	stake, err := s.wr.GetStakeByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not get wallet lock pool: %w", err)
	}

	unstakeDate := stake.CreatedAt.Add(time.Hour * 24 * time.Duration(stake.StakeDuration.Int32))
	if stake.UpdatedAt.Valid {
		unstakeDate = stake.UpdatedAt.Time.Add(time.Hour * 24 * time.Duration(stake.StakeDuration.Int32))
	}

	if time.Now().Before(unstakeDate) {
		return fmt.Errorf("unlock time has not yet come, unlock will be availabe at: %s", unstakeDate.String())
	}

	for i := 0; i < 5; i++ {
		newCtx, cancel := context.WithCancel(context.Background())
		if tx, err := s.sc.Unstake(newCtx, feePayer, userWallet, stakePool, asset); err != nil {
			if i < 4 {
				log.Println(err)
			} else {
				cancel()
				return fmt.Errorf("transaction: %w", err)
			}
			cancel()
			time.Sleep(time.Second * 10)
		} else {
			cancel()
			log.Printf("successful transaction: %s", tx)
			break
		}
	}

	err = s.wr.DeleteStakeByUserID(ctx, userID)
	if err != nil {
		log.Printf("could not delete stake by user id: %v", err)
		return nil
	}

	return nil
}

func castSolanaTxToTransaction(tx lib_solana.ConfirmedTransactionResponse, walletID uuid.UUID) Transaction {
	return Transaction{
		ID:        tx.TxHash,
		WalletID:  walletID.String(),
		TxHash:    tx.TxHash,
		Amount:    tx.Amount,
		CreatedAt: tx.CreatedAt.Format(time.RFC3339),
	}
}

// PayForService draft
func (s *Service) PayForService(ctx context.Context, uid uuid.UUID, amount float64, info string) error {
	w, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     uid,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return fmt.Errorf("could not make payment for %s: %w", info, err)
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey)
	if err != nil {
		return fmt.Errorf("could not get wallet balance")
	}

	if bal < amount {
		return fmt.Errorf("not enough balance for payment: %v", bal)
	}

	if _, err := s.execTransfer(ctx, w.ID, s.tokenHolderSolanaAddr, amount, &lib_solana.SendAssetsConfig{}); err != nil {
		return fmt.Errorf("could not make payment for %s: %w", info, err)
	}

	return nil
}

// PayForNFT draft
func (s *Service) PayForNFT(ctx context.Context, uid uuid.UUID, amount float64, info string, creatorAddr string, creatorShare int32) error {
	w, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     uid,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return fmt.Errorf("could not make payment for %s: %w", info, err)
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey)
	if err != nil {
		return fmt.Errorf("could not get wallet balance")
	}

	if bal < amount {
		return fmt.Errorf("not enough balance for payment: %v", bal)
	}

	if creatorShare > 100 {
		creatorShare = 100
	}

	satorShare := amount

	if creatorShare > 0 && creatorAddr != "" {
		creatorAmount := amount / 100 * float64(creatorShare)

		if _, err := s.execTransfer(ctx, w.ID, creatorAddr, creatorAmount, &lib_solana.SendAssetsConfig{}); err != nil {
			return fmt.Errorf("could not make payment for %s: %w", info, err)
		}

		satorShare = amount - creatorAmount
	}

	if satorShare > 0 {
		if _, err := s.execTransfer(ctx, w.ID, s.tokenHolderSolanaAddr, satorShare, &lib_solana.SendAssetsConfig{}); err != nil {
			return fmt.Errorf("could not make payment for %s: %w", info, err)
		}
	}

	return nil
}

// P2PTransfer draft
func (s *Service) P2PTransfer(ctx context.Context, uid, recipientID uuid.UUID, amount float64, cfg *lib_solana.SendAssetsConfig, info string) error {
	w, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     uid,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return fmt.Errorf("could not make payment for %s: %w", info, err)
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	bal, err := s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey)
	if err != nil {
		return fmt.Errorf("could not get wallet balance")
	}

	if bal < amount {
		return fmt.Errorf("not enough balance for payment: %v", bal)
	}

	wr, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     recipientID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return fmt.Errorf("could not get wallet by recipient id %s: %w", info, err)
	}

	sar, err := s.wr.GetSolanaAccountByID(ctx, wr.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("%w solana recipient account for this wallet", ErrNotFound)
		}
		return fmt.Errorf("could not get recipient solana account for this wallet: %w", err)
	}

	if _, err := s.execTransfer(ctx, w.ID, sar.PublicKey, amount, cfg); err != nil {
		return fmt.Errorf("could not make payment for %s: %w", info, err)
	}

	return nil
}

// GetMultiplier returns multiplier according to wallet's stake level.
func (s *Service) GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error) {
	stake, err := s.wr.GetStakeByUserID(ctx, userID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return 0, nil
		}

		return 0, fmt.Errorf("could not get stake by user id: %w", err)
	}

	lvl, err := s.wr.GetStakeLevelByAmount(ctx, stake.StakeAmount)
	if err != nil {
		if db.IsNotFoundError(err) {
			return 0, nil
		}

		return 0, fmt.Errorf("could not get user's stake level: %w", err)
	}

	return lvl.Multiplier.Int32, nil
}

// PossibleMultiplier returns multiplier that will be applied to user in stake will be increased on additionalAmount value.
func (s *Service) PossibleMultiplier(ctx context.Context, additionalAmount float64, userID, walletID uuid.UUID) (int32, error) {
	amount := additionalAmount

	staked, err := s.wr.GetStakeByUserID(ctx, userID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return 0, fmt.Errorf("could not get staked amount: %w", err)
		}
	} else {
		amount += staked.StakeAmount
	}

	lvl, err := s.wr.GetStakeLevelByAmount(ctx, amount)
	if err != nil {
		if db.IsNotFoundError(err) {
			return 0, nil
		}

		return 0, fmt.Errorf("could not get user's stake level: %w", err)
	}

	return lvl.Multiplier.Int32, nil
}

// GetEnabledStakeLevelsList returns enabled stake levels list
func (s *Service) GetEnabledStakeLevelsList(ctx context.Context, userID uuid.UUID) ([]StakeLevel, error) {
	lvls, err := s.wr.GetAllEnabledStakeLevels(ctx)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, err
		}
	}

	var currentStakeLvlID uuid.UUID
	if stake, err := s.wr.GetStakeByUserID(ctx, userID); err == nil {
		lvl, err := s.wr.GetStakeLevelByAmount(ctx, stake.StakeAmount)
		if err == nil {
			currentStakeLvlID = lvl.ID
		}
	}

	levels := make([]StakeLevel, 0, len(lvls))
	for _, l := range lvls {
		levels = append(levels, StakeLevel{
			ID:             l.ID.String(),
			MinAmount:      l.MinStakeAmount.Float64,
			MinDaysToStake: int(l.MinDaysAmount.Int32),
			Title:          l.Title,
			SubTitle:       l.Subtitle,
			IsCurrent:      currentStakeLvlID == l.ID,
			Rewards:        fmt.Sprintf("+%d", l.Multiplier.Int32) + "%",
		})
	}

	return levels, nil
}

// GetSAOBalance returns SAO account balance or error if it's not found.
func (s *Service) GetSAOBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	w, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return 0, fmt.Errorf("could not get wallet by user id %s: %w", userID, err)
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return 0, fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return 0, fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	return s.sc.GetTokenAccountBalanceWithAutoDerive(ctx, s.satorAssetSolanaAddr, sa.PublicKey)
}

func (s *Service) GetUserSolanaAccount(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	w, err := s.wr.GetWalletByUserIDAndType(ctx, repository.GetWalletByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get wallet by user id %s: %w", userID, err)
	}

	acc, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}

		return nil, fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	return acc.PrivateKey, nil
}
