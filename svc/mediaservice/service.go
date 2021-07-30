package example

import (
	"context"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService() *Service {
	return &Service{}
}

// Example ...
func (s *Service) Example(ctx context.Context, uid uuid.UUID) (interface{}, error) {
	return nil, nil
}
