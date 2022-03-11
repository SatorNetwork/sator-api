package challenge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

type ChallengesClient struct{}

func New() *ChallengesClient {
	return new(ChallengesClient)
}

type Empty struct{}

type ChallengeWrapper struct {
	Data *Challenge `json:"data"`
}

type Challenge struct {
	Id                 string      `json:"id"`
	ShowId             string      `json:"show_id"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	PrizePool          string      `json:"prize_pool"`
	PrizePoolAmount    int         `json:"prize_pool_amount"`
	Players            int         `json:"players"`
	Winners            string      `json:"winners"`
	TimePerQuestion    string      `json:"time_per_question"`
	TimePerQuestionSec int         `json:"time_per_question_sec"`
	Play               string      `json:"play"`
	EpisodeId          interface{} `json:"episode_id"`
	Kind               int         `json:"kind"`
	UserMaxAttempts    int         `json:"user_max_attempts"`
	AttemptsLeft       int         `json:"attempts_left"`
	ReceivedReward     int         `json:"received_reward"`
	ReceivedRewardStr  string      `json:"received_reward_str"`
}

func (a *ChallengesClient) GetChallengeById(accessToken, challengeID string) (*Challenge, error) {
	url := fmt.Sprintf("http://localhost:8080/challenges/%v", challengeID)
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

	var resp ChallengeWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}

// TODO(evg): debug this method
func (a *ChallengesClient) GetQuestionsByChallengeID(accessToken, challengeID string) (*Empty, error) {
	url := fmt.Sprintf("http://localhost:8080/challenges/%v/questions", challengeID)
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

	return &Empty{}, nil
}
