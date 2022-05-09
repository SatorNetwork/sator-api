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
func (c *Client) GetConfirmedTransactionForAccount(ctx context.Context, accPubKey, txhash string) (lib_solana.ConfirmedTransactionResponse, error) {
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

	amount, err := getTransactionAmountForAccountIdx(tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances, accPubKey)
	if err != nil {
		err := errors.Wrap(err, "can't get transaction amount for account idx")
		fmt.Println(err)
		return lib_solana.ConfirmedTransactionResponse{}, err
	}

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

func getTransactionAmountForAccountIdx(pre, post []rpc.TransactionMetaTokenBalance, accPubKey string) (float64, error) {
	var preTokenBalance, postTokenBalance int64
	for _, b := range pre {
		if b.Owner == accPubKey {
			a, err := strconv.ParseInt(b.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return 0, err
			}

			preTokenBalance = a
			break
		}
	}

	for _, b := range post {
		if b.Owner == accPubKey {
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
