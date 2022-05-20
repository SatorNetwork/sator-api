package worker

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/svc/rewards"
	"github.com/google/uuid"
	"github.com/portto/solana-go-sdk/client"

	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
)

type (
	TransactionWorker struct {
		repo rewardsRepository
		sc   solanaClient
		ws   walletService
	}

	rewardsRepository interface {
		GetRequestedTransactions(ctx context.Context) ([]repository.Reward, error)
		UpdateTransactionTxHash(ctx context.Context, arg repository.UpdateTransactionTxHashParams) error
		UpdateTransactionStatusByTxHash(ctx context.Context, arg repository.UpdateTransactionStatusByTxHashParams) error
	}

	solanaClient interface {
		CheckTransaction(ctx context.Context, txHash string) (bool, error)
		GetTransaction(ctx context.Context, txHash string) (*client.GetTransactionResponse, error)
	}

	walletService interface {
		WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (string, error)
	}
)

func New(repo rewardsRepository, sc solanaClient, ws walletService) *TransactionWorker {
	return &TransactionWorker{
		repo: repo,
		sc:   sc,
		ws:   ws,
	}
}

func (t *TransactionWorker) Start(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			transactions, err := t.repo.GetRequestedTransactions(ctx)
			if err != nil {
				log.Println(err)
				return
			}

			for i := range transactions {
				resp, err := t.sc.GetTransaction(ctx, transactions[i].TxHash.String)
				if err != nil {
					log.Println(err)
					continue
				}

				if resp != nil {
					if err = t.repo.UpdateTransactionStatusByTxHash(ctx, repository.UpdateTransactionStatusByTxHashParams{
						Status: rewards.TransactionStatusWithdrawn.String(),
						TxHash: transactions[i].TxHash,
					}); err != nil {
						log.Println(err)
					}
					continue
				}

				txHash, err := t.ws.WithdrawRewards(ctx, transactions[i].UserID, transactions[i].Amount)
				if err != nil {
					log.Println(err)
					continue
				}
				if err := t.repo.UpdateTransactionTxHash(ctx, repository.UpdateTransactionTxHashParams{
					TxHashNew: sql.NullString{
						String: txHash,
						Valid:  true,
					},
					TxHash: transactions[i].TxHash,
				}); err != nil {
					log.Println(err)
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
