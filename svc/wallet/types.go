package wallet

// Predefined wallet types
const (
	WalletTypeSolana   string = "sol"
	WalletTypeSator    string = "sao"
	WalletTypeRewards  string = "rewards"
	WalletTypeEthereum string = "eth"
)

// Predefined  solana account types
const (
	TokenAccount       SolanaAccountType = "token_account"   // custom token account with sator tokens
	GeneralAccount     SolanaAccountType = "general_account" // general account with SOL
	FeePayerAccount    SolanaAccountType = "fee_payer"       // general account with SOL to pay transaction commission
	IssuerAccount      SolanaAccountType = "issuer"          // sator tokens issuer
	DistributorAccount SolanaAccountType = "distributor"     // sator tokens distributor
	AssetAccount       SolanaAccountType = "asset"           // sator token account
	StakePoolAccount   SolanaAccountType = "stake_pool"      // sator stake pool account
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
	ActionStakeTokens   ActionType = "stake_tokens"
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
	case ActionStakeTokens:
		return "Lock"
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
		Order              int32  `json:"order"`
	}
)

type (
	// Transactions list
	Transactions []Transaction

	// Transaction ...
	Transaction struct {
		ID        string  `json:"id"`
		WalletID  string  `json:"wallet_id"`
		TxHash    string  `json:"tx_hash"`
		Amount    float64 `json:"amount"`
		CreatedAt string  `json:"created_at"`
	}

	// PreparedTransferTransaction struct
	PreparedTransferTransaction struct {
		AssetName       string  `json:"asset_name,omitempty"`
		Amount          float64 `json:"amount,omitempty"`
		RecipientAddr   string  `json:"recipient_address,omitempty"`
		Fee             float64 `json:"fee,omitempty"`
		TransactionHash string  `json:"tx_hash,omitempty"`
		SenderWalletID  string  `json:"sender_wallet_id,omitempty"`
	}
)

// Wallet details
type (
	// Wallet ...
	Wallet struct {
		ID                     string    `json:"id"`
		Type                   string    `json:"type"`
		Order                  int32     `json:"order"`
		SolanaAccountAddress   string    `json:"solana_account_address"`
		EthereumAccountAddress string    `json:"ethereum_account_address"`
		Balance                []Balance `json:"balance"`
		Actions                []Action  `json:"actions"`
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

// Stake details
type Stake struct {
	TotalLocked       float64
	LockedByYou       float64
	CurrentMultiplier int32
	AvailableToLock   float64
}

// Predefined token transfer statuses
const (
	TokenTransferStatusPending int32 = iota
	TokenTransferStatusSuccess
	TokenTransferStatusFailed
)
