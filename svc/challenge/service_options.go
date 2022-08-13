package challenge

import "time"

// WithCustomVerificationAttempts ...
// Overwrites default number of verification attempts.
func WithCustomVerificationAttempts(n int) ServiceOption {
	return func(s *Service) {
		s.attemptsNumber = int64(n)
	}
}

// WithCustomActivatedPeriod ...
// Overwrites default period while realm is activated.
func WithCustomActivatedPeriod(p time.Duration) ServiceOption {
	return func(s *Service) {
		s.activatedRealmPeriod = p
	}
}

// WithChargeForUnlockFunc ...
func WithChargeForUnlockFunc(fn chargeForUnlockFunc) ServiceOption {
	return func(s *Service) {
		s.chargeForUnlockFn = fn
	}
}

// WithDisabledVerification ...
// Disables verification.
func WithDisabledVerification(disable bool) ServiceOption {
	return func(s *Service) {
		s.disableVerification = disable
	}
}
