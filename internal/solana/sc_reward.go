package solana

import (
	"context"
	"fmt"
	"strconv"

	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// InitializeShow generates and calls instruction that initializes show.
func (c *Client) InitializeShow(ctx context.Context, feePayer types.Account, asset common.PublicKey) (txHast string, show types.Account, err error) {
	show = types.NewAccount()
	systemProgram := c.PublicKeyFromString(c.config.SystemProgram)
	sysvarRent := c.PublicKeyFromString(c.config.SysvarRent)
	splToken := c.PublicKeyFromString(c.config.SplToken)
	programID := c.PublicKeyFromString(c.config.RewardProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not get recent block hash: %w", err)
	}

	input := InitializeShowInput{RewardLockTime: 100}
	data, err := borsh.Serialize(input)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, show.PublicKey.Bytes())
	showAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not get show authority: %w", err)
	}

	tokenAccountStakePool := common.CreateWithSeed(showAuthority, "Show::token_account", splToken)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: programID,
				Accounts: []types.AccountMeta{
					{PubKey: systemProgram, IsSigner: false, IsWritable: false},
					{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
					{PubKey: splToken, IsSigner: false, IsWritable: false},
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
					{PubKey: show.PublicKey, IsSigner: true, IsWritable: true},
					{PubKey: showAuthority, IsSigner: false, IsWritable: false},
					{PubKey: tokenAccountStakePool, IsSigner: false, IsWritable: true},
					{PubKey: asset, IsSigner: false, IsWritable: false},
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer, show},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", types.Account{}, fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, show, nil
}

// InitializeViewer generates and calls instruction that initializes viewer.
func (c *Client) InitializeViewer(ctx context.Context, feePayer, show types.Account, wallet common.PublicKey) (txHast string, err error) {
	systemProgram := c.PublicKeyFromString(c.config.SystemProgram)
	sysvarRent := c.PublicKeyFromString(c.config.SysvarRent)
	programID := c.PublicKeyFromString(c.config.RewardProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	input := InitializeViewer{UserPubKey: wallet}
	data, err := borsh.Serialize(input)
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, show.PublicKey.Bytes())
	showAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", fmt.Errorf("could not get show authority: %w", err)
	}

	viewer := common.CreateWithSeed(showAuthority, wallet.ToBase58(), programID)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: programID,
				Accounts: []types.AccountMeta{
					{PubKey: systemProgram, IsSigner: false, IsWritable: false},
					{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
					{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
					{PubKey: show.PublicKey, IsSigner: true, IsWritable: false},
					{PubKey: showAuthority, IsSigner: false, IsWritable: false},
					{PubKey: viewer, IsSigner: false, IsWritable: true},
				},
				Data: data,
			},
		},
		Signers:         []types.Account{feePayer, show},
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

// InitializeQuiz generates and calls instruction that initializes quiz.
func (c *Client) InitializeQuiz(ctx context.Context, feePayer, show types.Account, wallet common.PublicKey, quizIndex uint64, winners []WinnerInput, amount uint64) (txHast string, err error) {
	systemProgram := c.PublicKeyFromString(c.config.SystemProgram)
	sysvarRent := c.PublicKeyFromString(c.config.SysvarRent)
	sysvarClock := c.PublicKeyFromString(c.config.SysvarClock)
	programID := c.PublicKeyFromString(c.config.RewardProgramID)

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	input := InitializeQuizInput{Winners: winners, TokenAmount: amount}
	data, err := borsh.Serialize(input)
	if err != nil {
		return "", fmt.Errorf("could not serialize data with borsh: %w", err)
	}

	var seeds [][]byte
	seeds = append(seeds, show.PublicKey.Bytes())
	showAuthority, _, err := common.FindProgramAddress(seeds, programID)
	if err != nil {
		return "", fmt.Errorf("could not get show authority: %w", err)
	}

	quizNumber := strconv.FormatUint(quizIndex, 10)
	quizPubkey := common.CreateWithSeed(showAuthority, wallet.ToBase58()+quizNumber, programID)

	accounts := []types.AccountMeta{
		{PubKey: systemProgram, IsSigner: false, IsWritable: false},
		{PubKey: sysvarRent, IsSigner: false, IsWritable: false},
		{PubKey: sysvarClock, IsSigner: false, IsWritable: false},
		{PubKey: feePayer.PublicKey, IsSigner: true, IsWritable: false},
		{PubKey: show.PublicKey, IsSigner: true, IsWritable: false},
		{PubKey: showAuthority, IsSigner: false, IsWritable: false},
		{PubKey: quizPubkey, IsSigner: false, IsWritable: true},
	}

	for _, v := range winners {
		accounts = append(accounts, types.AccountMeta{
			PubKey:     v.UserPubKey,
			IsSigner:   false,
			IsWritable: false,
		})
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			{
				ProgramID: programID,
				Accounts:  accounts,
				Data:      data,
			},
		},
		Signers:         []types.Account{feePayer, show},
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
