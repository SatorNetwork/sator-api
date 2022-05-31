package gapi

import (
	"time"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
)

func recoveryEnergyPoints(player repository.UnityGamePlayer, energyFull int32, energyRecoveryPeriod time.Duration) int32 {
	if energyFull == 0 {
		energyFull = 3
	}

	if energyRecoveryPeriod == 0 {
		energyRecoveryPeriod = time.Hour * 4
	}

	timeSince := time.Since(player.EnergyRefilledAt)
	if player.EnergyPoints >= energyFull || timeSince < energyRecoveryPeriod {
		return 0
	}

	recoveryPoints := int32(timeSince.Seconds()) / int32(energyRecoveryPeriod.Seconds())
	if recoveryPoints > 0 {
		if player.EnergyPoints+int32(recoveryPoints) > energyFull {
			return energyFull - player.EnergyPoints
		}

		return recoveryPoints
	}

	return 0
}
