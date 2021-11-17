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

// WithMasterOTP option
// Sets up master OTP code to use in dev environment
func WithMasterOTPCode(hash string) ServiceOption {
	return func(s *Service) {
		s.masterCode = hash
	}
}

// WithBlacklistEmailDomains option
// Sets up email domains which have to be blocked to signup
func WithBlacklistEmailDomains(domains ...string) ServiceOption {
	return func(s *Service) {
		s.blacklistEmailDomains = domains
	}
}
