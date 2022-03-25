package iap

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/awa/go-iap/appstore"
	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/types"

	lib_appstore "github.com/SatorNetwork/sator-api/lib/appstore"
	appstore_client "github.com/SatorNetwork/sator-api/lib/appstore/client"
	iap_repository "github.com/SatorNetwork/sator-api/svc/iap/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

const (
	maxRetries      = 5
	constantBackOff = 10 * time.Second

	testProductID   = "test2"
	testTokenAmount = 0
)

type (
	Service struct {
		client lib_appstore.Interface
		ir     iapRepository
		wr     walletRepository
		sc     solanaClient

		satorAssetSolanaAddr string
		feePayer             types.Account
		tokenHolder          types.Account
	}

	iapRepository interface {
		CreateIAPReceipt(ctx context.Context, arg iap_repository.CreateIAPReceiptParams) (iap_repository.IapReceipt, error)
		GetIAPReceiptByTxID(ctx context.Context, transactionID string) (iap_repository.IapReceipt, error)
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
	ir iapRepository,
	wr walletRepository,
	sc solanaClient,
	satorAssetSolanaAddr string,
	feePayer types.Account,
	tokenHolder types.Account,
) *Service {
	s := &Service{
		client: appstore_client.New(),
		ir:     ir,
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
	if len(iapResp.Receipt.InApp) != 1 {
		return nil, errors.Errorf("1 purchase is expected, got: %v", len(iapResp.Receipt.InApp))
	}
	purchase := iapResp.Receipt.InApp[0]

	solanaAccount, err := s.wr.GetSolanaAccountByUserIDAndType(ctx, repository.GetSolanaAccountByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: wallet.WalletTypeSator,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get solana account by user id and type")
	}

	_, err = s.ir.GetIAPReceiptByTxID(ctx, purchase.TransactionID)
	if err != nil && !strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
		return nil, errors.Wrap(err, "can't get iap receipt by txid")
	}
	if err == nil {
		return nil, errors.Errorf("receipt already processed")
	}

	ctxb := context.Background()
	switch purchase.ProductID {
	case testProductID:
		err := backoff.Retry(func() error {
			_, err = s.sc.GiveAssetsWithAutoDerive(
				ctxb,
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
		return nil, errors.Errorf("unknown product id: %v", purchase.ProductID)
	}

	receiptInJson, err := json.Marshal(iapResp)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal iap receipt")
	}
	_, err = s.ir.CreateIAPReceipt(ctxb, iap_repository.CreateIAPReceiptParams{
		TransactionID: purchase.TransactionID,
		ReceiptData:   req.ReceiptData,
		ReceiptInJson: string(receiptInJson),
		UserID:        userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create IAP receipt")
	}

	return &Empty{}, nil
}
