package wallet

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/google/uuid"
)

type walletRepoMock struct {
	userStakeAmount float64
	multiplier      int32

	// errors
	GetStakeByUserIDErr      error
	GetStakeLevelByAmountErr error
}

func (r *walletRepoMock) CreateWallet(ctx context.Context, arg repository.CreateWalletParams) (repository.Wallet, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetWalletBySolanaAccountID(ctx context.Context, solanaAccountID uuid.UUID) (repository.Wallet, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetWalletByUserIDAndType(ctx context.Context, arg repository.GetWalletByUserIDAndTypeParams) (repository.Wallet, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) AddSolanaAccount(ctx context.Context, arg repository.AddSolanaAccountParams) (repository.SolanaAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetSolanaAccountByID(ctx context.Context, id uuid.UUID) (repository.SolanaAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetSolanaAccountByType(ctx context.Context, accountType string) (repository.SolanaAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetSolanaAccountTypeByPublicKey(ctx context.Context, publicKey string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetSolanaAccountByUserIDAndType(ctx context.Context, arg repository.GetSolanaAccountByUserIDAndTypeParams) (repository.SolanaAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) AddEthereumAccount(ctx context.Context, arg repository.AddEthereumAccountParams) (repository.EthereumAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetEthereumAccountByID(ctx context.Context, id uuid.UUID) (repository.EthereumAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetEthereumAccountByUserIDAndType(ctx context.Context, arg repository.GetEthereumAccountByUserIDAndTypeParams) (repository.EthereumAccount, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) AddStake(ctx context.Context, arg repository.AddStakeParams) (repository.Stake, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) DeleteStakeByUserID(ctx context.Context, userID uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetStakeByUserID(ctx context.Context, userID uuid.UUID) (repository.Stake, error) {
	if r.GetStakeByUserIDErr != nil {
		return repository.Stake{}, r.GetStakeByUserIDErr
	}
	return repository.Stake{
		ID:          uuid.New(),
		UserID:      userID,
		StakeAmount: r.userStakeAmount,
	}, nil
}

func (r *walletRepoMock) GetTotalStake(ctx context.Context) (float64, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) UpdateStake(ctx context.Context, arg repository.UpdateStakeParams) error {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetAllStakeLevels(ctx context.Context) ([]repository.StakeLevel, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetAllEnabledStakeLevels(ctx context.Context) ([]repository.StakeLevel, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetMinimalStakeLevel(ctx context.Context) (repository.StakeLevel, error) {
	panic("not implemented") // TODO: Implement
}

func (r *walletRepoMock) GetStakeLevelByAmount(ctx context.Context, amount float64) (repository.GetStakeLevelByAmountRow, error) {
	if r.GetStakeLevelByAmountErr != nil {
		return repository.GetStakeLevelByAmountRow{}, r.GetStakeLevelByAmountErr
	}

	return repository.GetStakeLevelByAmountRow{
		ID: uuid.New(),
		Multiplier: sql.NullInt32{
			Int32: r.multiplier,
			Valid: true,
		},
	}, nil
}

func TestService_GetMultiplier(t *testing.T) {
	type fields struct {
		wr                          walletRepository
		sc                          solanaClient
		ec                          ethereumClient
		satorAssetName              string
		solanaAssetName             string
		satorAssetSolanaAddr        string
		feePayerSolanaAddr          string
		feePayerSolanaPrivateKey    []byte
		stakePoolSolanaPublicKey    string
		tokenHolderSolanaAddr       string
		tokenHolderSolanaPrivateKey []byte
		walletDetailsURL            string
		walletTransactionsURL       string
		rewardsWalletDetailsURL     string
		rewardsTransactionsURL      string
		minAmountToTransfer         float64
	}
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int32
		wantErr bool
	}{
		{
			"no levels found",
			fields{wr: &walletRepoMock{GetStakeLevelByAmountErr: sql.ErrNoRows}},
			args{context.TODO(), uuid.New()},
			0,
			false,
		},
		{
			"user does not hold any tokens",
			fields{wr: &walletRepoMock{GetStakeByUserIDErr: sql.ErrNoRows}},
			args{context.TODO(), uuid.New()},
			0,
			false,
		},
		{
			"user holds 100 tokens",
			fields{wr: &walletRepoMock{userStakeAmount: 100, multiplier: 1}},
			args{context.TODO(), uuid.New()},
			1,
			false,
		},
		{
			"unexpected db error while getting user hold amount",
			fields{wr: &walletRepoMock{GetStakeByUserIDErr: errors.New("unexpected error")}},
			args{context.TODO(), uuid.New()},
			0,
			true,
		},
		{
			"unexpected db error while getting stake level",
			fields{wr: &walletRepoMock{GetStakeLevelByAmountErr: errors.New("unexpected error")}},
			args{context.TODO(), uuid.New()},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				wr:                          tt.fields.wr,
				sc:                          tt.fields.sc,
				ec:                          tt.fields.ec,
				satorAssetName:              tt.fields.satorAssetName,
				solanaAssetName:             tt.fields.solanaAssetName,
				satorAssetSolanaAddr:        tt.fields.satorAssetSolanaAddr,
				feePayerSolanaAddr:          tt.fields.feePayerSolanaAddr,
				feePayerSolanaPrivateKey:    tt.fields.feePayerSolanaPrivateKey,
				stakePoolSolanaPublicKey:    tt.fields.stakePoolSolanaPublicKey,
				tokenHolderSolanaAddr:       tt.fields.tokenHolderSolanaAddr,
				tokenHolderSolanaPrivateKey: tt.fields.tokenHolderSolanaPrivateKey,
				walletDetailsURL:            tt.fields.walletDetailsURL,
				walletTransactionsURL:       tt.fields.walletTransactionsURL,
				rewardsWalletDetailsURL:     tt.fields.rewardsWalletDetailsURL,
				rewardsTransactionsURL:      tt.fields.rewardsTransactionsURL,
				minAmountToTransfer:         tt.fields.minAmountToTransfer,
			}
			got, err := s.GetMultiplier(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetMultiplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetMultiplier() = %v, want %v", got, tt.want)
			}
		})
	}
}
