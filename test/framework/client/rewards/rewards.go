package rewards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
	"github.com/pkg/errors"
)

type RewardsClient struct{}

func New() *RewardsClient {
	return new(RewardsClient)
}

type Wallet struct {
	ID                     string    `json:"id"`
	Type                   string    `json:"type"`
	Order                  int32     `json:"order"`
	SolanaAccountAddress   string    `json:"solana_account_address"`
	EthereumAccountAddress string    `json:"ethereum_account_address"`
	Balance                []Balance `json:"balance"`
	Actions                []Action  `json:"actions"`
}

type Action struct {
	Type string `json:"type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Balance struct
type Balance struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

func (a *RewardsClient) GetRewardsWallet(accessToken, walletID string) (*Wallet, error) {
	url := fmt.Sprintf("http://localhost:8080/rewards/wallet/%s", walletID)
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
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
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp struct {
		Data Wallet `json:"data"`
	}
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp.Data, nil
}

type ClaimRewardsResult struct {
	DisplayAmount   string  `json:"amount"`
	TransactionURL  string  `json:"transaction_url"`
	Amount          float64 `json:"-"`
	TransactionHash string  `json:"-"`
}

func (a *RewardsClient) ClaimRewards(accessToken string) (*ClaimRewardsResult, error) {
	url := fmt.Sprint("http://localhost:8080/rewards/claim")
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
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
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp struct {
		Data ClaimRewardsResult `json:"data"`
	}
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp.Data, nil
}
