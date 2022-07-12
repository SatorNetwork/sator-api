package shows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

type ShowsClient struct{}

func New() *ShowsClient {
	return new(ShowsClient)
}

type Show struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	Cover          string      `json:"cover"`
	HasNewEpisode  bool        `json:"has_new_episode"`
	Categories     []uuid.UUID `json:"categories"`
	Description    string      `json:"description"`
	Claps          int64       `json:"claps"`
	RealmsTitle    string      `json:"realms_title"`
	RealmsSubtitle string      `json:"realms_subtitle"`
	Watch          string      `json:"watch"`
	HasNFT         bool        `json:"has_nft"`
}

type Shows struct {
	Data []*Show `json:"data"`
}

type AddShowRequest struct {
	Title          string   `json:"title,omitempty"`
	Cover          string   `json:"cover,omitempty"`
	HasNewEpisode  bool     `json:"has_new_episode,omitempty"`
	Categories     []string `json:"categories,omitempty"`
	Description    string   `json:"description,omitempty"`
	RealmsTitle    string   `json:"realms_title,omitempty"`
	RealmsSubtitle string   `json:"realms_subtitle,omitempty"`
	Watch          string   `json:"watch,omitempty"`
	Status         string   `json:"status,omitempty"`
}

type AddShowResponseWrapper struct {
	Data *AddShowResponse `json:"data"`
}

type AddShowResponse struct {
	Id             string      `json:"id"`
	Title          string      `json:"title"`
	Cover          string      `json:"cover"`
	HasNewEpisode  bool        `json:"has_new_episode"`
	Categories     interface{} `json:"categories"`
	Description    string      `json:"description"`
	Claps          int         `json:"claps"`
	RealmsTitle    string      `json:"realms_title"`
	RealmsSubtitle string      `json:"realms_subtitle"`
	Watch          string      `json:"watch"`
	HasNft         bool        `json:"has_nft"`
}

type AddSeasonRequest struct {
	ShowID       string `json:"show_id" validate:"required,uuid"`
	SeasonNumber int32  `json:"season_number"`
}

type AddSeasonResponseWrapper struct {
	Data *AddSeasonResponse `json:"data"`
}

type AddSeasonResponse struct {
	Id           string      `json:"id"`
	Title        string      `json:"title"`
	SeasonNumber int         `json:"season_number"`
	Episodes     interface{} `json:"episodes"`
	ShowId       string      `json:"show_id"`
}

type AddEpisodeRequest struct {
	ShowID                  string `json:"show_id"`
	SeasonID                string `json:"season_id"`
	EpisodeNumber           int32  `json:"episode_number"`
	Cover                   string `json:"cover,omitempty"`
	Title                   string `json:"title"`
	Description             string `json:"description,omitempty"`
	ReleaseDate             string `json:"release_date,omitempty"`
	ChallengeID             string `json:"challenge_id,omitempty"`
	VerificationChallengeID string `json:"verification_challenge_id,omitempty"`
	HintText                string `json:"hint_text,omitempty"`
	Watch                   string `json:"watch,omitempty"`
	Status                  string `json:"status,omitempty"`
}

type SendTipsRequest struct {
	ReviewID string  `json:"review_id"`
	Amount   float64 `json:"amount"`
}

func (a *ShowsClient) GetShows(apiKey string) (*Shows, error) {
	url := fmt.Sprintf("http://localhost:8080/nft-marketplace/shows")
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

	var resp Shows
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func (a *ShowsClient) AddShow(apiKey string, req *AddShowRequest) (*AddShowResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/shows")
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
	var resp AddShowResponseWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (a *ShowsClient) AddSeason(apiKey string, req *AddSeasonRequest) (*AddSeasonResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/shows/%v/seasons", req.ShowID)
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

	var resp AddSeasonResponseWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (a *ShowsClient) AddEpisode(apiKey string, req *AddEpisodeRequest) error {
	url := fmt.Sprintf("http://localhost:8080/shows/%v/episodes", req.ShowID)
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
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
	var resp AddShowResponseWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return err
	}

	return nil
}

func (a *ShowsClient) SendTipsToReviewAuthor(apiKey string, req *SendTipsRequest) error {
	url := fmt.Sprintf("http://localhost:8080/shows/reviews/%v/tips", req.ReviewID)
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "can't marshal request")
	}
	reader := bytes.NewReader(body)
	httpReq, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return errors.Wrap(err, "can't create http request")
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
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
