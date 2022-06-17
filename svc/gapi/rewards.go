package gapi

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

func calculateUserRewardsForGame(conf configer, nftType string, complexity, result int32) (float64, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var key string

	switch result {
	case int32(GameResultWin):
		key = fmt.Sprintf("viewers_%s_%s", nftType, getGameLevelName(complexity))
	case int32(GameResultLose):
		key = fmt.Sprintf("viewers_%s_%s", nftType, GameResultLose.String())
	default:
		return 0, 0, fmt.Errorf("invalid game result")
	}

	viewers, err := conf.GetInt(ctx, key)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get %s: %w", key, err)
	}

	mltpl, err := conf.GetFloat64(ctx, "viewers_multiplier")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get viewers_multiplier: %w", err)
	}

	v := big.NewFloat(float64(viewers))
	rewardsAmount, _ := big.NewFloat(0).Mul(v, big.NewFloat(mltpl)).Float64()

	return rewardsAmount, viewers, nil
}
