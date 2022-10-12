package iap

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/awa/go-iap/appstore"
	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/types"
	log "github.com/sirupsen/logrus"

	lib_appstore "github.com/SatorNetwork/sator-api/lib/appstore"
	appstore_client "github.com/SatorNetwork/sator-api/lib/appstore/client"
	lib_errors "github.com/SatorNetwork/sator-api/lib/errors"
	lib_nft_marketplace "github.com/SatorNetwork/sator-api/lib/nft_marketplace"
	"github.com/SatorNetwork/sator-api/svc/exchange_rates"
	iap_repository "github.com/SatorNetwork/sator-api/svc/iap/repository"
	"github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
)

const (
	maxRetries      = 100
	constantBackOff = 10 * time.Second
)

type (
	Service struct {
		client         lib_appstore.Interface
		nftMarketplace lib_nft_marketplace.Interface
		ir             iapRepository
		wr             walletRepository
		sc             solanaClient
		exchange_rates exchangeRatesClient

		satorAssetSolanaAddr string
		feePayer             types.Account
		tokenHolder          types.Account
	}

	exchangeRatesClient interface {
		GetAssetPrice(ctx context.Context, req *exchange_rates.Asset) (*exchange_rates.Price, error)
	}

	iapRepository interface {
		CreateIAPReceipt(ctx context.Context, arg iap_repository.CreateIAPReceiptParams) (iap_repository.IapReceipt, error)
		GetIAPReceiptByTxID(ctx context.Context, transactionID string) (iap_repository.IapReceipt, error)
		GetIapProductByID(ctx context.Context, id string) (iap_repository.IapProduct, error)
	}

	walletRepository interface {
		GetSolanaAccountByUserIDAndType(
			ctx context.Context,
			arg repository.GetSolanaAccountByUserIDAndTypeParams,
		) (repository.SolanaAccount, error)
	}

	solanaClient interface {
		AccountFromPrivateKeyBytes(pk []byte) (types.Account, error)
		GiveAssetsWithAutoDerive(
			ctx context.Context,
			assetAddr string,
			feePayer types.Account,
			issuer types.Account,
			recipientAddr string,
			amount float64,
		) (string, error)
		TransactionDeserialize(tx []byte) (types.Transaction, error)
		SerializeTxMessage(message types.Message) ([]byte, error)
	}

	Empty struct{}

	RegisterInAppPurchaseRequest struct {
		ReceiptData string `json:"receipt_data"`
		MintAddress string `json:"mint_address"`
	}
)

func NewService(
	nftMarketplace lib_nft_marketplace.Interface,
	ir iapRepository,
	wr walletRepository,
	sc solanaClient,
	satorAssetSolanaAddr string,
	feePayer types.Account,
	tokenHolder types.Account,
	exchange_rates exchangeRatesClient,
) *Service {
	s := &Service{
		client:         appstore_client.New(),
		nftMarketplace: nftMarketplace,
		ir:             ir,
		wr:             wr,
		sc:             sc,
		exchange_rates: exchange_rates,

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

	fmt.Printf("Processing receipt with %v transaction ID\n", purchase.TransactionID)

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

	iapProduct, err := s.ir.GetIapProductByID(ctx, purchase.ProductID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get iap product by id")
	}

	saoPrice, err := s.exchange_rates.GetAssetPrice(ctx, &exchange_rates.Asset{
		AssetType: exchange_rates.AssetTypeSAO,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get sao price")
	}

	saoAmount := (iapProduct.PriceInUsd / saoPrice.Usd) * 0.7

	ctxb := context.Background()
	err = backoff.Retry(func() error {
		_, err := s.sc.GiveAssetsWithAutoDerive(
			ctxb,
			s.satorAssetSolanaAddr,
			s.feePayer,
			s.tokenHolder,
			solanaAccount.PublicKey,
			saoAmount,
		)
		if err != nil {
			return errors.Wrap(err, "can't send sator tokens")
		}
		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(constantBackOff), maxRetries))
	if err != nil {
		log.Error(errors.Wrap(err, "can't send transaction"))
		return nil, lib_errors.ErrCantSendSolanaTransaction
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
		return nil, errors.Wrapf(err, "can't create IAP receipt with %v transaction ID", purchase.TransactionID)
	}

	log.Printf("Preparing buy tx, mint address: %v, charge tokens from: %v", req.MintAddress, solanaAccount.PublicKey)
	prepareBuyTxResp, err := s.nftMarketplace.PrepareBuyTx(&lib_nft_marketplace.PrepareBuyTxRequest{
		MintAddress:      req.MintAddress,
		ChargeTokensFrom: solanaAccount.PublicKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare buy tx")
	}

	nftBuyer, err := s.sc.AccountFromPrivateKeyBytes(solanaAccount.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "can't get account from private key bytes")
	}

	nftBuyerSignature, err := s.getNFTBuyerSignature(nftBuyer, prepareBuyTxResp.EncodedTx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get nft buyer signature")
	}

	log.Printf("Sending prepared buy tx, txid: %v", prepareBuyTxResp.PreparedBuyTxId)
	_, err = s.nftMarketplace.SendPreparedBuyTx(&lib_nft_marketplace.SendPreparedBuyTxRequest{
		PreparedBuyTxId: prepareBuyTxResp.PreparedBuyTxId,
		BuyerSignature:  base64.StdEncoding.EncodeToString(nftBuyerSignature),
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't send prepared buy tx")
	}

	return &Empty{}, nil
}

func (s *Service) getNFTBuyerSignature(nftBuyer types.Account, encodedTx string) ([]byte, error) {
	var nftBuyerSignature []byte
	{
		serializedTx, err := base64.StdEncoding.DecodeString(encodedTx)
		if err != nil {
			return nil, errors.Wrap(err, "can't decode transaction")
		}
		tx, err := s.sc.TransactionDeserialize(serializedTx)
		if err != nil {
			return nil, errors.Wrap(err, "can't deserialize transaction")
		}
		serializedMessage, err := s.sc.SerializeTxMessage(tx.Message)
		if err != nil {
			return nil, errors.Wrap(err, "can't serialize message")
		}
		nftBuyerSignature = nftBuyer.Sign(serializedMessage)
	}

	return nftBuyerSignature, nil
}
