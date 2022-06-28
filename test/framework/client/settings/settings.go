package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type SettingsClient struct{}

func New() *SettingsClient {
	return new(SettingsClient)
}

type Setting struct {
	Key         string      `json:"key"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	ValueType   string      `json:"value_type"`
	Description string      `json:"description,omitempty"`
}

type GetSettingsResponse struct {
	Data []*Setting
}

type UpdateSettingResponse struct {
	Data *Setting
}

func (a *SettingsClient) AddSetting(apiKey string, req *Setting) (*Setting, error) {
	url := fmt.Sprintf("http://localhost:8080/settings/")
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
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
	var resp UpdateSettingResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (a *SettingsClient) GetSettings(apiKey string) ([]*Setting, error) {
	url := fmt.Sprintf("http://localhost:8080/settings/")
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

	var resp GetSettingsResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (a *SettingsClient) UpdateSetting(apiKey string, req *Setting) (*Setting, error) {
	url := fmt.Sprintf("http://localhost:8080/settings/")
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
	var resp UpdateSettingResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
