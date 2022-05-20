//go:build !mock_solana

package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	pkg_errors "github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/types"

	"github.com/SatorNetwork/sator-api/lib/fee_accumulator"
	"github.com/SatorNetwork/sator-api/lib/solana"
)

func (c *Client) SendAssetsWithAutoDerive(
	ctx context.Context,
	assetAddr string,
	feePayer types.Account,
	source types.Account,
	recipientAddr string,
	amount float64,
	cfg *solana.SendAssetsConfig,
) (string, error) {
	resp, err := c.PrepareSendAssetsTx(ctx, assetAddr, feePayer, source, recipientAddr, amount, cfg)
	if err != nil {
		return "", pkg_errors.Wrap(err, "can't prepare send assets tx")
	}

	txHash, err := c.solana.SendTransaction(ctx, resp.Tx)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}

func (c *Client) PrepareSendAssetsTx(
	ctx context.Context,
	assetAddr string,
	feePayer types.Account,
	source types.Account,
	recipientAddr string,
	amount float64,
	cfg *solana.SendAssetsConfig,
) (*solana.PrepareTxResponse, error) {
	feeAccumulator, err := fee_accumulator.New(c.exchangeRatesClient)
	if err != nil {
		return nil, pkg_errors.Wrap(err, "can't create new fee accumulator")
	}

	if !(cfg.PercentToCharge >= 0 && cfg.PercentToCharge <= 100) {
		return nil, fmt.Errorf("percent to charge fees invalid: %v", cfg.PercentToCharge)
	}

	feeAccumulator.AddSAO(amount * cfg.PercentToCharge / 100)
	asset := common.PublicKeyFromString(assetAddr)

	sourceAta, _, err := common.FindAssociatedTokenAddress(source.PublicKey, asset)
	if err != nil {
		return nil, pkg_errors.Wrap(err, "can't find associated token address for source account")
	}

	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, err := c.deriveATAPublicKey(ctx, recipientPublicKey, asset)
	if err != nil {
		if !errors.Is(err, ErrATANotCreated) {
			return nil, err
		}
		// Add instruction to create token account
		// instructions = append(instructions,
		// 	assotokenprog.CreateAssociatedTokenAccount(
		// 		feePayer.PublicKey,
		// 		recipientPublicKey,
		// 		common.PublicKeyFromString(assetAddr),
		// 	))

		_, err := c.CreateAccountWithATA(ctx, assetAddr, recipientPublicKey.ToBase58(), feePayer)
		if err != nil {
			// return "", fmt.Errorf("CreateAccountWithATA: %w", err)
			log.Printf("CreateAccountWithATA: %v", err)
		}
		recipientAta, err = c.deriveATAPublicKey(ctx, recipientPublicKey, asset)
		if err != nil {
			return nil, err
		}
	}

	message, err := c.prepareSendAssetsMessage(
		ctx,
		feePayer,
		sourceAta,
		recipientAta,
		asset,
		source.PublicKey,
		amount-feeAccumulator.GetFeeInSAO(),
		feeAccumulator.GetFeeInSAO(),
	)
	if err != nil {
		return nil, pkg_errors.Wrap(err, "can't prepare send assets message (before adding blockchain fee)")
	}

	var solanaTxFee uint64
	if false {
		//if cfg.ChargeSolanaFeeFromSender {
		solanaTxFee, err = c.GetFeeForMessage(ctx, message, cfg.AllowFallbackToDefaultFee, cfg.DefaultFee)
		if err != nil {
			return nil, pkg_errors.Wrap(err, "can't get fee for message")
		}
		feeAccumulator.AddSOL(float64(solanaTxFee) / fee_accumulator.SolMltpl)

		if amount <= feeAccumulator.GetFeeInSAO() {
			return nil, pkg_errors.Errorf("amount <= fee, amount: %v, fee: %v", amount, feeAccumulator.GetFeeInSAO())
		}

		message, err = c.prepareSendAssetsMessage(
			ctx,
			feePayer,
			sourceAta,
			recipientAta,
			asset,
			source.PublicKey,
			amount-feeAccumulator.GetFeeInSAO(),
			feeAccumulator.GetFeeInSAO(),
		)
		if err != nil {
			return nil, pkg_errors.Wrap(err, "can't prepare send assets message (after adding blockchain fee)")
		}
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{feePayer, source},
	})
	if err != nil {
		return nil, fmt.Errorf("could not create new raw transaction: %w", err)
	}

	return &solana.PrepareTxResponse{
		Tx:                      tx,
		FeeInSAO:                feeAccumulator.GetFeeInSAO(),
		BlockchainFeeInSOLMltpl: solanaTxFee,
	}, nil
}

func (c *Client) prepareSendAssetsMessage(
	ctx context.Context,
	feePayer types.Account,
	sourceAta common.PublicKey,
	recipientAta common.PublicKey,
	asset common.PublicKey,
	sourcePublicKey common.PublicKey,
	amount float64,
	satorFee float64,
) (types.Message, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	satorFeeToSend := uint64(satorFee * float64(c.mltpl))

	instructions := make([]types.Instruction, 0, 2)
	instructions = append(instructions, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
		From:     sourceAta,
		To:       recipientAta,
		Mint:     asset,
		Auth:     sourcePublicKey,
		Signers:  []common.PublicKey{},
		Amount:   amountToSend,
		Decimals: c.decimals,
	}))

	if satorFee > 0 {
		if c.config.FeeAccumulatorAddress == "" {
			return types.Message{}, pkg_errors.Errorf("Fee accumulator address is empty")
		}

		feeAccumulatorPublicKey := common.PublicKeyFromString(c.config.FeeAccumulatorAddress)
		feeAccumulatorAta, err := c.deriveATAPublicKey(ctx, feeAccumulatorPublicKey, asset)
		if err != nil {
			return types.Message{}, pkg_errors.Wrapf(err, "can't derive ata public key for fee accumulator, addr: %v", c.config.FeeAccumulatorAddress)
		}

		instructions = append(instructions, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
			From:     sourceAta,
			To:       feeAccumulatorAta,
			Mint:     asset,
			Auth:     sourcePublicKey,
			Signers:  []common.PublicKey{},
			Amount:   satorFeeToSend,
			Decimals: c.decimals,
		}))
	}

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return types.Message{}, fmt.Errorf("could not get recent block hash: %w", err)
	}

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        feePayer.PublicKey,
		Instructions:    instructions,
		RecentBlockhash: res.Blockhash,
	})

	return message, nil
}
