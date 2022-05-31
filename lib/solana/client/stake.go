//go:build !mock_solana

package client

import (
	"context"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// InitializeStakePool generates and calls instruction that initializes stake pool.
func (c *Client) InitializeStakePool(ctx context.Context, feePayer, issuer types.Account, asset common.PublicKey) (txHast string, stakePool types.Account, err error) {
	stakePool = types.NewAccount()
	systemProgram := c.PublicKeyFromString(c.config.SystemProgram)
	sysvarRent := c.PublicKeyFromString(c.config.SysvarRent)
	splToken := c.PublicKeyFromString(c.config.SplToken)
	programID := c.PublicKeyFromString(c.config.StakeProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not get recent block hash: %w", err)
	}

	input := InitializeStakePoolInput{Number: 0, Ranks: [4]Rank{
		{0, 100},
		{1800, 200},
		{3600, 300},
		{7200, 500}}}
	data, err := borsh.Serialize(input)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, stakePool.PublicKey.Bytes()[0:32])
	stakeAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not get stake authority: %w", err)
	}

	tokenAccountStakePool := common.CreateWithSeed(stakeAuthority, "ViewerStakePool::token_account", splToken)

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer: feePayer.PublicKey,
			Instructions: []types.Instruction{
				{
					ProgramID: programID,
					Accounts: []types.AccountMeta{
						{PubKey: systemProgram, IsSigner: false, IsWritable: false},
						{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
						{PubKey: splToken, IsSigner: false, IsWritable: false},
						{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
						{PubKey: issuer.PublicKey, IsSigner: true, IsWritable: false},
						{PubKey: stakePool.PublicKey, IsSigner: true, IsWritable: true},
						{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},
						{PubKey: tokenAccountStakePool, IsSigner: false, IsWritable: true},
						{PubKey: asset, IsSigner: false, IsWritable: false},
					},
					Data: data,
				},
			},
			RecentBlockhash: res.Blockhash,
		}),
		Signers: []types.Account{feePayer, issuer, stakePool},
	})
	if err != nil {
		return "", types.Account{}, err
	}

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, stakePool, nil
}

// Stake ...
func (c *Client) Stake(ctx context.Context, feePayer, userWallet types.Account, pool, asset common.PublicKey, duration int64, amount float64) (string, error) {
	sysvarClock := c.PublicKeyFromString(c.config.SysvarClock)
	sysvarRent := c.PublicKeyFromString(c.config.SysvarRent)
	systemProgram := c.PublicKeyFromString(c.config.SystemProgram)
	splToken := c.PublicKeyFromString(c.config.SplToken)
	programID := c.PublicKeyFromString(c.config.StakeProgramID)

	amountUint := uint64(amount * float64(c.mltpl))

	data, err := borsh.Serialize(StakeInput{
		Number:   1,
		Duration: duration,
		Amount:   amountUint,
	})
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, pool.Bytes()[0:32])
	stakeAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", fmt.Errorf("could not get stake authority: %w", err)
	}

	tokenAccountStakeTarget := common.CreateWithSeed(stakeAuthority, "ViewerStakePool::token_account", splToken)
	seed := userWallet.PublicKey.Bytes()
	seedString := base58.Encode(seed[0:20])
	stakeAccount := common.CreateWithSeed(stakeAuthority, seedString, programID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	ataPub, err := c.deriveATAPublicKey(ctx, userWallet.PublicKey, asset)
	if err != nil {
		return "", fmt.Errorf("could not derive ATA pub key: %w", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer: feePayer.PublicKey,
			Instructions: []types.Instruction{
				{
					ProgramID: programID,
					Accounts: []types.AccountMeta{
						{PubKey: systemProgram, IsSigner: false, IsWritable: false},
						{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
						{PubKey: sysvarClock, IsSigner: false, IsWritable: false},
						{PubKey: splToken, IsSigner: false, IsWritable: false},
						{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
						{PubKey: pool, IsSigner: false, IsWritable: false},
						{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},
						{PubKey: ataPub, IsSigner: false, IsWritable: true},
						{PubKey: tokenAccountStakeTarget, IsSigner: false, IsWritable: true},
						{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
						{PubKey: userWallet.PublicKey, IsSigner: true, IsWritable: false},
					},
					Data: data,
				},
			},
			RecentBlockhash: res.Blockhash,
		}),
		Signers: []types.Account{feePayer, userWallet},
	})
	if err != nil {
		return "", fmt.Errorf("could not create new transaction: %w", err)
	}

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, nil
}

// Unstake ...
func (c *Client) Unstake(ctx context.Context, feePayer, userWallet types.Account, stakePool, asset common.PublicKey) (string, error) {
	sysvarClock := c.PublicKeyFromString(c.config.SysvarClock)
	splToken := c.PublicKeyFromString(c.config.SplToken)
	programID := c.PublicKeyFromString(c.config.StakeProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	data, err := borsh.Serialize(UnstakeInput{Number: 2})
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, stakePool.Bytes()[0:32])
	stakeAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", fmt.Errorf("could not get stake authority: %w", err)
	}

	tokenAccountStakeTarget := common.CreateWithSeed(stakeAuthority, "ViewerStakePool::token_account", splToken)

	seed := userWallet.PublicKey.Bytes()
	seedString := base58.Encode(seed[0:20])
	stakeAccount := common.CreateWithSeed(stakeAuthority, seedString, programID)

	ataPub, err := c.deriveATAPublicKey(ctx, userWallet.PublicKey, asset)
	if err != nil {
		return "", fmt.Errorf("could not derive ATA pub key: %w", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer: feePayer.PublicKey,
			Instructions: []types.Instruction{
				{
					ProgramID: c.PublicKeyFromString(c.config.StakeProgramID),
					Accounts: []types.AccountMeta{
						{PubKey: sysvarClock, IsSigner: false, IsWritable: false},
						{PubKey: splToken, IsSigner: false, IsWritable: false},
						{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: true},
						{PubKey: stakePool, IsSigner: false, IsWritable: false},
						{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},
						{PubKey: ataPub, IsSigner: false, IsWritable: true},
						{PubKey: tokenAccountStakeTarget, IsSigner: false, IsWritable: true},
						{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
						{PubKey: userWallet.PublicKey, IsSigner: true, IsWritable: false},
					},
					Data: data,
				},
			},
			RecentBlockhash: res.Blockhash,
		}),
		Signers: []types.Account{feePayer, userWallet},
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, nil
}
