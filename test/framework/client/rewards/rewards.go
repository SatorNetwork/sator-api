package rewards

import (
	"encoding/json"
	"fmt"
	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type RewardsClient struct {}

func New() *RewardsClient {
	return new(RewardsClient)
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

	var resp ClaimRewardsResult
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}