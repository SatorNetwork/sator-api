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

type InitializeQuizInput struct {
	Winners     []WinnerInput
	TokenAmount uint64
}

type WinnerInput struct {
	UserPubKey common.PublicKey
	Points     uint32
}
