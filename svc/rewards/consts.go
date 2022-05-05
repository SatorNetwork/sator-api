package rewards

const (
	// TransactionTypeDeposit indicates that transaction type deposit.
	TransactionTypeDeposit = iota + 1
	// TransactionTypeWithdraw indicates that transaction type withdraw.
	TransactionTypeWithdraw
)

type TransactionStatus uint8

const (
	// TransactionStatusAvailable indicates that transaction available to withdraw
	TransactionStatusAvailable TransactionStatus = iota
	// TransactionStatusRequested indicates that transaction requested to withdraw
	TransactionStatusRequested
	// TransactionStatusFailed indicates that transaction failed to withdraw
	TransactionStatusFailed
	// TransactionStatusWithdrawn indicates that transaction withdrawn
	TransactionStatusWithdrawn
)

func NewTransactionStatus(s string) TransactionStatus {
	switch s {
	case "TransactionStatusAvailable":
		return TransactionStatusAvailable
	case "TransactionStatusRequested":
		return TransactionStatusRequested
	case "TransactionStatusFailed":
		return TransactionStatusFailed
	case "TransactionStatusWithdrawn":
		return TransactionStatusWithdrawn
	default:
		return 255
	}
}

func (s TransactionStatus) String() string {
	switch s {
	case TransactionStatusAvailable:
		return "TransactionStatusAvailable"
	case TransactionStatusRequested:
		return "TransactionStatusRequested"
	case TransactionStatusFailed:
		return "TransactionStatusFailed"
	case TransactionStatusWithdrawn:
		return "TransactionStatusWithdrawn"
	default:
		return ""
	}
}
