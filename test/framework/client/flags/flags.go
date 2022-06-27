package flags

import (
	"bytes"
	"encoding/json"
	"fmt"
	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type FlagsClient struct{}

func New() *FlagsClient {
	return new(FlagsClient)
}

type Flag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateFlagResponse struct {
	Data *Flag
}

func (a *FlagsClient) GetFlags(apiKey string) ([]Flag, error) {
	url := fmt.Sprintf("http://localhost:8080/flags/")
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
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

	var resp []Flag
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp, nil
}

func (a *FlagsClient) UpdateFlag(apiKey string, req *Flag) (*Flag, error) {
	url := fmt.Sprintf("http://localhost:8080/flags/")
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPut, url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
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
	var resp UpdateFlagResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
