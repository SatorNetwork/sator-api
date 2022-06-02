//go:build !mock_solana

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
)

// GetConfirmedTransaction returns extended transaction details
func (c *Client) GetConfirmedTransaction(ctx context.Context, txhash string) (lib_solana.GetConfirmedTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result lib_solana.GetConfirmedTransactionResponse `json:"result"`
	}{}
	err := c.request(ctx, "getConfirmedTransaction", []interface{}{txhash, "json"}, &res)
	if err != nil {
		return lib_solana.GetConfirmedTransactionResponse{}, err
	}
	return res.Result, nil
}

// GetConfirmedTransactionForAccount returns transactions list for given account
func (c *Client) GetConfirmedTransactionForAccount(ctx context.Context, assetAddr, rootPubKey, txhash string) (lib_solana.ConfirmedTransactionResponse, error) {
	tx, err := c.solana.GetTransaction(ctx, txhash)
	if err != nil {
		return lib_solana.ConfirmedTransactionResponse{}, err
	}

	txInJson, err := json.Marshal(tx)
	if err != nil {
		return lib_solana.ConfirmedTransactionResponse{}, err
	}
	fmt.Printf("txInJson: %s\n", txInJson)

	if err := checkIfTxIsValid(tx); err != nil {
		err := errors.Wrap(err, "tx is invalid")
		fmt.Println(err)
		return lib_solana.ConfirmedTransactionResponse{}, err
	}

	amount, err := getTransactionAmountForAccountIdx(tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances, assetAddr, rootPubKey)
	if err != nil {
		err := errors.Wrap(err, "can't get transaction amount for account idx")
		fmt.Println(err)
		return lib_solana.ConfirmedTransactionResponse{}, err
	}
	fmt.Printf("amount: %v\n", amount)

	var blockTime int64
	if tx.BlockTime != nil {
		blockTime = *tx.BlockTime
	}

	tr := lib_solana.ConfirmedTransactionResponse{
		TxHash:        txhash,
		Amount:        amount,
		AmountString:  fmt.Sprintf("%f", amount),
		CreatedAtUnix: blockTime,
		CreatedAt:     time.Unix(blockTime, 0),
	}

	return tr, nil
}

func checkIfTxIsValid(tx *client.GetTransactionResponse) error {
	if tx.Meta == nil {
		return fmt.Errorf("tx.Meta should not be nil")
	}
	if tx.Meta.PreTokenBalances == nil {
		return fmt.Errorf("tx.Meta.PreTokenBalances should not be nil")
	}
	if tx.Meta.PostTokenBalances == nil {
		return fmt.Errorf("tx.Meta.PostTokenBalances should not be nil")
	}

	return nil
}

func getTransactionAmountForAccountIdx(pre, post []rpc.TransactionMetaTokenBalance, assetAddr, rootPubKey string) (float64, error) {
	var preTokenBalance, postTokenBalance int64
	for _, b := range pre {
		if b.Owner == rootPubKey && b.Mint == assetAddr {
			a, err := strconv.ParseInt(b.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return 0, err
			}

			preTokenBalance = a
			break
		}
	}

	for _, b := range post {
		if b.Owner == rootPubKey && b.Mint == assetAddr {
			a, err := strconv.ParseInt(b.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return 0, err
			}

			postTokenBalance = a
			break
		}
	}

	return float64(postTokenBalance-preTokenBalance) / 1e9, nil
}

func (c *Client) IsTransactionSuccessful(ctx context.Context, txhash string) (bool, error) {
	ss, err := c.solana.GetSignatureStatusWithConfig(ctx, txhash, rpc.GetSignatureStatusesConfig{
		SearchTransactionHistory: true,
	})
	if err != nil {
		return false, errors.Wrap(err, "can't get signature status")
	}
	ok1 := ss != nil && ss.ConfirmationStatus != nil && *ss.ConfirmationStatus == rpc.CommitmentFinalized

	tx, err := c.solana.GetTransactionWithConfig(ctx, txhash, rpc.GetTransactionConfig{
		Encoding:   rpc.GetTransactionConfigEncodingBase64,
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return false, errors.Wrap(err, "can't get transaction by txhash")
	}
	ok2 := tx != nil

	return ok1 || ok2, nil
}

func (s *Client) NeedToRetry(ctx context.Context, latestValidBlockHeight int64) (bool, error) {
	cbh, err := s.GetBlockHeight(ctx)
	if err != nil {
		return false, errors.Wrap(err, "can't get block height")
	}

	return int64(cbh) > latestValidBlockHeight, nil
}

func (c *Client) GetBlockHeight(ctx context.Context) (uint64, error) {
	res := struct {
		GeneralResponse
		Result uint64 `json:"result"`
	}{}

	if err := c.request(ctx, "getBlockHeight", nil, &res); err != nil {
		return 0, err
	}

	return res.Result, nil
}
