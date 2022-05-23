package gapi

import "time"

// WithMinVersion sets the minimum version of the unity game client
func WithMinVersion(minVersion string) ServiceOption {
	return func(s *Service) {
		s.minVersion = minVersion
	}
}

// WithEnergyFull sets the energy full value
func WithEnergyFull(energyFull int32) ServiceOption {
	return func(s *Service) {
		s.energyFull = energyFull
	}
}

// WithEnergyRecoveryPeriod sets the energy recovery period
func WithEnergyRecoveryPeriod(energyRecoveryPeriod time.Duration) ServiceOption {
	return func(s *Service) {
		s.energyRecoveryPeriod = energyRecoveryPeriod
	}
}

// WithMinRewardsToClaim sets the minimum rewards to claim
func WithMinRewardsToClaim(minRewardsToClaim float64) ServiceOption {
	return func(s *Service) {
		s.minRewardsToClaim = minRewardsToClaim
	}
}
