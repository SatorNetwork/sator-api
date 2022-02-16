package solana

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
