package rewards

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/lib/sumsub"
	tx_watcher_svc "github.com/SatorNetwork/sator-api/svc/tx_watcher"
	tx_watcher_repository "github.com/SatorNetwork/sator-api/svc/tx_watcher/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/accounts"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestClaimRewards_ImmediateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()
	solanaMock.EXPECT().
		GetTokenAccountBalanceWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(float64(100), nil).
		AnyTimes()
	solanaMock.EXPECT().
		SerializeTxMessage(gomock.Any()).
		Return([]byte{}, nil).
		Times(1)
	solanaMock.EXPECT().
		DeserializeTxMessage(gomock.Any()).
		Return(types.Message{}, nil).
		Times(1)
	solanaMock.EXPECT().
		NewTransaction(gomock.Any()).
		Return(types.Transaction{}, nil).
		Times(1)
	solanaMock.EXPECT().
		GetLatestBlockhash(gomock.Any()).
		Return(rpc.GetLatestBlockhashValue{}, nil).
		Times(1)
	solanaMock.EXPECT().
		PrepareSendAssetsTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&lib_solana.PrepareTxResponse{}, nil).
		Times(1)
	solanaMock.EXPECT().
		SendConstructedTransaction(gomock.Any(), gomock.Any()).
		Return("", nil).
		Times(1)
	solanaMock.EXPECT().
		IsTransactionSuccessful(gomock.Any(), gomock.Any()).
		Return(true, nil).
		Times(1)

	defer app_config.RunAndWait()()

	c := client.NewClient()

	txWatcherRepository, err := tx_watcher_repository.Prepare(context.Background(), c.DB.Client())
	if err != nil {
		log.Fatalf("can't prepare tx watcher repository: %v", err)
	}
	err = txWatcherRepository.CleanTransactions(context.Background())
	require.NoError(t, err)

	user := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	err = c.DB.AuthDB().UpdateKYCStatus(context.TODO(), user.Email(), sumsub.KYCStatusApproved)
	require.NoError(t, err)

	id, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user.Email())
	require.NoError(t, err)

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(user.AccessToken())
	require.NoError(t, err)
	require.NotEqual(t, "", resp.TransactionURL)

	var txWatcherSvc *tx_watcher_svc.Service
	{
		feePayer, tokenHolder, _ := accounts.GetAccounts()

		txWatcherSvc = tx_watcher_svc.NewService(
			txWatcherRepository,
			solanaMock,
			feePayer,
			tokenHolder,
		)
	}
	err = txWatcherSvc.ResendSolanaDBTXsIfNeeded(context.Background())
	require.NoError(t, err)

	allTransactions, err := txWatcherRepository.GetAllTransactions(context.Background())
	require.NoError(t, err)
	require.Len(t, allTransactions, 1)
	require.Equal(t, allTransactions[0].Status, "successful")
}

func TestClaimRewards_SuccessAfterSomeTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()
	solanaMock.EXPECT().
		GetTokenAccountBalanceWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(float64(100), nil).
		AnyTimes()
	solanaMock.EXPECT().
		SerializeTxMessage(gomock.Any()).
		Return([]byte{}, nil).
		Times(1)
	solanaMock.EXPECT().
		DeserializeTxMessage(gomock.Any()).
		Return(types.Message{}, nil).
		Times(1)
	solanaMock.EXPECT().
		NewTransaction(gomock.Any()).
		Return(types.Transaction{}, nil).
		Times(1)
	solanaMock.EXPECT().
		GetLatestBlockhash(gomock.Any()).
		Return(rpc.GetLatestBlockhashValue{}, nil).
		Times(1)
	solanaMock.EXPECT().
		PrepareSendAssetsTx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&lib_solana.PrepareTxResponse{}, nil).
		Times(1)
	solanaMock.EXPECT().
		SendConstructedTransaction(gomock.Any(), gomock.Any()).
		Return("", nil).
		Times(1)
	{
		var cnt int
		callback := func(ctx context.Context, txhash string) (bool, error) {
			cnt++
			if cnt <= 2 {
				return false, nil
			}

			return true, nil
		}
		solanaMock.EXPECT().
			IsTransactionSuccessful(gomock.Any(), gomock.Any()).
			DoAndReturn(callback).
			Times(3)
	}
	solanaMock.EXPECT().
		NeedToRetry(gomock.Any(), gomock.Any()).
		Return(false, nil).
		Times(2)

	defer app_config.RunAndWait()()

	c := client.NewClient()

	txWatcherRepository, err := tx_watcher_repository.Prepare(context.Background(), c.DB.Client())
	if err != nil {
		log.Fatalf("can't prepare tx watcher repository: %v", err)
	}
	err = txWatcherRepository.CleanTransactions(context.Background())
	require.NoError(t, err)

	user := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	err = c.DB.AuthDB().UpdateKYCStatus(context.TODO(), user.Email(), sumsub.KYCStatusApproved)
	require.NoError(t, err)

	id, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user.Email())
	require.NoError(t, err)

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(user.AccessToken())
	require.NoError(t, err)
	require.NotEqual(t, "", resp.TransactionURL)

	var txWatcherSvc *tx_watcher_svc.Service
	{
		feePayer, tokenHolder, _ := accounts.GetAccounts()

		txWatcherSvc = tx_watcher_svc.NewService(
			txWatcherRepository,
			solanaMock,
			feePayer,
			tokenHolder,
		)
	}
	err = txWatcherSvc.ResendSolanaDBTXsIfNeeded(context.Background())
	require.NoError(t, err)

	allTransactions, err := txWatcherRepository.GetAllTransactions(context.Background())
	require.NoError(t, err)
	require.Len(t, allTransactions, 1)
	require.Equal(t, "registered", allTransactions[0].Status)

	err = txWatcherSvc.ResendSolanaDBTXsIfNeeded(context.Background())
	require.NoError(t, err)
	err = txWatcherSvc.ResendSolanaDBTXsIfNeeded(context.Background())
	require.NoError(t, err)
}
