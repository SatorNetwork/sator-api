package solana_multiprovider

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
)

var (
	ErrSolanaProvidersDontRespond = errors.New("solana providers dont respond")
)

type solanaMultiProvider struct {
	providers []lib_solana.Interface
	m         *metricsRegistrator
}

func New(providers []lib_solana.Interface, mr metricsRepository) (lib_solana.Interface, error) {
	if len(providers) == 0 {
		return nil, errors.New("at least one solana provider should be specified")
	}

	return &solanaMultiProvider{
		providers: providers,
		m:         newMetricsRegistrator(mr),
	}, nil
}

func tryNextProvider(err error) bool {
	return strings.Contains(err.Error(), `{"jsonrpc":"2.0","error":{"code":503,"message":"Service unavailable"}`)
}

func (s *solanaMultiProvider) Endpoint() string {
	return "solana multiprovider"
}

func (s *solanaMultiProvider) IssueAsset(ctx context.Context, feePayer, issuer, asset types.Account, dest common.PublicKey, amount float64) (string, error) {
	for _, p := range s.providers {
		resp, err := p.IssueAsset(ctx, feePayer, issuer, asset, dest, amount)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) DeriveATAPublicKey(ctx context.Context, recipientPK, assetPK common.PublicKey) (common.PublicKey, error) {
	for _, p := range s.providers {
		resp, err := p.DeriveATAPublicKey(ctx, recipientPK, assetPK)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return common.PublicKey{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) CreateAccountWithATA(ctx context.Context, assetAddr, initAccAddr string, feePayer types.Account) (string, error) {
	for _, p := range s.providers {
		resp, err := p.CreateAccountWithATA(ctx, assetAddr, initAccAddr, feePayer)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetConfirmedTransaction(ctx context.Context, txhash string) (lib_solana.GetConfirmedTransactionResponse, error) {
	for _, p := range s.providers {
		resp, err := p.GetConfirmedTransaction(ctx, txhash)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return lib_solana.GetConfirmedTransactionResponse{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetConfirmedTransactionForAccount(ctx context.Context, assetAddr string, rootPubKey string, txhash string) (lib_solana.ConfirmedTransactionResponse, error) {
	for _, p := range s.providers {
		resp, err := p.GetConfirmedTransactionForAccount(ctx, assetAddr, rootPubKey, txhash)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return lib_solana.ConfirmedTransactionResponse{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) IsTransactionSuccessful(ctx context.Context, txhash string) (bool, error) {
	for _, p := range s.providers {
		resp, err := p.IsTransactionSuccessful(ctx, txhash)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return false, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) NeedToRetry(ctx context.Context, latestValidBlockHeight int64) (bool, error) {
	for _, p := range s.providers {
		resp, err := p.NeedToRetry(ctx, latestValidBlockHeight)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return false, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetBlockHeight(ctx context.Context) (uint64, error) {
	for _, p := range s.providers {
		resp, err := p.GetBlockHeight(ctx)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return 0, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) NewAccount() types.Account {
	return s.providers[0].NewAccount()
}

func (s *solanaMultiProvider) PublicKeyFromString(pk string) common.PublicKey {
	return s.providers[0].PublicKeyFromString(pk)
}

func (s *solanaMultiProvider) AccountFromPrivateKeyBytes(pk []byte) (types.Account, error) {
	for _, p := range s.providers {
		resp, err := p.AccountFromPrivateKeyBytes(pk)
		if err != nil && tryNextProvider(err) {
			continue
		}
		return resp, err
	}
	return types.Account{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) CheckPrivateKey(addr string, pk []byte) error {
	for _, p := range s.providers {
		err := p.CheckPrivateKey(addr, pk)
		if err != nil && tryNextProvider(err) {
			continue
		}
		return err
	}
	return ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) FeeAccumulatorAddress() string {
	return s.providers[0].FeeAccumulatorAddress()
}

func (s *solanaMultiProvider) RequestAirdrop(ctx context.Context, pubKey string, amount float64) (string, error) {
	for _, p := range s.providers {
		resp, err := p.RequestAirdrop(ctx, pubKey, amount)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) SendConstructedTransaction(ctx context.Context, tx types.Transaction) (string, error) {
	for _, p := range s.providers {
		resp, err := p.SendConstructedTransaction(ctx, tx)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) SendTransaction(ctx context.Context, feePayer, signer types.Account, instructions ...types.Instruction) (string, error) {
	for _, p := range s.providers {
		resp, err := p.SendTransaction(ctx, feePayer, signer, instructions...)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetAccountBalanceSOL(ctx context.Context, accPubKey string) (float64, error) {
	for _, p := range s.providers {
		resp, err := p.GetAccountBalanceSOL(ctx, accPubKey)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return 0, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetTokenAccountBalance(ctx context.Context, accPubKey string) (float64, error) {
	for _, p := range s.providers {
		resp, err := p.GetTokenAccountBalance(ctx, accPubKey)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return 0, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetTokenAccountBalanceWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (float64, error) {
	for _, p := range s.providers {
		resp, err := p.GetTokenAccountBalanceWithAutoDerive(ctx, assetAddr, accountAddr)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return 0, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetTransactions(ctx context.Context, assetAddr, rootPubKey, ataPubKey string) (txList []lib_solana.ConfirmedTransactionResponse, err error) {
	for _, p := range s.providers {
		resp, err := p.GetTransactions(ctx, assetAddr, rootPubKey, ataPubKey)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetTransactionsWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (txList []lib_solana.ConfirmedTransactionResponse, err error) {
	for _, p := range s.providers {
		resp, err := p.GetTransactionsWithAutoDerive(ctx, assetAddr, accountAddr)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) CreateAsset(ctx context.Context, feePayer, issuer, asset types.Account) (string, error) {
	for _, p := range s.providers {
		resp, err := p.CreateAsset(ctx, feePayer, issuer, asset)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) InitAccountToUseAsset(ctx context.Context, feePayer, issuer, asset, initAcc types.Account) (string, error) {
	for _, p := range s.providers {
		resp, err := p.InitAccountToUseAsset(ctx, feePayer, issuer, asset, initAcc)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GiveAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, issuer types.Account, recipientAddr string, amount float64) (string, error) {
	for _, p := range s.providers {
		resp, err := p.GiveAssetsWithAutoDerive(ctx, assetAddr, feePayer, issuer, recipientAddr, amount)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) PrepareSendAssetsTx(
	ctx context.Context,
	assetAddr string,
	feePayer types.Account,
	source types.Account,
	recipientAddr string,
	amount float64,
	cfg *lib_solana.SendAssetsConfig,
) (*lib_solana.PrepareTxResponse, error) {
	for _, p := range s.providers {
		resp, err := p.PrepareSendAssetsTx(ctx, assetAddr, feePayer, source, recipientAddr, amount, cfg)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) SendAssetsWithAutoDerive(
	ctx context.Context,
	assetAddr string,
	feePayer types.Account,
	source types.Account,
	recipientAddr string,
	amount float64,
	cfg *lib_solana.SendAssetsConfig,
) (string, error) {
	for _, p := range s.providers {
		resp, err := p.SendAssetsWithAutoDerive(ctx, assetAddr, feePayer, source, recipientAddr, amount, cfg)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) TransactionDeserialize(tx []byte) (types.Transaction, error) {
	for _, p := range s.providers {
		resp, err := p.TransactionDeserialize(tx)
		if err != nil && tryNextProvider(err) {
			continue
		}
		return resp, err
	}
	return types.Transaction{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) SerializeTxMessage(message types.Message) ([]byte, error) {
	for _, p := range s.providers {
		resp, err := p.SerializeTxMessage(message)
		if err != nil && tryNextProvider(err) {
			continue
		}
		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) DeserializeTxMessage(message []byte) (types.Message, error) {
	for _, p := range s.providers {
		resp, err := p.DeserializeTxMessage(message)
		if err != nil && tryNextProvider(err) {
			continue
		}
		return resp, err
	}
	return types.Message{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) NewTransaction(param types.NewTransactionParam) (types.Transaction, error) {
	for _, p := range s.providers {
		resp, err := p.NewTransaction(param)
		if err != nil && tryNextProvider(err) {
			continue
		}
		return resp, err
	}
	return types.Transaction{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetLatestBlockhash(ctx context.Context) (rpc.GetLatestBlockhashValue, error) {
	for _, p := range s.providers {
		resp, err := p.GetLatestBlockhash(ctx)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return rpc.GetLatestBlockhashValue{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetNFTsByWalletAddress(ctx context.Context, walletAddr string) ([]*lib_solana.ArweaveNFTMetadata, error) {
	for _, p := range s.providers {
		resp, err := p.GetNFTsByWalletAddress(ctx, walletAddr)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetNFTMintAddrs(ctx context.Context, walletAddr string) ([]string, error) {
	for _, p := range s.providers {
		resp, err := p.GetNFTMintAddrs(ctx, walletAddr)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) GetNFTMetadata(mintAddr string) (*lib_solana.ArweaveNFTMetadata, error) {
	for _, p := range s.providers {
		resp, err := p.GetNFTMetadata(mintAddr)
		if err != nil {
			s.m.registerError(context.Background(), p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(context.Background(), p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(context.Background(), p.Endpoint())
		} else {
			s.m.registerSuccessCall(context.Background(), p.Endpoint())
		}

		return resp, err
	}
	return nil, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) InitializeStakePool(ctx context.Context, feePayer, issuer types.Account, asset common.PublicKey) (txHash string, stakePool types.Account, err error) {
	for _, p := range s.providers {
		txHash, stakePool, err := p.InitializeStakePool(ctx, feePayer, issuer, asset)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return txHash, stakePool, err
	}
	return "", types.Account{}, ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) Stake(ctx context.Context, feePayer, userWallet types.Account, pool, asset common.PublicKey, duration int64, amount float64) (string, error) {
	for _, p := range s.providers {
		resp, err := p.Stake(ctx, feePayer, userWallet, pool, asset, duration, amount)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}

func (s *solanaMultiProvider) Unstake(ctx context.Context, feePayer, userWallet types.Account, stakePool, asset common.PublicKey) (string, error) {
	for _, p := range s.providers {
		resp, err := p.Unstake(ctx, feePayer, userWallet, stakePool, asset)
		if err != nil {
			s.m.registerError(ctx, p.Endpoint(), err.Error())
		}
		if err != nil && tryNextProvider(err) {
			s.m.registerNotAvailableError(ctx, p.Endpoint())
			continue
		}

		if err != nil {
			s.m.registerOtherError(ctx, p.Endpoint())
		} else {
			s.m.registerSuccessCall(ctx, p.Endpoint())
		}

		return resp, err
	}
	return "", ErrSolanaProvidersDontRespond
}
