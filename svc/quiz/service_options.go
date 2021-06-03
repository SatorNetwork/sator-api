package quiz

// WithCustomTokenTTL service option
func WithCustomTokenTTL(ttl int64) ServiceOption {
	return func(s *Service) {
		s.tokenTTL = ttl
	}
}

// WithCustomTokenGenerateFunction service option
func WithCustomTokenGenerateFunction(fn tokenGenFunc) ServiceOption {
	return func(s *Service) {
		s.tokenGenFunc = fn
	}
}

// WithCustomTokenParseFunction service option
func WithCustomTokenParseFunction(fn tokenParseFunc) ServiceOption {
	return func(s *Service) {
		s.tokenParseFunc = fn
	}
}

// // WithCountdown service option
// func WithCountdown(n int) ServiceOption {
// 	return func(s *Service) {
// 		s.countdown = n
// 	}
// }

// // WithTimeBtwQuestion service option
// func WithTimeBtwQuestion(t time.Duration) ServiceOption {
// 	return func(s *Service) {
// 		s.timeBtwQuestion = t
// 	}
// }
