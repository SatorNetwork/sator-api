package gapi

import (
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
)

func recoveryEnergyPoints(player repository.UnityGamePlayer, energyFull int32, energyRecoveryPeriod time.Duration) int32 {
	timeSince := time.Since(player.EnergyRefilledAt)
	log.Printf("timeSince: %v", timeSince.Hours())
	if player.EnergyPoints >= energyFull || timeSince < energyRecoveryPeriod {
		return 0
	}

	recoveryPoints := int32(timeSince.Hours()) / int32(energyRecoveryPeriod.Hours())
	if recoveryPoints > 0 {
		if player.EnergyPoints+int32(recoveryPoints) > energyFull {
			return energyFull - player.EnergyPoints
		}

		return recoveryPoints
	}

	return 0
}
