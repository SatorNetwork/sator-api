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
		rs rewardsService
	}

	qrcodeRepository interface {
		GetDataByQRCodeID(ctx context.Context, id uuid.UUID) (repository.Qrcode, error)
	}

	rewardsService interface {
		AddDepositTransaction(ctx context.Context, userID uuid.UUID, amount float64, relationID uuid.UUID, relationType string) error
	}

	// Qrcode struct
	Qrcode struct {
		ID           string  `json:"id"`
		ShowID       string  `json:"show_id"`
		EpisodeID    string  `json:"episode_id"`
		RewardAmount float64 `json:"reward_amount"`
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
func (s *Service) GetDataByQRCodeID(ctx context.Context, qrcodeID uuid.UUID, userID uuid.UUID) (interface{}, error) {
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
	if qrcodeData.RewardAmount.Float64 > 0 {
		err := s.rs.AddDepositTransaction(ctx, userID, qrcodeData.RewardAmount.Float64, qrcodeID, "??")
		if err != nil {
			return nil, fmt.Errorf("could not add transaction for user_id=%s and qrcode_id=%s: %w", userID.String(), qrcodeID.String(), err)
		}
	}

	qrcode := &Qrcode{
		ID:           qrcodeData.ID.String(),
		ShowID:       qrcodeData.ShowID.String(),
		EpisodeID:    qrcodeData.EpisodeID.String(),
		RewardAmount: qrcodeData.RewardAmount.Float64,
	}

	return qrcode, nil
}
