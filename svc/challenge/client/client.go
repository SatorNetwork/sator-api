package client

import (
	"context"
	"encoding/json"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		GetByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge.Challenge, error)
		GetChallengesByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)

		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetOneRandomQuestionByChallengeID(ctx context.Context, id uuid.UUID, excludeIDs ...uuid.UUID) (*challenge.Question, error)
		CheckAnswer(ctx context.Context, aid, uid uuid.UUID) (bool, error)

		StoreChallengeAttempt(ctx context.Context, challengeID, userID uuid.UUID) error
		StoreChallengeReceivedRewardAmount(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error
		GetChallengeReceivedRewardAmount(ctx context.Context, challengeID uuid.UUID) (float64, error)
		GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error)
		GetPassedChallengeAttempts(ctx context.Context, challengeID, userID uuid.UUID) (int64, error)

		NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error)
		ListAvailableUserEpisodes(ctx context.Context, userID uuid.UUID) ([]challenge.EpisodeAccess, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetListByShowID returns challenges list filtered by show id
func (c *Client) GetListByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error) {
	if limit < 1 {
		limit = 20
	}
	return c.s.GetChallengesByShowID(ctx, showID, userID, limit, offset)
}

// GetChallengeByID returns Challenge struct
func (c *Client) GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge.Challenge, error) {
	return c.s.GetByID(ctx, challengeID, userID)
}

// GetQuestionsByChallengeID returns questions list filtered by challenge id
func (c *Client) GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]challenge.Question, error) {
	res, err := c.s.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	list := make([]challenge.Question, 0)
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// CheckAnswer ...
func (c *Client) CheckAnswer(ctx context.Context, aid, uid uuid.UUID) (bool, error) {
	return c.s.CheckAnswer(ctx, aid, uid)
}

// GetOneRandomQuestionByChallengeID ...
func (c *Client) GetOneRandomQuestionByChallengeID(ctx context.Context, id uuid.UUID, excludeIDs ...uuid.UUID) (*challenge.Question, error) {
	return c.s.GetOneRandomQuestionByChallengeID(ctx, id, excludeIDs...)
}

// StoreChallengeAttempt user to store challenge attempts.
func (c *Client) StoreChallengeAttempt(ctx context.Context, challengeID, userID uuid.UUID) error {
	return c.s.StoreChallengeAttempt(ctx, challengeID, userID)
}

// StoreChallengeReceivedRewardAmount user to store challenge received reward amount.
func (c *Client) StoreChallengeReceivedRewardAmount(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error {
	return c.s.StoreChallengeReceivedRewardAmount(ctx, challengeID, userID, rewardAmount)
}

// GetChallengeReceivedRewardAmount user to store challenge received reward amount.
func (c *Client) GetChallengeReceivedRewardAmount(ctx context.Context, challengeID uuid.UUID) (float64, error) {
	res, err := c.s.GetChallengeReceivedRewardAmount(ctx, challengeID)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GetChallengeReceivedRewardAmountByUserID user to store challenge received reward amount.
func (c *Client) GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error) {
	res, err := c.s.GetChallengeReceivedRewardAmountByUserID(ctx, challengeID, userID)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// GetPassedChallengeAttempts user to get challenge attempts passed.
func (c *Client) GetPassedChallengeAttempts(ctx context.Context, challengeID, userID uuid.UUID) (int64, error) {
	res, err := c.s.GetPassedChallengeAttempts(ctx, challengeID, userID)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// NumberUsersWhoHaveAccessToEpisode ...
func (c *Client) NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error) {
	res, err := c.s.NumberUsersWhoHaveAccessToEpisode(ctx, episodeID)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// ListAvailableUserEpisodes ...
func (c *Client) ListAvailableUserEpisodes(ctx context.Context, userID uuid.UUID) ([]challenge.EpisodeAccess, error) {
	res, err := c.s.ListAvailableUserEpisodes(ctx, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
