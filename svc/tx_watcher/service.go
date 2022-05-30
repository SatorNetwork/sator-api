package tx_watcher

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"github.com/robfig/cron/v3"

	txw_repository "github.com/SatorNetwork/sator-api/svc/tx_watcher/repository"
)

type alias uint8

const (
	undefinedAlias alias = iota
	feePayerAlias
	tokenHolderAlias
)

func newAliasFromString(s string) (alias, error) {
	switch s {
	case "fee_payer":
		return feePayerAlias, nil
	case "token_holder":
		return tokenHolderAlias, nil
	default:
		return undefinedAlias, errors.Errorf("alias with such name %v doesn't exist", s)
	}
}

func (a alias) String() string {
	switch a {
	case undefinedAlias:
		return "undefined"
	case feePayerAlias:
		return "fee_payer"
	case tokenHolderAlias:
		return "token_holder"
	default:
		return "undefined"
	}
}

type aliases []alias

func newAliasesFromStrings(strings []string) (aliases, error) {
	aliases := make(aliases, 0, len(strings))
	for _, s := range strings {
		alias, err := newAliasFromString(s)
		if err != nil {
			return nil, errors.Wrap(err, "can't get new alias from string")
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

func (as aliases) ToStrings() []string {
	strings := make([]string, 0, len(as))
	for _, a := range as {
		strings = append(strings, a.String())
	}

	return strings
}

type status uint8

const (
	undefinedStatus status = iota
	registeredStatus
	successfulStatus
)

func (s status) String() string {
	switch s {
	case undefinedStatus:
		return "undefined"
	case registeredStatus:
		return "registered"
	case successfulStatus:
		return "successful"
	default:
		return "undefined"
	}
}

type (
	Service struct {
		txwr txwRepository
		sc   solanaClient

		feePayer    types.Account
		tokenHolder types.Account
	}

	txwRepository interface {
		GetTransactionsByStatus(ctx context.Context, status string) ([]txw_repository.WatcherTransaction, error)
		RegisterTransaction(ctx context.Context, arg txw_repository.RegisterTransactionParams) (txw_repository.WatcherTransaction, error)
		UpdateTransaction(ctx context.Context, arg txw_repository.UpdateTransactionParams) error
		UpdateTransactionStatus(ctx context.Context, arg txw_repository.UpdateTransactionStatusParams) error
	}

	solanaClient interface {
		SendConstructedTransaction(ctx context.Context, tx types.Transaction) (string, error)
		IsTransactionSuccessful(ctx context.Context, txhash string) (bool, error)
		GetBlockHeight(ctx context.Context) (uint64, error)
		GetLatestBlockhash(ctx context.Context) (rpc.GetLatestBlockhashValue, error)
	}
)

func NewService(
	txwr txwRepository,
	sc solanaClient,
	feePayer types.Account,
	tokenHolder types.Account,
) *Service {
	s := &Service{
		txwr:        txwr,
		sc:          sc,
		feePayer:    feePayer,
		tokenHolder: tokenHolder,
	}

	return s
}

func (s *Service) accountByAlias(a alias) (types.Account, error) {
	switch a {
	case undefinedAlias:
		return types.Account{}, errors.Errorf("alias %v is undefined", a)
	case feePayerAlias:
		return s.feePayer, nil
	case tokenHolderAlias:
		return s.tokenHolder, nil
	default:
		return types.Account{}, errors.Errorf("alias %v is undefined", a)
	}
}

func (s *Service) accountsByAliases(aliases []alias) ([]types.Account, error) {
	accounts := make([]types.Account, 0, len(aliases))
	for _, alias := range aliases {
		account, err := s.accountByAlias(alias)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *Service) SendAndWatchTx(ctx context.Context, message types.Message, accountAliases []alias) error {
	serializedMessage, err := message.Serialize()
	if err != nil {
		return errors.Wrap(err, "can't serialize message")
	}

	resp, err := s.sendSolanaTx(ctx, string(serializedMessage), aliases(accountAliases).ToStrings())
	if err != nil {
		return err
	}

	_, err = s.txwr.RegisterTransaction(ctx, txw_repository.RegisterTransactionParams{
		SerializedMessage:      string(serializedMessage),
		LatestValidBlockHeight: int64(resp.LatestValidBlockHeight),
		AccountAliases:         aliases(accountAliases).ToStrings(),
		TxHash:                 resp.TxHash,
		Status:                 registeredStatus.String(),
	})
	if err != nil {
		return errors.Wrap(err, "can't register transaction")
	}

	return nil
}

func (s *Service) resendSolanaDBTX(ctx context.Context, tx txw_repository.WatcherTransaction) error {
	resp, err := s.sendSolanaTx(ctx, tx.SerializedMessage, tx.AccountAliases)
	if err != nil {
		return err
	}

	err = s.txwr.UpdateTransaction(ctx, txw_repository.UpdateTransactionParams{
		ID:                     tx.ID,
		LatestValidBlockHeight: int64(resp.LatestValidBlockHeight),
		TxHash:                 resp.TxHash,
	})
	if err != nil {
		return errors.Wrap(err, "can't update transaction")
	}

	return nil
}

type sendSolanaTxResp struct {
	TxHash                 string
	LatestValidBlockHeight uint64
}

func (s *Service) sendSolanaTx(ctx context.Context, serializedMessage string, accountAliases []string) (*sendSolanaTxResp, error) {
	message, err := types.MessageDeserialize([]byte(serializedMessage))
	if err != nil {
		return nil, errors.Wrap(err, "can't deserialize message")
	}
	latestBlockhash, err := s.sc.GetLatestBlockhash(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get latest blockhash")
	}
	message.RecentBlockHash = latestBlockhash.Blockhash

	aliases, err := newAliasesFromStrings(accountAliases)
	if err != nil {
		return nil, errors.Wrap(err, "can't get new aliases from strings")
	}
	accounts, err := s.accountsByAliases(aliases)
	if err != nil {
		return nil, errors.Wrap(err, "can't get accounts by aliases")
	}
	solanaTx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: accounts,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create new transaction")
	}
	txHash, err := s.sc.SendConstructedTransaction(ctx, solanaTx)
	if err != nil {
		return nil, errors.Wrap(err, "can't send transaction")
	}

	return &sendSolanaTxResp{
		TxHash:                 txHash,
		LatestValidBlockHeight: latestBlockhash.LatestValidBlockHeight,
	}, nil
}

func (s *Service) start() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		if err := s.resendSolanaDBTXsIfNeeded(context.Background()); err != nil {
			log.Printf("can't resend solana DBTXs: %v", err)
		}
	})
	if err != nil {
		log.Printf("can't register resend-solana-dbtxs-if-needed callback")
	}

	c.Start()
}

