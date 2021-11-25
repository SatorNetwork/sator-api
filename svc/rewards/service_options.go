package rewards

import "time"

// WithAssetName option
// Default value: SAO
func WithAssetName(assetName string) Option {
	return func(s *Service) {
		s.assetName = assetName
	}
}

// WithExplorerURLTmpl option
// Default value: "https://explorer.solana.com/tx/%s?cluster=devnet"
func WithExplorerURLTmpl(explorerURLTmpl string) Option {
	return func(s *Service) {
		s.explorerURLTmpl = explorerURLTmpl
	}
}

// WithHoldRewardsPeriod option
// Default value: 30 days
func WithHoldRewardsPeriod(period time.Duration) Option {
	return func(s *Service) {
		s.holdRewardsPeriod = period
	}
}

// WithMinAmountToClaim ...
func WithMinAmountToClaim(amount float64) Option {
	return func(s *Service) {
		s.minAmountToClaim = amount
	}
}
