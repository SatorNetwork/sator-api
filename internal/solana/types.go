package solana

import "github.com/portto/solana-go-sdk/common"

// InitializeStakePoolInput struct
type InitializeStakePoolInput struct {
	Number uint8
	Ranks  [4]Rank
}

// Rank struct
type Rank struct {
	MinimalStakingTime int64
	Amount             uint64
}

// StakeInput struct
type StakeInput struct {
	Number   uint8
	Duration int64
	Amount   uint64
}

// UnstakeInput struct
type UnstakeInput struct {
	Number uint8
}

// InitializeShowInput struct
type InitializeShowInput struct {
	RewardLockTime int64
}

// InitializeViewer struct
type InitializeViewer struct {
	UserPubKey common.PublicKey
}

const (
	SystemProgram   = "11111111111111111111111111111111"
	SysvarRent      = "SysvarRent111111111111111111111111111111111"
	SysvarClock     = "SysvarC1ock11111111111111111111111111111111"
	SplToken        = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	StakeProgramID  = "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u"
	RewardProgramID = "DajevvE6uo5HtST4EDguRUcbdEMNKNcLWjjNowMRQvZ1"
)
