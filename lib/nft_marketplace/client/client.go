//go:build !mock_nft_marketplace

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	lib_nft_marketplace "github.com/SatorNetwork/sator-api/lib/nft_marketplace"
)

type Client struct {
	serverHost string
	serverPort int
}

func New(serverHost string, serverPort int) lib_nft_marketplace.Interface {
	return &Client{
		serverHost: serverHost,
		serverPort: serverPort,
	}
}

func (c *Client) getServerEndpoint() string {
	return fmt.Sprintf("%v:%v", c.serverHost, c.serverPort)
}

func (c *Client) PrepareBuyTx(req *lib_nft_marketplace.PrepareBuyTxRequest) (*lib_nft_marketplace.PrepareBuyTxResponse, error) {
	url := fmt.Sprintf("%v/v1/service/market/buy_tx/prepare", c.getServerEndpoint())
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
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

	var resp lib_nft_marketplace.PrepareBuyTxResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func (c *Client) SendPreparedBuyTx(req *lib_nft_marketplace.SendPreparedBuyTxRequest) (*lib_nft_marketplace.SendPreparedBuyTxResponse, error) {
	url := fmt.Sprintf("%v/v1/service/market/buy_tx/prepared/send", c.getServerEndpoint())
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
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

	var resp lib_nft_marketplace.SendPreparedBuyTxResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func IsStatusCodeSuccess(code int) bool {
	return code >= http.StatusOK && code < 300
}
