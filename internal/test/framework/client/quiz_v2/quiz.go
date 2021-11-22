package quiz_v2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/internal/test/framework/client/utils"
)

type QuizClient struct{}

func NewQuizClient() *QuizClient {
	return new(QuizClient)
}

type Empty struct{}

type GetQuizLinkResponse struct {
	Data struct {
		BaseQuizURL     string `json:"base_quiz_url"`
		RecvMessageSubj string `json:"recv_message_subj"`
		SendMessageSubj string `json:"send_message_subj"`
		UserID          string `json:"user_id"`
	} `json:"data"`
}

func (a *QuizClient) GetQuizLink(accessToken, challengeID string) (*GetQuizLinkResponse, error) {
	url := fmt.Sprintf("http://localhost:8080/quiz_v2/%v/play", challengeID)
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

	var resp GetQuizLinkResponse
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}
