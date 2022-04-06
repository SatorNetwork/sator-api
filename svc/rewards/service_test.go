package rewards

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
	wallet "github.com/SatorNetwork/sator-api/svc/wallet/client"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestService_ClaimRewards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUID, err := uuid.NewUUID()
	require.NoError(t, err)
	txHash := "test"

	rewardsRepositoryMock := NewMockRewardsRepository(ctrl)
	mock.RegisterMockObject("RewardsRepository", rewardsRepositoryMock)
	rewardsRepositoryMock.EXPECT().
		GetTotalAmount(gomock.Any(), gomock.Any()).
		Return(100.0, nil).
		AnyTimes()
	rewardsRepositoryMock.EXPECT().
		AddTransaction(gomock.Any(), repository.AddTransactionParams{
		UserID:          userUID,
		TransactionType: TransactionTypeWithdraw,
		Amount:          100.0,
		TxHash:          sql.NullString{
			String: txHash,
			Valid:  true,
		},
	}).
		Return(nil).
		AnyTimes()
	rewardsRepositoryMock.EXPECT().
		Withdraw(gomock.Any(), userUID).
		Return(nil).
		AnyTimes()

	walletSvcClientMock := wallet.NewMockService(ctrl)
	mock.RegisterMockObject("WalletClient-1", walletSvcClientMock)
	walletSvcClientMock.EXPECT().
		WithdrawRewards(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(txHash, nil).
		AnyTimes()

	svc := NewService(rewardsRepositoryMock, walletSvcClientMock, nil)
	_, err = svc.ClaimRewards(nil, userUID)
	require.NoError(t, err)
}
