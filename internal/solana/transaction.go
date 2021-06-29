package solana

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/portto/solana-go-sdk/client"
)

type (
	// GetConfirmedTransactionResponse ...
	GetConfirmedTransactionResponse struct {
		BlockTime   int64              `json:"blockTime"`
		Slot        uint64             `json:"slot"`
		Meta        TransactionMeta    `json:"meta"`
		Transaction client.Transaction `json:"transaction"`
	}

	// TransactionMeta ...
	TransactionMeta struct {
		Fee               uint64   `json:"fee"`
		PreBalances       []int64  `json:"preBalances"`
		PostBalances      []int64  `json:"postBalances"`
		LogMessages       []string `json:"logMesssages"`
		InnerInstructions []struct {
			Index        uint64               `json:"index"`
			Instructions []client.Instruction `json:"instructions"`
		} `json:"innerInstructions"`
		Err    interface{}            `json:"err"`
		Status map[string]interface{} `json:"status"`

		// custom fields
		PostTokenBalances []TokenBalance `json:"postTokenBalances"`
		PreTokenBalances  []TokenBalance `json:"preTokenBalances"`
	}

	// TokenBalance ...
	TokenBalance struct {
		AccountIndex  int           `json:"accountIndex"`
		Mint          string        `json:"mint"`
		UITokenAmount UITokenAmount `json:"uiTokenAmount"`
	}

	// UITokenAmount ...
	UITokenAmount struct {
		Amount         string  `json:"amount"`
		Decimals       int     `json:"decimals"`
		UIAmount       float64 `json:"uiAmount"`
		UIAmountString string  `json:"uiAmountString"`
	}

	// ConfirmedTransactionResponse ...
	ConfirmedTransactionResponse struct {
		TxHash        string    `json:"tx_hash"`
		Amount        float64   `json:"amount"`
		AmountString  string    `json:"amount_string"`
		CreatedAtUnix int64     `json:"created_at_unix"`
		CreatedAt     time.Time `json:"created_at"`
	}
)

// GetConfirmedTransaction returns extended transaction details
func (c *Client) GetConfirmedTransaction(ctx context.Context, txhash string) (GetConfirmedTransactionResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmedTransactionResponse `json:"result"`
	}{}
	err := c.request(ctx, "getConfirmedTransaction", []interface{}{txhash, "json"}, &res)
	if err != nil {
		return GetConfirmedTransactionResponse{}, err
	}
	return res.Result, nil
}

// GetConfirmedTransactionForAccount returns transactions list for given account
func (c *Client) GetConfirmedTransactionForAccount(ctx context.Context, accPubKey, txhash string) (ConfirmedTransactionResponse, error) {
	tx, err := c.GetConfirmedTransaction(ctx, txhash)
	if err != nil {
		return ConfirmedTransactionResponse{}, err
	}

	var accountIndex int
	for idx, acc := range tx.Transaction.Message.AccountKeys {
		if acc == accPubKey {
			accountIndex = idx
			break
		}
	}

	amount, err := getTransactionAmountForAccountIdx(tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances, accountIndex)
	if err != nil {
		return ConfirmedTransactionResponse{}, err
	}

	tr := ConfirmedTransactionResponse{
		TxHash:        txhash,
		Amount:        amount,
		AmountString:  fmt.Sprintf("%f", amount),
		CreatedAtUnix: tx.BlockTime,
		CreatedAt:     time.Unix(tx.BlockTime, 0),
	}

	return tr, nil
}

func getTransactionAmountForAccountIdx(pre, post []TokenBalance, accountIndex int) (float64, error) {
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
			a, err :=  strconv.ParseInt(b.UITokenAmount.Amount, 10, 64)
			if err != nil {
				return 0, err
			}

			postTokenBalance = a
			break
		}
	}

	return float64(postTokenBalance - preTokenBalance)/100000000, nil
}
