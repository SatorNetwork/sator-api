package challenge

import "time"

// WithCustomVerificationAttempts ...
// Overvrites default number of varification attempts
func WithCustomVerificationAttempts(n int) ServiceOption {
	return func(s *Service) {
		s.attemptsNumber = int64(n)
	}
}

// WithCustomActivatedPeriod ...
// Overvrites default period while realm is activated
func WithCustomActivatedPeriod(p time.Duration) ServiceOption {
	return func(s *Service) {
		s.activatedRealmPeriod = p
	}
}
