package wallet

// Predefined wallet types
const (
	TypePersonal string = "personal"
	TypeRewards  string = "rewards"
)

// SolanaAccountType solana account type
type SolanaAccountType string

func (t SolanaAccountType) String() string {
	return string(t)
}

// ActionType of wallet
type ActionType string

// Predefined wallet action types
const (
	ActionClaimRewards  ActionType = "claim_rewards"
	ActionSendTokens    ActionType = "send_tokens"
	ActionReceiveTokens ActionType = "receive_tokens"
)

// Name of action type
func (at ActionType) Name() string {
	switch at {
	case ActionClaimRewards:
		return "Claim rewards"
	case ActionSendTokens:
		return "Send"
	case ActionReceiveTokens:
		return "Receive"
	}
	return "Undefined"
}

func (at ActionType) String() string {
	return string(at)
}

type (
	// Wallets list
	Wallets []WalletsListItem

	// WalletsListItem ...
	WalletsListItem struct {
		ID                 string `json:"id"`
		Type               string `json:"type"`
		GetDetailsURL      string `json:"get_details_url"`      // url to get wallet details
		GetTransactionsURL string `json:"get_transactions_url"` // url to get transactions list
	}
)

type (
	// Transactions list
	Transactions []Transaction

	// Transaction ...
	Transaction struct {
		TxHash    string  `json:"tx_hash"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"created_at"`
	}
)

// Wallet details
type (
	// Wallet ...
	Wallet struct {
		ID                   string    `json:"id"`
		SolanaAccountAddress string    `json:"solana_account_address"`
		Balance              []Balance `json:"balance"`
		Actions              []Action  `json:"actions"`
	}

	// Action ...
	Action struct {
		Type string `json:"type"`
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// Balance struct
	Balance struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
)
