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
		GetQRCodeDataByID(ctx context.Context, id uuid.UUID) (repository.Qrcode, error)
		GetQRCodesData(ctx context.Context, arg repository.GetQRCodesDataParams) ([]repository.Qrcode, error)
	}

	rewardsClient interface {
		AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
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
func (s *Service) GetDataByQRCodeID(ctx context.Context, id, userID uuid.UUID) (Qrcode, error) {
	qrcodeData, err := s.qr.GetQRCodeDataByID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Qrcode{}, fmt.Errorf("could not get qrcode by id: %w", err)
		}
		return Qrcode{}, fmt.Errorf("no qrcode with such id:%s, error:%w", id, err)
	}

	now := time.Now()
	if now.Before(qrcodeData.StartsAt) {
		return Qrcode{}, ErrQRCodeInvalid
	}
	if now.After(qrcodeData.ExpiresAt) {
		return Qrcode{}, ErrQRCodeExpired
	}
	if qrcodeData.RewardAmount.Float64 > 0 {
		err := s.rc.AddDepositTransaction(ctx, userID, id, RelationTypeQRcodes, qrcodeData.RewardAmount.Float64)
		if err != nil {
			return Qrcode{}, fmt.Errorf("could not add transaction for user_id=%s and qrcode_id=%s: %w", userID.String(), id.String(), err)
		}
	}

	qrcode := Qrcode{
		ID:           qrcodeData.ID.String(),
		ShowID:       qrcodeData.ShowID.String(),
		EpisodeID:    qrcodeData.EpisodeID.String(),
		RewardAmount: qrcodeData.RewardAmount.Float64,
	}

	return qrcode, nil
}

// GetQRCodesData returns list qrcodes.
func (s *Service) GetQRCodesData(ctx context.Context, limit, offset int32) ([]Qrcode, error) {
	list, err := s.qr.GetQRCodesData(ctx, repository.GetQRCodesDataParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Qrcode{}, fmt.Errorf("could not get qrcodes data list: %w", err)
	}

	return castToListQRCodes(list), nil
}

// castToListQRCodes cast []repository.Qrcode to service []Qrcode structure
func castToListQRCodes(source []repository.Qrcode) []Qrcode {
	result := make([]Qrcode, 0, len(source))
	for _, s := range source {
		result = append(result, Qrcode{
			ID:           s.ID.String(),
			ShowID:       s.ShowID.String(),
			EpisodeID:    s.EpisodeID.String(),
			RewardAmount: s.RewardAmount.Float64,
		})
	}

	return result
}
