package errors

var (
	ErrCantSendSolanaTransaction = newServiceError("cant send solana transaction", 1000)
)
