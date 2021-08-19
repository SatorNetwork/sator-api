package solana

import (
	"context"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// InitializeStakePoolInput ...
type InitializeStakePoolInput struct {
	Number uint8 // 0
	Ranks  []Rank
}

// StakeInput ...
type StakeInput struct {
	Number   uint8 // 1
	Duration int64
	Amount   uint64
}

// UnstakeInput ...
type UnstakeInput struct {
	Number uint8 // 2
}

// InitializeStakePool ...
func (c *Client) InitializeStakePool(ctx context.Context, feePayer, signer, asset, issuer types.Account) (string, error) {
	stakePool := c.NewAccount()
	systemProgram := c.PublicKeyFromString("11111111111111111111111111111111")
	sysvarRent := c.PublicKeyFromString("SysvarRent111111111111111111111111111111111")
	splToken := c.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	programID := c.PublicKeyFromString("CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u")

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	data, err := borsh.Serialize(InitializeStakePoolInput{Number: 0, Ranks: []Rank{{0, 1000}, {3600, 2000}, {7200, 3000}, {10800, 4000}}})
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, stakePool.PublicKey.Bytes()[0:32])
	stakeAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", fmt.Errorf("could not get stake authority: %w", err)
	}

	tokenAccountStakePool := common.CreateWithSeed(stakeAuthority, "ViewerStakePool::token_account", splToken)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: programID,
				Accounts: []types.AccountMeta{
					{PubKey: systemProgram, IsSigner: false, IsWritable: false},        // system - link incl
					{PubKey: sysvarRent, IsSigner: false, IsWritable: false},           // system - link incl
					{PubKey: splToken, IsSigner: false, IsWritable: false},             // SPL "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},    // fee_payer
					{PubKey: issuer.PublicKey, IsSigner: true, IsWritable: false},      // issuer
					{PubKey: stakePool.PublicKey, IsSigner: true, IsWritable: true},    // new acc for stake pool
					{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},       // to generate from rules #1
					{PubKey: tokenAccountStakePool, IsSigner: false, IsWritable: true}, // to generate from rules #2
					{PubKey: asset.PublicKey, IsSigner: false, IsWritable: true},       // MINT
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer, signer},
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

// Stake ...
func (c *Client) Stake(ctx context.Context, feePayer, signer types.Account, pool, stakeAuthority, userWallet, tokenAccountStakeTarget common.PublicKey) (string, error) {
	sysvarClock := c.PublicKeyFromString("SysvarC1ock11111111111111111111111111111111")
	sysvarRent := c.PublicKeyFromString("SysvarRent111111111111111111111111111111111")
	systemProgram := c.PublicKeyFromString("11111111111111111111111111111111")
	splToken := c.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	programID := c.PublicKeyFromString("CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u")
	tokenAccountSource := c.NewAccount()

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	data, err := borsh.Serialize(StakeInput{
		Number:   1,
		Duration: 3600,
		Amount:   1000,
	})
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	seed := userWallet.Bytes()
	seedString := base58.Encode(seed[0:20])
	stakeAccount := common.CreateWithSeed(stakeAuthority, seedString, programID)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: c.PublicKeyFromString("CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u"),
				Accounts: []types.AccountMeta{
					{PubKey: systemProgram, IsSigner: false, IsWritable: false},
					{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
					{PubKey: sysvarClock, IsSigner: false, IsWritable: false},
					{PubKey: splToken, IsSigner: false, IsWritable: false},                    // spl
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},           // feepayer // fee payer, lock owner
					{PubKey: pool, IsSigner: false, IsWritable: true},                         // from init
					{PubKey: stakeAuthority, IsSigner: false, IsWritable: false},              // предварительно юзер на кого делаем стейк
					{PubKey: tokenAccountSource.PublicKey, IsSigner: false, IsWritable: true}, // spl acc with Wallet PK
					{PubKey: tokenAccountStakeTarget, IsSigner: false, IsWritable: true},      // generated #3
					{PubKey: stakeAccount, IsSigner: false, IsWritable: true},                 // gen from wallet PK, time when locked.
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer, signer},
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

// Unstake ...
func (c *Client) Unstake(ctx context.Context, feePayer, signer, asset, issuer, pool, userWallet, tokenAccount, stake, stakeAuthority types.Account) (string, error) {
	sysvarClock := c.PublicKeyFromString("SysvarC1ock11111111111111111111111111111111")
	splToken := c.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

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
				ProgramID: c.PublicKeyFromString("CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u"),
				Accounts: []types.AccountMeta{
					{PubKey: sysvarClock, IsSigner: false, IsWritable: false},              // sysvar clock const
					{PubKey: splToken, IsSigner: false, IsWritable: false},                 // spl token const
					{PubKey: pool.PublicKey, IsSigner: false, IsWritable: false},           // pool from init
					{PubKey: stakeAuthority.PublicKey, IsSigner: false, IsWritable: false}, // from init
					{PubKey: userWallet.PublicKey, IsSigner: false, IsWritable: true},      // wallet
					{PubKey: tokenAccount.PublicKey, IsSigner: false, IsWritable: true},    // token account from init
					{PubKey: stake.PublicKey, IsSigner: false, IsWritable: true},           // gen in stake
					{PubKey: issuer.PublicKey, IsSigner: true, IsWritable: false},          // issuer
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer, signer},
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
