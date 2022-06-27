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
		CreateFlag(ctx context.Context, arg repository.CreateFlagParams) error
		UpdateFlag(ctx context.Context, arg repository.UpdateFlagParams) (repository.Flag, error)
		GetFlags(ctx context.Context) ([]repository.Flag, error)
	}
)

func NewService(
	fr flagRepository,
) *Service {
	return &Service{
		fr: fr,
	}
}

var initFlags = []repository.CreateFlagParams{
	{
		alias.FlagKeyPuzzleGameRewards.String(),
		alias.FlagValueEnabled.String(),
	},
	{
		alias.FlagKeyPuzzleGamePaidSteps.String(),
		alias.FlagValueEnabled.String(),
	},
}

func (s *Service) Init(ctx context.Context) error {
	for i := range initFlags {
		if err := s.fr.CreateFlag(ctx, initFlags[i]); err != nil {
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

func (s *Service) UpdateFlag(ctx context.Context, flag *repository.Flag) (*repository.Flag, error) {
	newFlag, err := s.fr.UpdateFlag(ctx, repository.UpdateFlagParams{
		Value: flag.Value,
		Key:   flag.Key,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't update flags")
	}

	return &newFlag, nil
}

func (s *Service) GetFlags(ctx context.Context) ([]repository.Flag, error) {
	flags, err := s.fr.GetFlags(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get flags")
	}

	return flags, nil
}
