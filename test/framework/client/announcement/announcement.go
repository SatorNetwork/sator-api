package announcement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

type (
	AnnouncementClient struct{}

	CreateAnnouncementRequest struct {
		Title              string            `json:"title"`
		Description        string            `json:"description"`
		ActionUrl          string            `json:"action_url"`
		StartsAt           int64             `json:"starts_at"`
		EndsAt             int64             `json:"ends_at"`
		Type               string            `json:"type"`
		TypeSpecificParams map[string]string `json:"type_specific_params"`
	}

	CreateAnnouncementResponseWrapper struct {
		Data *CreateAnnouncementResponse `json:"data"`
	}

	CreateAnnouncementResponse struct {
		ID string `json:"id"`
	}

	GetAnnouncementByIDRequest struct {
		ID string `json:"id"`
	}

	AnnouncementsWrapper struct {
		Data []*Announcement `json:"data"`
	}

	AnnouncementWrapper struct {
		Data *Announcement `json:"data"`
	}

	Announcement struct {
		ID                 string            `json:"id"`
		Title              string            `json:"title"`
		Description        string            `json:"description"`
		ActionUrl          string            `json:"action_url"`
		StartsAt           int64             `json:"starts_at"`
		EndsAt             int64             `json:"ends_at"`
		Type               string            `json:"type"`
		TypeSpecificParams map[string]string `json:"type_specific_params"`
	}

	UpdateAnnouncementRequest struct {
		ID                 string            `json:"id"`
		Title              string            `json:"title"`
		Description        string            `json:"description"`
		ActionUrl          string            `json:"action_url"`
		StartsAt           int64             `json:"starts_at"`
		EndsAt             int64             `json:"ends_at"`
		Type               string            `json:"type"`
		TypeSpecificParams map[string]string `json:"type_specific_params"`
	}

	DeleteAnnouncementRequest struct {
		ID string `json:"id"`
	}

	MarkAsReadRequest struct {
		AnnouncementID string `json:"announcement_id"`
	}

	GetAnnouncementTypesResponseWrapper struct {
		Data *GetAnnouncementTypesResponse `json:"data"`
	}

	GetAnnouncementTypesResponse struct {
		Types []string `json:"types"`
	}
)

func New() *AnnouncementClient {
	return new(AnnouncementClient)
}

func (c *AnnouncementClient) CreateAnnouncement(accessToken string, req *CreateAnnouncementRequest) (*CreateAnnouncementResponse, error) {
	url := "http://localhost:8080/announcement"
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

	var resp CreateAnnouncementResponseWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (c *AnnouncementClient) GetAnnouncementByID(accessToken string, req *GetAnnouncementByIDRequest) (*Announcement, error) {
	if req.ID == "" {
		return nil, errors.Errorf("ID should not be empty")
	}

	url := fmt.Sprintf("http://localhost:8080/announcement/%v", req.ID)
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

	var resp AnnouncementWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (c *AnnouncementClient) UpdateAnnouncement(accessToken string, req *UpdateAnnouncementRequest) error {
	url := fmt.Sprintf("http://localhost:8080/announcement/%v", req.ID)
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPut, url, reader)
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

func (c *AnnouncementClient) DeleteAnnouncement(accessToken string, req *DeleteAnnouncementRequest) error {
	url := fmt.Sprintf("http://localhost:8080/announcement/%v", req.ID)
	httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
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

func (c *AnnouncementClient) ListAnnouncements(accessToken string) ([]*Announcement, error) {
	url := "http://localhost:8080/announcement"
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

	var resp AnnouncementsWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (c *AnnouncementClient) ListUnreadAnnouncements(accessToken string) ([]*Announcement, error) {
	url := "http://localhost:8080/announcement/unread"
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

	var resp AnnouncementsWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (c *AnnouncementClient) ListActiveAnnouncements(accessToken string) ([]*Announcement, error) {
	url := "http://localhost:8080/announcement/active"
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

	var resp AnnouncementsWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

func (c *AnnouncementClient) MarkAsRead(accessToken string, req *MarkAsReadRequest) error {
	url := fmt.Sprintf("http://localhost:8080/announcement/%v/read", req.AnnouncementID)
	httpReq, err := http.NewRequest(http.MethodPost, url, nil)
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

func (c *AnnouncementClient) MarkAllAsRead(accessToken string) error {
	url := fmt.Sprintf("http://localhost:8080/announcement/read_all")
	httpReq, err := http.NewRequest(http.MethodPost, url, nil)
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

func (c *AnnouncementClient) GetAnnouncementTypes(accessToken string) (*GetAnnouncementTypesResponse, error) {
	url := "http://localhost:8080/announcement/types"
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

	var resp GetAnnouncementTypesResponseWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}
