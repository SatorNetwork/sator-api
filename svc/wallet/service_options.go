package wallet

// WithAssetSolanaAddress ...
func WithAssetSolanaAddress(addr string) ServiceOption {
	return func(s *Service) {
		s.satorAssetSolanaAddr = addr
	}
}

// WithStakePoolSolanaAddress ...
func WithStakePoolSolanaAddress(addr string) ServiceOption {
	return func(s *Service) {
		s.stakePoolSolanaPublicKey = addr
	}
}

// WithSolanaFeePayer ...
func WithSolanaFeePayer(addr string, pk []byte) ServiceOption {
	return func(s *Service) {
		s.feePayerSolanaAddr = addr
		s.feePayerSolanaPrivateKey = pk
	}
}

// WithSolanaTokenHolder ...
func WithSolanaTokenHolder(addr string, pk []byte) ServiceOption {
	return func(s *Service) {
		s.tokenHolderSolanaAddr = addr
		s.tokenHolderSolanaPrivateKey = pk
	}
}

// WithMinAmountToTransfer ...
func WithMinAmountToTransfer(amount float64) ServiceOption {
	return func(s *Service) {
		s.minAmountToTransfer = amount
	}
}

// WithFraudDetectionMode ...
func WithFraudDetectionMode(mode bool) ServiceOption {
	return func(s *Service) {
		s.fraudDetectionMode = mode
	}
}
