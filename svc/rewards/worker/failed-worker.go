package worker

import (
	"context"
	"database/sql"
	"log"

	"github.com/SatorNetwork/sator-api/svc/rewards/interfaces"
	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
)

type FailedTransactionStatusWorker struct {
	ipWorker *InProgressTransactionStatusWorker

	ctx  context.Context
	repo interfaces.RewardsRepository
	ws   interfaces.WalletService
}

func NewFailedTransactionStatusWorker(
	ctx context.Context,
	repo interfaces.RewardsRepository,
	ipWorker *InProgressTransactionStatusWorker,
) *FailedTransactionStatusWorker {
	return &FailedTransactionStatusWorker{
		ipWorker: ipWorker,
		ctx:      ctx,
		repo:     repo,
	}
}

func (w *FailedTransactionStatusWorker) Start() {
	transactions, err := w.repo.GetFailedTransactions(w.ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for i := range transactions {
		userID := transactions[i].UserID
		amount := transactions[i].Amount
		txHash, err := w.ws.WithdrawRewards(w.ctx, userID, transactions[i].Amount)
		if err != nil {
			log.Println(err)
			continue
		}

		if err := w.repo.AddTransaction(w.ctx, repository.AddTransactionParams{
			UserID:          userID,
			Amount:          amount,
			TransactionType: 2, // TransactionTypeWithdraw
			TxHash: sql.NullString{
				String: txHash,
				Valid:  true,
			},
			Status: 2, // TransactionStatusInProgress
		}); err != nil {
			log.Println(err)
			continue
		}

		w.ipWorker.AddTransaction(InProgressTransactionStatusWorkerJob{
			UserID: userID,
			Amount: amount,
			TxHash: txHash,
		})
	}
}