func (s *Service) resendSolanaDBTXsIfNeeded(ctx context.Context) error {
	txs, err := s.txwr.GetTransactionsByStatus(ctx, registeredStatus.String())
	if err != nil {
		return errors.Wrap(err, "can't get transactions by status")
	}

	for _, tx := range txs {
		if err := s.processTx(ctx, tx); err != nil {
			log.Printf("can't process tx: %v\n", err)
			continue
		}
	}

	return nil
}

func (s *Service) processTx(ctx context.Context, tx txw_repository.WatcherTransaction) error {
	success, err := s.sc.IsTransactionSuccessful(ctx, tx.TxHash)
	if err != nil {
		return errors.Wrap(err, "can't check if transaction is successful")
	}
	if success {
		err := s.txwr.UpdateTransactionStatus(ctx, txw_repository.UpdateTransactionStatusParams{
			ID:     tx.ID,
			Status: successfulStatus.String(),
		})
		if err != nil {
			return errors.Wrap(err, "can't update transaction status")
		}
		return nil
	}

	needToRetry, err := s.needToRetry(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "can't check if tx need to be retried")
	}
	if !needToRetry {
		return nil
	}

	err = s.resendSolanaDBTX(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) needToRetry(ctx context.Context, tx txw_repository.WatcherTransaction) (bool, error) {
	cbh, err := s.sc.GetBlockHeight(ctx)
	if err != nil {
		return false, errors.Wrap(err, "can't get block height")
	}

	return int64(cbh) > tx.LatestValidBlockHeight, nil
}
