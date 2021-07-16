package qrcodes

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/qrcodes/repository"

	"github.com/google/uuid"
)

const (
	// RelationTypeQRcodes indicates that relation type is "qrcodes".
	RelationTypeQRcodes = "qrcodes"
)

type (
	// Service struct
	Service struct {
		qr qrcodeRepository
		rc rewardsClient
	}

	qrcodeRepository interface {
		GetDataByQRCodeID(ctx context.Context, id uuid.UUID) (repository.Qrcode, error)
		AddQRCode(ctx context.Context, arg repository.AddQRCodeParams) (repository.Qrcode, error)
		DeleteQRCodeByID(ctx context.Context, id uuid.UUID) error
		UpdateQRCode(ctx context.Context, arg repository.UpdateQRCodeParams) error
	}

	rewardsClient interface {
		AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
	}

	// Qrcode struct
	Qrcode struct {
		ID           uuid.UUID `json:"id"`
		ShowID       uuid.UUID `json:"show_id"`
		EpisodeID    uuid.UUID `json:"episode_id"`
		StartsAt     time.Time `json:"starts_at"`
		ExpiresAt    time.Time `json:"expires_at"`
		RewardAmount float64   `json:"reward_amount"`
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation
func NewService(qr qrcodeRepository, rc rewardsClient) *Service {
	if qr == nil {
		log.Fatalln("qrcode service is not set")
	}
	if rc == nil {
		log.Fatalln("rewards client is not set")
	}

	return &Service{qr: qr, rc: rc}
}

// GetDataByQRCodeID returns show id and episode id by qrcode id
func (s *Service) GetDataByQRCodeID(ctx context.Context, id, userID uuid.UUID) (interface{}, error) {
	qrcodeData, err := s.qr.GetDataByQRCodeID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get qrcode by id: %w", err)
		}
		return nil, fmt.Errorf("no qrcode with such id:%s, error:%w", id, err)
	}

	now := time.Now()
	if now.Before(qrcodeData.StartsAt) {
		return nil, ErrQRCodeInvalid
	}
	if now.After(qrcodeData.ExpiresAt) {
		return nil, ErrQRCodeExpired
	}
	if qrcodeData.RewardAmount.Float64 > 0 {
		err := s.rc.AddDepositTransaction(ctx, userID, id, RelationTypeQRcodes, qrcodeData.RewardAmount.Float64)
		if err != nil {
			return nil, fmt.Errorf("could not add transaction for user_id=%s and qrcode_id=%s: %w", userID.String(), id.String(), err)
		}
	}

	qrcode := &Qrcode{
		ID:           qrcodeData.ID,
		ShowID:       qrcodeData.ShowID,
		EpisodeID:    qrcodeData.EpisodeID,
		RewardAmount: qrcodeData.RewardAmount.Float64,
	}

	return qrcode, nil
}

// AddQRCode ..
func (s *Service) AddQRCode(ctx context.Context, qr Qrcode) (Qrcode, error) {
	qrCode, err := s.qr.AddQRCode(ctx, repository.AddQRCodeParams{
		ShowID:    qr.ShowID,
		EpisodeID: qr.EpisodeID,
		StartsAt:  qr.StartsAt,
		ExpiresAt: qr.ExpiresAt,
		RewardAmount: sql.NullFloat64{
			Float64: qr.RewardAmount,
			Valid:   true,
		},
	})
	if err != nil {
		return Qrcode{}, fmt.Errorf("could not add qrcode with episode id=%s: %w", qr.EpisodeID, err)
	}

	return castToQrcode(qrCode), nil
}

func castToQrcode(qr repository.Qrcode) Qrcode {
	return Qrcode{
		ID:           qr.ID,
		ShowID:       qr.ShowID,
		EpisodeID:    qr.EpisodeID,
		StartsAt:     qr.StartsAt,
		ExpiresAt:    qr.ExpiresAt,
		RewardAmount: qr.RewardAmount.Float64,
	}
}

// DeleteQRCodeByID ...
func (s *Service) DeleteQRCodeByID(ctx context.Context, id uuid.UUID) error {
	if err := s.qr.DeleteQRCodeByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete qrcode with id=%s:%w", id, err)
	}

	return nil
}

// UpdateQRCode ..
func (s *Service) UpdateQRCode(ctx context.Context, qr Qrcode) error {
	if err := s.qr.UpdateQRCode(ctx, repository.UpdateQRCodeParams{
		ShowID:    qr.ShowID,
		EpisodeID: qr.EpisodeID,
		StartsAt:  qr.StartsAt,
		ExpiresAt: qr.ExpiresAt,
		RewardAmount: sql.NullFloat64{
			Float64: qr.RewardAmount,
			Valid:   true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: qr.ID,
	}); err != nil {
		return fmt.Errorf("could not update qrcode with id=%s:%w", qr.ID, err)
	}

	return nil
}
