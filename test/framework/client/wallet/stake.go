package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

type Stake struct{}

type SetStakeRequest struct {
	Amount   float64 `json:"amount"`
	WalletID string  `json:"wallet_id"`
	Duration int64   `json:"duration"`
}

type WrappedSetStakeResponse struct {
	Data *SetStakeResponse `json:"data"`
}

type SetStakeResponse struct {
	TotalLocked       int `json:"TotalLocked"`
	LockedByYou       int `json:"LockedByYou"`
	CurrentMultiplier int `json:"CurrentMultiplier"`
	AvailableToLock   int `json:"AvailableToLock"`
}

func (w *WalletClient) GetStake(accessToken, walletID string) (*SetStakeResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/wallets/%v/stake", walletID)

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
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	var resp WrappedSetStakeResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (w *WalletClient) SetStake(accessToken string, req *SetStakeRequest) error {
	url := fmt.Sprintf("http://localhost:8080/wallets/%v/stake", req.WalletID)

	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal create transfer request")
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
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	return nil
}

func (w *WalletClient) Unstake(accessToken, walletID string) error {
	url := fmt.Sprintf("http://localhost:8080/wallets/%v/unstake", walletID)

	reader := bytes.NewReader([]byte{})
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "can't make http request")
	}
	rawBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "can't read response body")
	}
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	return nil
}
