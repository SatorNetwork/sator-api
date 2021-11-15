package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet"
)

type WalletClient struct {
	solanaClient *solana.Client
}

func NewWalletClient() *WalletClient {
	return &WalletClient{
		solanaClient: solana.New("http://localhost:8899"),
	}
}

type WrappedCreateTransferResponse struct {
	Data *CreateTransferResponse `json:"data"`
}

type CreateTransferResponse struct {
	AssetName        string  `json:"asset_name"`
	Amount           float64 `json:"amount"`
	RecipientAddress string  `json:"recipient_address"`
	TxHash           string  `json:"tx_hash"`
	SenderWalletId   string  `json:"sender_wallet_id"`
}

type GetWalletsResponse struct {
	Data []*Wallet `json:"data"`
}

type Wallet struct {
	Id                 string `json:"id"`
	Type               string `json:"type"`
	GetDetailsUrl      string `json:"get_details_url"`
	GetTransactionsUrl string `json:"get_transactions_url"`
	Order              int    `json:"order"`
}

type GetWalletByIDResponse struct {
	Data *WalletDetails `json:"data"`
}

type WalletDetails struct {
	Id                     string `json:"id"`
	Order                  int    `json:"order"`
	SolanaAccountAddress   string `json:"solana_account_address"`
	EthereumAccountAddress string `json:"ethereum_account_address"`
	Balance                []struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	} `json:"balance"`
	Actions []struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"actions"`
}

type GetWalletTxs struct {
	Data []*Tx `json:"data"`
}

type Tx struct {
	Id        string    `json:"id"`
	WalletId  string    `json:"wallet_id"`
	TxHash    string    `json:"tx_hash"`
	Amount    float64   `json:"amount"`
	CreatedAt string `json:"created_at"`
}

func IsStatusCodeSuccess(code int) bool {
	return code >= http.StatusOK && code < 300
}

func (w *WalletClient) GetWallets(accessToken string) ([]*Wallet, error) {
	url := "http://localhost:8080/wallets"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	if !IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp GetWalletsResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (w *WalletClient) GetWalletByID(accessToken string, walletDetailsUrl string) (*WalletDetails, error) {
	url := fmt.Sprintf("http://localhost:8080/%v", walletDetailsUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	if !IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp GetWalletByIDResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (w *WalletClient) GetWalletTxs(accessToken string, walletTransactionsUrl string) ([]*Tx, error) {
	url := fmt.Sprintf("http://localhost:8080/%v", walletTransactionsUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	if !IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp GetWalletTxs
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}


func (w *WalletClient) CreateTransfer(accessToken string, req *wallet.CreateTransferRequest) (*CreateTransferResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/wallets/%v/create-transfer", req.SenderWalletID)
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	if !IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp WrappedCreateTransferResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (w *WalletClient) ConfirmTransfer(accessToken string, req *wallet.ConfirmTransferRequest) error {
	url := fmt.Sprintf("http://localhost:8080/wallets/%v/confirm-transfer", req.SenderWalletID)
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal confirm transfer request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "can't read response body")
	}
	if !IsStatusCodeSuccess(httpResp.StatusCode) {
		return errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	return nil
}