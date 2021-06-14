package qrcodes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/qrcodes/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		pr qrcodeRepository
	}

	qrcodeRepository interface {
		GetDataByQRCodeID(ctx context.Context, id uuid.UUID) (repository.Qrcode, error)
	}

	// Qrcode struct
	Qrcode struct {
		ID        string `json:"id"`
		ShowID    string `json:"show_id"`
		EpisodeID string `json:"episode_id"`
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation
func NewService(pr qrcodeRepository) *Service {
	if pr == nil {
		log.Fatalln("qrcode repository is not set")
	}
	return &Service{pr: pr}
}

// GetDataByQRCodeID returns show id and episode id by qrcode id
func (s *Service) GetDataByQRCodeID(ctx context.Context, qrcodeID uuid.UUID) (interface{}, error) {
	qrcodeData, err := s.pr.GetDataByQRCodeID(ctx, qrcodeID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not qrcode data: %w", err)
		}
	}

	now := time.Now()
	if now.Before(qrcodeData.StartsAt) {
		return nil, ErrQRCodeInvalid
	}
	if now.After(qrcodeData.ExpiresAt) {
		return nil, ErrQRCodeExpired
	}

	qrcode := &Qrcode{
		ID:        qrcodeData.ID.String(),
		ShowID:    qrcodeData.ShowID.String(),
		EpisodeID: qrcodeData.EpisodeID.String(),
	}

	return qrcode, nil
}
