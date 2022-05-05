//go:build !mock_solana

package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
)

// GetTransaction returns transaction details
func (c *Client) GetTransaction(ctx context.Context, txhash string) (lib_solana.GetConfirmedTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result lib_solana.GetConfirmedTransactionResponse `json:"result"`
	}{}
	err := c.request(ctx, "getTransaction", []interface{}{txhash, "json"}, &res)
	if err != nil {
		return lib_solana.GetConfirmedTransactionResponse{}, err
	}
	return res.Result, nil
}

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
func (c *Client) GetConfirmedTransactionForAccount(ctx context.Context, accPubKey, txhash string) (lib_solana.ConfirmedTransactionResponse, error) {
	tx, err := c.GetConfirmedTransaction(ctx, txhash)
	if err != nil {
		return lib_solana.ConfirmedTransactionResponse{}, err
	}

	var accountIndex int
	for idx, acc := range tx.Transaction.Message.Accounts {
		if acc.ToBase58() == accPubKey {
			accountIndex = idx
			break
		}
	}

	amount, err := getTransactionAmountForAccountIdx(tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances, accountIndex)
	if err != nil {
		return lib_solana.ConfirmedTransactionResponse{}, err
	}

	tr := lib_solana.ConfirmedTransactionResponse{
		TxHash:        txhash,
		Amount:        amount,
		AmountString:  fmt.Sprintf("%f", amount),
		CreatedAtUnix: tx.BlockTime,
		CreatedAt:     time.Unix(tx.BlockTime, 0),
	}

	return tr, nil
}

func getTransactionAmountForAccountIdx(pre, post []lib_solana.TokenBalance, accountIndex int) (float64, error) {
	var preTokenBalance, postTokenBalance int64
	for _, b := range pre {
		if b.AccountIndex == accountIndex {
			a, err := strconv.ParseInt(b.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return 0, err
			}

			preTokenBalance = a
			break
		}
	}

	for _, b := range post {
		if b.AccountIndex == accountIndex {
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

// Ismail Ibragim
// Backend Developer at Sator with competencies in payment transactions, acquiring, banking and blockchain.