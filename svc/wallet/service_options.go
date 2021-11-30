package wallet

// WithAssetSolanaAddress ...
func WithAssetSolanaAddress(addr string) ServiceOption {
	return func(s *Service) {
		s.satorAssetSolanaAddr = addr
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
