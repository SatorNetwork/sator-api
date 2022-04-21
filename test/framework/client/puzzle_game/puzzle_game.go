package puzzle_game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SatorNetwork/gopuzzlegame"
	client_utils "github.com/SatorNetwork/sator-api/test/framework/client/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PuzzleGameClient struct{}

func New() *PuzzleGameClient {
	return new(PuzzleGameClient)
}

type UnlockPuzzleGameRequest struct {
	UnlockOption string `json:"unlock_option" validate:"required"`
}

type TapTileRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PuzzleGame struct {
	// general info
	ID           uuid.UUID `json:"id"`
	EpisodeID    uuid.UUID `json:"episode_id"`
	PrizePool    float64   `json:"prize_pool"`
	Rewards      float64   `json:"rewards,omitempty"`
	BonusRewards float64   `json:"bonus_rewards,omitempty"`
	PartsX       int32     `json:"parts_x"`
	// PartsY     int32     `json:"parts_y"`
	Steps      int32                `json:"steps"`
	StepsTaken int32                `json:"steps_taken,omitempty"`
	Status     int32                `json:"status"`
	Result     int32                `json:"result,omitempty"`
	Tiles      []*gopuzzlegame.Tile `json:"tiles,omitempty"`

	// depends on user role
	Images []PuzzleGameImage `json:"images,omitempty"`
	Image  string            `json:"image,omitempty"`
}

type PuzzleGameImage struct {
	ID      uuid.UUID `json:"id"`
	FileURL string    `json:"file_url"`
}

func (a *PuzzleGameClient) UnlockPuzzleGame(accessToken string, puzzleGameID uuid.UUID, req *UnlockPuzzleGameRequest) (*PuzzleGame, error) {
	url := fmt.Sprintf("http://localhost:8080/puzzle-game/%s/unlock", puzzleGameID.String())
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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

	var resp PuzzleGame
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func (a *PuzzleGameClient) Start(accessToken string, puzzleGameID uuid.UUID) (*PuzzleGame, error) {
	url := fmt.Sprintf("http://localhost:8080/puzzle-game/%s/start", puzzleGameID.String())
	httpReq, err := http.NewRequest(http.MethodPost, url, nil)
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

	var resp PuzzleGame
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}

func (a *PuzzleGameClient) TapTile(accessToken string, puzzleGameID uuid.UUID, req *TapTileRequest) (*PuzzleGame, error) {
	url := fmt.Sprintf("http://localhost:8080/puzzle-game/%s/tap-tile", puzzleGameID.String())
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal create transfer request")
	}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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

	var resp PuzzleGame
	if err := json.Unmarshal(rawBody, &resp); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}

	return &resp, nil
}
