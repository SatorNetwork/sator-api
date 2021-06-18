package wallet

// Predefined wallet types
const (
	WalletTypeSolana  string = "sol"
	WalletTypeSator   string = "sao"
	WalletTypeRewards string = "rewards"
)

// Predefined  solana account types
const (
	TokenAccount       SolanaAccountType = "token_account"   // custom token account with sator tokens
	GeneralAccount     SolanaAccountType = "general_account" // general account with SOL
	FeePayerAccount    SolanaAccountType = "fee_payer"       // general account with SOL to pay transaction comission
	IssuerAccount      SolanaAccountType = "issuer"          // sator tokens issuer
	DistributorAccount SolanaAccountType = "distributor"     // sator tokens distributor
	AssetAccount       SolanaAccountType = "asset"           // sator token account
)

// SolanaAccountType solana account type
type SolanaAccountType string

func (t SolanaAccountType) String() string {
	return string(t)
}

// Predefined wallet action types
const (
	ActionClaimRewards  string = "claim_rewards"
	ActionSendTokens    string = "send_tokens"
	ActionReceiveTokens string = "receive_tokens"
)

type (
	// Wallets list
	Wallets []WalletsListItem

	// WalletsListItem ...
	WalletsListItem struct {
		ID                 string `json:"id"`
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
