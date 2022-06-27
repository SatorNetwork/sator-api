package alias

import "github.com/pkg/errors"

type Key string

const (
	FlagKeyUndefined           Key = "UNDEFINED"
	FlagKeyPuzzleGameRewards   Key = "PUZZLE_GAME_REWARDS"
	FlagKeyPuzzleGamePaidSteps Key = "PUZZLE_GAME_PAID_STEPS"
)

func NewFlagKeyFromString(s string) (Key, error) {
	switch s {
	case "PUZZLE_GAME_REWARDS":
		return FlagKeyPuzzleGameRewards, nil
	case "PUZZLE_GAME_PAID_STEPS":
		return FlagKeyPuzzleGamePaidSteps, nil
	default:
		return FlagKeyUndefined, errors.Errorf("flags key with such name %v doesn't exist", s)
	}
}

func (k Key) String() string {
	switch k {
	case FlagKeyPuzzleGameRewards:
		return "PUZZLE_GAME_REWARDS"
	case FlagKeyPuzzleGamePaidSteps:
		return "PUZZLE_GAME_PAID_STEPS"
	default:
		return "UNDEFINED"
	}
}
