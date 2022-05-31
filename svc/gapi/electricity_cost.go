package gapi

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

const (
	electricityCostModeWin    = "win"
	electricityCostModeLose   = "lose"
	electricityCostModeAlways = "always"
)

func calculateElectricityCost(conf configer, nftType string, result int32, rewards float64) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mode, err := conf.GetString(ctx, "electricity_percent_mode")
	if err != nil {
		mode = electricityCostModeAlways
	}

	if (mode == electricityCostModeWin && result == 1) ||
		(mode == electricityCostModeLose && result == 0) ||
		mode == electricityCostModeAlways {

		feePercent, err := conf.GetInt(ctx, fmt.Sprintf("electricity_percent_%s", nftType))
		if err != nil || feePercent == 0 {
			feePercent = 40
		}

		rew := big.NewFloat(rewards)
		fee := big.NewFloat(float64(feePercent))
		res := big.NewFloat(0).Mul(rew, fee)
		result, _ := big.NewFloat(0).Quo(res, big.NewFloat(100)).Float64()

		return result, nil
	}

	return 0, nil
}
