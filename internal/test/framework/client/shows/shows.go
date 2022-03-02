package shows

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/internal/test/framework/client/utils"
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
