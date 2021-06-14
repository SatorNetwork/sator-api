package rewards

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
