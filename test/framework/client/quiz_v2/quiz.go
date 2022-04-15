package quiz_v2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
)

type QuizClient struct{}

func New() *QuizClient {
	return new(QuizClient)
}

type Empty struct{}

type GetQuizLinkResponse struct {
	Data struct {
		BaseQuizURL     string `json:"base_quiz_url"`
		RecvMessageSubj string `json:"recv_message_subj"`
		SendMessageSubj string `json:"send_message_subj"`
		UserID          string `json:"user_id"`
		ServerPublicKey string `json:"server_public_key"`
	} `json:"data"`
}

type ChallengeWrapper struct {
	Data *Challenge `json:"data"`
}

type Challenge struct {
	Id                     string      `json:"id"`
	ShowId                 string      `json:"show_id"`
	Title                  string      `json:"title"`
	Description            string      `json:"description"`
	PrizePool              string      `json:"prize_pool"`
	PrizePoolAmount        int         `json:"prize_pool_amount"`
	Players                int         `json:"players"`
	Winners                string      `json:"winners"`
	TimePerQuestion        string      `json:"time_per_question"`
	TimePerQuestionSec     int         `json:"time_per_question_sec"`
	Play                   string      `json:"play"`
	EpisodeId              interface{} `json:"episode_id"`
	Kind                   int         `json:"kind"`
	UserMaxAttempts        int         `json:"user_max_attempts"`
	AttemptsLeft           int         `json:"attempts_left"`
	ReceivedReward         int         `json:"received_reward"`
	ReceivedRewardStr      string      `json:"received_reward_str"`
	RegisteredPlayersInDB  int         `json:"registered_players_in_db"`
	CurrentPrizePool       string      `json:"current_prize_pool"`
	CurrentPrizePoolAmount float64     `json:"current_prize_pool_amount"`
}

type ChallengeWithPlayer struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	PlayersToStart   int    `json:"players_to_start"`
	PlayersNumber    int    `json:"players_number"`
	PrizePool        string `json:"prize_pool"`
	IsRealmActivated bool   `json:"is_realm_activated"`
}

type ChallengesWithPlayerWrapper struct {
	Data []*ChallengeWithPlayer
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

func (a *QuizClient) GetChallengeById(accessToken, challengeID string) (*Challenge, error) {
	url := fmt.Sprintf("http://localhost:8080/quiz_v2/challenges/%v", challengeID)
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

func (a *QuizClient) GetChallengesSortedByPlayers(accessToken string) ([]*ChallengeWithPlayer, error) {
	url := "http://localhost:8080/quiz_v2/challenges/sorted_by_players"
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

	var resp ChallengesWithPlayerWrapper
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return resp.Data, nil
}
