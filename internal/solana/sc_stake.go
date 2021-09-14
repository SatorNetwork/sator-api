package solana

import (
	"context"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// InitializeStakePool generates and calls instruction that initializes stake pool
func (c *Client) InitializeStakePool(ctx context.Context, feePayer types.Account, asset common.PublicKey) (txHast string, stakePool types.Account, stakeAuthority, tokenAccountStakePool common.PublicKey, err error) {
	stakePool = types.NewAccount()
	systemProgram := c.PublicKeyFromString(SystemProgram)
	sysvarRent := c.PublicKeyFromString(SysvarRent)
	splToken := c.PublicKeyFromString(SplToken)
	programID := c.PublicKeyFromString(ProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", types.Account{}, common.PublicKey{}, common.PublicKey{}, fmt.Errorf("could not get recent block hash: %w", err)
	}

	input := InitializeStakePoolInput{Number: 0, Ranks: [4]Rank{
		{0, 100},
		{1800, 200},
		{3600, 300},
		{7200, 500}}}
	data, err := borsh.Serialize(input)
	if err != nil {
		return "", types.Account{}, common.PublicKey{}, common.PublicKey{}, fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, stakePool.PublicKey.Bytes()[0:32])
	stakeAuthority, _, err = common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", types.Account{}, common.PublicKey{}, common.PublicKey{}, fmt.Errorf("could not get stake authority: %w", err)
	}

	tokenAccountStakePool = common.CreateWithSeed(stakeAuthority, "ViewerStakePool::token_account", splToken)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: programID,
				Accounts: []types.AccountMeta{
					{PubKey: systemProgram, IsSigner: false, IsWritable: false},
					{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
					{PubKey: splToken, IsSigner: false, IsWritable: false},
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
					{PubKey: stakePool.PublicKey, IsSigner: true, IsWritable: true},
					{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},
					{PubKey: tokenAccountStakePool, IsSigner: false, IsWritable: true},
					{PubKey: asset, IsSigner: false, IsWritable: false},
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer, stakePool},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", types.Account{}, common.PublicKey{}, common.PublicKey{}, fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", types.Account{}, common.PublicKey{}, common.PublicKey{}, fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, stakePool, stakeAuthority, tokenAccountStakePool, nil
}

// Stake stakes given amount for given period.
func (c *Client) Stake(ctx context.Context, feePayer types.Account, pool, stakeAuthority, userWallet, tokenAccountStakeTarget common.PublicKey, duration int64, amount uint64) (string, common.PublicKey, error) {
	sysvarClock := c.PublicKeyFromString(SysvarClock)
	sysvarRent := c.PublicKeyFromString(SysvarRent)
	systemProgram := c.PublicKeyFromString(SystemProgram)
	splToken := c.PublicKeyFromString(SplToken)
	programID := c.PublicKeyFromString(ProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", common.PublicKey{}, fmt.Errorf("could not get recent block hash: %w", err)
	}

	data, err := borsh.Serialize(StakeInput{
		Number:   1,
		Duration: duration,
		Amount:   amount,
	})
	if err != nil {
		return "", common.PublicKey{}, fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	seed := userWallet.Bytes()
	seedString := base58.Encode(seed[0:20])
	stakeAccount := common.CreateWithSeed(stakeAuthority, seedString, programID)
	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
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
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
					{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},
					{PubKey: userWallet, IsSigner: false, IsWritable: true},
					{PubKey: tokenAccountStakeTarget, IsSigner: false, IsWritable: true},
					{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", common.PublicKey{}, fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", common.PublicKey{}, fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, stakeAccount, nil
}

// Unstake unstake.
func (c *Client) Unstake(ctx context.Context, feePayer types.Account, stakePool, userWallet, tokenAccount, stakeAccount, stakeAuthority common.PublicKey) (string, error) {
	sysvarClock := c.PublicKeyFromString(SysvarClock)
	splToken := c.PublicKeyFromString(SplToken)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	data, err := borsh.Serialize(UnstakeInput{Number: 2})
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: c.PublicKeyFromString(ProgramID),
				Accounts: []types.AccountMeta{
					{PubKey: sysvarClock, IsSigner: false, IsWritable: false},
					{PubKey: splToken, IsSigner: false, IsWritable: false},
					{PubKey: stakePool, IsSigner: false, IsWritable: false},
					{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},
					{PubKey: userWallet, IsSigner: false, IsWritable: true},
					{PubKey: tokenAccount, IsSigner: false, IsWritable: true},
					{PubKey: stakeAccount, IsSigner: false, IsWritable: true},
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, nil
}
