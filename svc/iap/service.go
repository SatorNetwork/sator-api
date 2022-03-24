package iap

import (
	"context"
	"time"

	"github.com/awa/go-iap/appstore"
	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/types"

	lib_appstore "github.com/SatorNetwork/sator-api/lib/appstore"
	appstore_client "github.com/SatorNetwork/sator-api/lib/appstore/client"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

const (
	maxRetries      = 5
	constantBackOff = 10 * time.Second

	testProductID   = "test2"
	testTokenAmount = 1
)

type (
	Service struct {
		client lib_appstore.Interface
		wr     walletRepository
		sc     solanaClient

		satorAssetSolanaAddr string
		feePayer             types.Account
		tokenHolder          types.Account
	}

	walletRepository interface {
		GetSolanaAccountByUserIDAndType(
			ctx context.Context,
			arg repository.GetSolanaAccountByUserIDAndTypeParams,
		) (repository.SolanaAccount, error)
	}

	solanaClient interface {
		GiveAssetsWithAutoDerive(
			ctx context.Context,
			assetAddr string,
			feePayer types.Account,
			issuer types.Account,
			recipientAddr string,
			amount float64,
		) (string, error)
	}

	Empty struct{}

	RegisterInAppPurchaseRequest struct {
		ReceiptData string `json:"receipt_data"`
	}
)

func NewService(
	wr walletRepository,
	sc solanaClient,
	satorAssetSolanaAddr string,
	feePayer types.Account,
	tokenHolder types.Account,
) *Service {
	s := &Service{
		client: appstore_client.New(),
		wr:     wr,
		sc:     sc,

		satorAssetSolanaAddr: satorAssetSolanaAddr,
		feePayer:             feePayer,
		tokenHolder:          tokenHolder,
	}

	return s
}

func (s *Service) RegisterInAppPurchase(ctx context.Context, userID uuid.UUID, req *RegisterInAppPurchaseRequest) (*Empty, error) {
	iapReq := appstore.IAPRequest{
		ReceiptData:            req.ReceiptData,
		Password:               "",
		ExcludeOldTransactions: false,
	}
	iapResp := &appstore.IAPResponse{}
	if err := s.client.Verify(ctx, iapReq, iapResp); err != nil {
		return nil, errors.Wrap(err, "can't verify in-app-purchase request")
	}
	if err := appstore.HandleError(iapResp.Status); err != nil {
		return nil, errors.Wrap(err, "invalid in-app-purchase status code")
	}

	solanaAccount, err := s.wr.GetSolanaAccountByUserIDAndType(ctx, repository.GetSolanaAccountByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: wallet.WalletTypeSator,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get solana account by user id and type")
	}

	for _, inApp := range iapResp.Receipt.InApp {
		switch inApp.ProductID {
		case testProductID:
			err := backoff.Retry(func() error {
				_, err = s.sc.GiveAssetsWithAutoDerive(
					ctx,
					s.satorAssetSolanaAddr,
					s.feePayer,
					s.tokenHolder,
					solanaAccount.PublicKey,
					testTokenAmount,
				)
				if err != nil {
					return errors.Wrap(err, "can't send sator tokens")
				}
				return nil
			}, backoff.WithMaxRetries(backoff.NewConstantBackOff(constantBackOff), maxRetries))
			if err != nil {
				return nil, err
			}

		default:
			return nil, errors.Errorf("unknown product id: %v", inApp.ProductID)
		}
	}

	return &Empty{}, nil
}
