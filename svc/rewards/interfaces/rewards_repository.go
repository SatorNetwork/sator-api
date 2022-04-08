package interfaces

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
)

type RewardsRepository interface {
	AddTransaction(ctx context.Context, arg repository.AddTransactionParams) error
	Withdraw(ctx context.Context, uid uuid.UUID) error
	GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
	GetTransactionsByUserIDPaginated(ctx context.Context, arg repository.GetTransactionsByUserIDPaginatedParams) ([]repository.Reward, error)
	// GetAmountAvailableToWithdraw(ctx context.Context, arg repository.GetAmountAvailableToWithdrawParams) (float64, error)
	GetScannedQRCodeByUserID(ctx context.Context, arg repository.GetScannedQRCodeByUserIDParams) (repository.Reward, error)
	RequestTransactionsByUserID(ctx context.Context, userID uuid.UUID) error
	SetInProgressTransaction(ctx context.Context, arg repository.SetInProgressTransactionParams) error
	UpdateTransactionStatusByTxHash(ctx context.Context, arg repository.UpdateTransactionStatusByTxHashParams) error
	GetFailedTransactions(ctx context.Context) ([]repository.Reward, error)
}
