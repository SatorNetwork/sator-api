package worker

import (
	"context"
	"database/sql"
	"log"

	"github.com/SatorNetwork/sator-api/svc/rewards/consts"

	"github.com/google/uuid"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/svc/rewards/interfaces"
	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
)

type (
	InProgressTransactionStatusWorkerJob struct {
		UserID uuid.UUID
		Amount float64
		TxHash string
	}

	InProgressTransactionStatusWorker struct {
		transactions chan InProgressTransactionStatusWorkerJob

		ctx  context.Context
		repo interfaces.RewardsRepository
		sc   solanaClient

		done chan struct{}
	}

	solanaClient interface {
		GetTransaction(ctx context.Context, txhash string) (lib_solana.GetConfirmedTransactionResponse, error)
	}
)

func NewInProgressTransactionStatusWorker(
	ctx context.Context,
	repo interfaces.RewardsRepository,
	sc solanaClient,
) *InProgressTransactionStatusWorker {
	return &InProgressTransactionStatusWorker{
		transactions: make(chan InProgressTransactionStatusWorkerJob),
		ctx:          ctx,
		repo:         repo,
		sc:           sc,
		done:         make(chan struct{}),
	}
}

func (w *InProgressTransactionStatusWorker) Start() {
LOOP:
	for {
		select {
		case transaction := <-w.transactions:
			resp, err := w.sc.GetTransaction(w.ctx, transaction.TxHash)
			if err != nil {
				log.Println(err)
				continue
			}

			if resp.Meta.Err != nil {
				err := w.repo.UpdateTransactionStatusByTxHash(w.ctx, repository.UpdateTransactionStatusByTxHashParams{
					Status: consts.TransactionStatusFailed.String(),
					TxHash: sql.NullString{
						String: transaction.TxHash,
						Valid:  true,
					},
				})
				if err != nil {
					log.Println(err)
					continue
				}
				break
			}

			err = w.repo.UpdateTransactionStatusByTxHash(w.ctx, repository.UpdateTransactionStatusByTxHashParams{
				Status: consts.TransactionStatusWithdrawn.String(),
				TxHash: sql.NullString{
					String: transaction.TxHash,
					Valid:  true,
				},
			})
			if err != nil {
				log.Println(err)
				continue
			}

			if err = w.repo.Withdraw(w.ctx, transaction.UserID); err != nil {
				log.Println(err)
				continue
			}
		case <-w.done:
			break LOOP
		}
	}
}

func (w *InProgressTransactionStatusWorker) Close() {
	close(w.done)
}

func (w *InProgressTransactionStatusWorker) AddTransaction(transaction InProgressTransactionStatusWorkerJob) {
	w.transactions <- transaction
}
