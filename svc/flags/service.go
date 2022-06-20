package flags

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/flags/alias"
	"github.com/SatorNetwork/sator-api/svc/flags/repository"
)

type (
	Service struct {
		fr flagRepository
	}

	flagRepository interface {
		GetFlagByKey(ctx context.Context, key string) (repository.Flag, error)
		CreateFlag(ctx context.Context, arg repository.CreateFlagParams) (repository.Flag, error)
	}
)

func NewService(
	fr flagRepository,
) *Service {
	return &Service{
		fr: fr,
	}
}

var flags = []repository.CreateFlagParams{
	{
		alias.FlagKeyPuzzleGameRewards.String(),
		alias.FlagValueEnabled.String(),
	},
}

func (s *Service) Init(ctx context.Context) error {
	for i := range flags {
		if _, err := s.fr.CreateFlag(ctx, flags[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetFlagValueByKey(ctx context.Context, key alias.Key) (alias.Value, error) {
	flag, err := s.fr.GetFlagByKey(ctx, key.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return alias.FlagValueDisabled, nil
		}
		return "", errors.Wrap(err, "can't get flags")
	}

	return alias.NewFlagValueFromString(flag.Value)
}
