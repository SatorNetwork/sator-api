package solana

// Rank struct
type Rank struct {
	MinimalStakingTime int64
	Amount             uint64
}

// InitializeStakePoolInput struct
type InitializeStakePoolInput struct {
	Number uint8 // 0
	Ranks  [4]Rank
}

// StakeInput struct
type StakeInput struct {
	Number   uint8 // 1
	Duration int64
	Amount   uint64
}

// UnstakeInput struct
type UnstakeInput struct {
	Number uint8 // 2
}

// InitializeShowInput struct
type InitializeShowInput struct {
	Number         uint8 // ?
	RewardLockTime int64
}
