package solana

import (
	"context"
	"time"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

//go:generate mockgen -destination=mock_client.go -package=solana github.com/SatorNetwork/sator-api/lib/solana Interface
type Interface interface {
	Endpoint() string
	IssueAsset(ctx context.Context, feePayer, issuer, asset types.Account, dest common.PublicKey, amount float64) (string, error)
	CreateAccountWithATA(ctx context.Context, assetAddr, initAccAddr string, feePayer types.Account) (string, error)
	DeriveATAPublicKey(ctx context.Context, recipientPK, assetPK common.PublicKey) (common.PublicKey, error)
	GetConfirmedTransaction(ctx context.Context, txhash string) (GetConfirmedTransactionResponse, error)
	GetConfirmedTransactionForAccount(ctx context.Context, assetAddr string, rootPubKey string, txhash string) (ConfirmedTransactionResponse, error)
	IsTransactionSuccessful(ctx context.Context, txhash string) (bool, error)
	NeedToRetry(ctx context.Context, latestValidBlockHeight int64) (bool, error)
	GetBlockHeight(ctx context.Context) (uint64, error)
	NewAccount() types.Account
	PublicKeyFromString(pk string) common.PublicKey
	AccountFromPrivateKeyBytes(pk []byte) (types.Account, error)
	CheckPrivateKey(addr string, pk []byte) error
	FeeAccumulatorAddress() string
	RequestAirdrop(ctx context.Context, pubKey string, amount float64) (string, error)
	SendConstructedTransaction(ctx context.Context, tx types.Transaction) (string, error)
	SendTransaction(ctx context.Context, feePayer, signer types.Account, instructions ...types.Instruction) (string, error)
	GetAccountBalanceSOL(ctx context.Context, accPubKey string) (float64, error)
	GetTokenAccountBalance(ctx context.Context, accPubKey string) (float64, error)
	GetTokenAccountBalanceWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (float64, error)
	GetTransactions(ctx context.Context, assetAddr, rootPubKey, ataPubKey string) (txList []ConfirmedTransactionResponse, err error)
	GetTransactionsWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (txList []ConfirmedTransactionResponse, err error)
	CreateAsset(ctx context.Context, feePayer, issuer, asset types.Account) (string, error)
	InitAccountToUseAsset(ctx context.Context, feePayer, issuer, asset, initAcc types.Account) (string, error)
	GiveAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, issuer types.Account, recipientAddr string, amount float64) (string, error)
	PrepareSendAssetsTx(
		ctx context.Context,
		assetAddr string,
		feePayer types.Account,
		source types.Account,
		recipientAddr string,
		amount float64,
		cfg *SendAssetsConfig,
	) (*PrepareTxResponse, error)
	SendAssetsWithAutoDerive(
		ctx context.Context,
		assetAddr string,
		feePayer types.Account,
		source types.Account,
		recipientAddr string,
		amount float64,
		cfg *SendAssetsConfig,
	) (string, error)
	TransactionDeserialize(tx []byte) (types.Transaction, error)
	SerializeTxMessage(message types.Message) ([]byte, error)
	DeserializeTxMessage(message []byte) (types.Message, error)
	NewTransaction(param types.NewTransactionParam) (types.Transaction, error)
	GetLatestBlockhash(ctx context.Context) (rpc.GetLatestBlockhashValue, error)
	GetNFTsByWalletAddress(ctx context.Context, walletAddr string) ([]*ArweaveNFTMetadata, error)
	GetNFTMintAddrs(ctx context.Context, walletAddr string) ([]string, error)
	GetNFTMetadata(mintAddr string) (*ArweaveNFTMetadata, error)
	InitializeStakePool(ctx context.Context, feePayer, issuer types.Account, asset common.PublicKey) (txHast string, stakePool types.Account, err error)
	Stake(ctx context.Context, feePayer, userWallet types.Account, pool, asset common.PublicKey, duration int64, amount float64) (string, error)
	Unstake(ctx context.Context, feePayer, userWallet types.Account, stakePool, asset common.PublicKey) (string, error)
}

type (
	// GetConfirmedTransactionResponse ...
	GetConfirmedTransactionResponse struct {
		BlockTime   int64             `json:"blockTime"`
		Slot        uint64            `json:"slot"`
		Meta        TransactionMeta   `json:"meta"`
		Transaction types.Transaction `json:"transaction"`
	}

	// TransactionMeta ...
	TransactionMeta struct {
		Fee               uint64   `json:"fee"`
		PreBalances       []int64  `json:"preBalances"`
		PostBalances      []int64  `json:"postBalances"`
		LogMessages       []string `json:"logMesssages"`
		InnerInstructions []struct {
			Index        uint64              `json:"index"`
			Instructions []types.Instruction `json:"instructions"`
		} `json:"innerInstructions"`
		Err    interface{}            `json:"err"`
		Status map[string]interface{} `json:"status"`

		// custom fields
		PostTokenBalances []TokenBalance `json:"postTokenBalances"`
		PreTokenBalances  []TokenBalance `json:"preTokenBalances"`
	}

	// TokenBalance ...
	TokenBalance struct {
		AccountIndex  int           `json:"accountIndex"`
		Mint          string        `json:"mint"`
		UITokenAmount UITokenAmount `json:"uiTokenAmount"`
	}

	// UITokenAmount ...
	UITokenAmount struct {
		Amount         string  `json:"amount"`
		Decimals       int     `json:"decimals"`
		UIAmount       float64 `json:"uiAmount"`
		UIAmountString string  `json:"uiAmountString"`
	}

	// ConfirmedTransactionResponse ...
	ConfirmedTransactionResponse struct {
		TxHash        string    `json:"tx_hash"`
		Amount        float64   `json:"amount"`
		AmountString  string    `json:"amount_string"`
		CreatedAtUnix int64     `json:"created_at_unix"`
		CreatedAt     time.Time `json:"created_at"`
	}

	PrepareTxResponse struct {
		Tx                      types.Transaction
		FeeInSAO                float64
		BlockchainFeeInSOLMltpl uint64
	}

	SendAssetsConfig struct {
		PercentToCharge           float64
		ChargeSolanaFeeFromSender bool
		AllowFallbackToDefaultFee bool
		DefaultFee                uint64
	}

	ArweaveNFTMetadata struct {
		Name                 string `json:"name"`
		Symbol               string `json:"symbol"`
		Description          string `json:"description"`
		SellerFeeBasisPoints int    `json:"seller_fee_basis_points"`
		Image                string `json:"image"`
		Attributes           []struct {
			TraitType string      `json:"trait_type"`
			Value     interface{} `json:"value"`
		} `json:"attributes"`
		Collection struct {
			Name   string `json:"name"`
			Family string `json:"family"`
		} `json:"collection"`
		Properties struct {
			Files []struct {
				Uri  string `json:"uri"`
				Type string `json:"type"`
			} `json:"files"`
			Category string `json:"category"`
			Creators []struct {
				Address string `json:"address"`
				Share   int    `json:"share"`
			} `json:"creators"`
		} `json:"properties"`
	}
)
