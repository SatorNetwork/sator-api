package puzzle_game

// WithChargeFunction sets the charge function for the service.
func WithChargeFunction(fn chargeForUnlockFunc) ServiceOption {
	return func(s *Service) {
		s.chargeForUnlock = fn
	}
}

// WithFileServiceClient sets the file service client for the service.
func WithFileServiceClient(fs filesService) ServiceOption {
	return func(s *Service) {
		s.filesSvc = fs
	}
}

// WithRewardsFunction sets the rewards function for the service.
func WithRewardsFunction(fn rewardsFunc) ServiceOption {
	return func(s *Service) {
		s.rewardsFn = fn
	}
}

// WithUserMultiplierFunction sets the get user multiplier function for the service.
func WithUserMultiplierFunction(fn getUserRewardsMultiplierFunc) ServiceOption {
	return func(s *Service) {
		s.getUserRewardsMultiplierFn = fn
	}
}

// WithSettingsServiceClient sets the settings service client for the service.
func WithSettingsServiceClient(sc settingsService) ServiceOption {
	return func(s *Service) {
		s.settingsSvc = sc
	}
}
