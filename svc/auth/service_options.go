package auth

// WithMailService option
// Sets up service to send emails
func WithMailService(m mailer) ServiceOption {
	return func(s *Service) {
		s.mail = m
	}
}

// WithCustomOTPLength option
// Sets up custom OTP code length
func WithCustomOTPLength(n int) ServiceOption {
	return func(s *Service) {
		s.otpLen = n
	}
}
