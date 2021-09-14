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

const (
	SystemProgram = "11111111111111111111111111111111"
	SysvarRent = "SysvarRent111111111111111111111111111111111"
	SysvarClock = "SysvarC1ock11111111111111111111111111111111"
	SplToken = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	ProgramID = "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u"

)