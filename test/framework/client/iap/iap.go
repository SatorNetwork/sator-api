package iap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

type InAppClient struct{}

func New() *InAppClient {
	return new(InAppClient)
}

type Empty struct{}

type RegisterInAppPurchaseRequest struct {
	ReceiptData string `json:"receipt_data"`
}

func (c *InAppClient) RegisterInAppPurchase(accessToken string, req *RegisterInAppPurchaseRequest) (*Empty, error) {
	url := "http://localhost:8080/iap"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal request")
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
	if !client_utils.IsStatusCodeSuccess(httpResp.StatusCode) {
		return nil, errors.Errorf("unexpected status code: %v, body: %s", httpResp.StatusCode, rawBody)
	}

	return &Empty{}, nil
}
