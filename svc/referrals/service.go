package referrals

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		rr referralsRepository
	}

	Data struct {
		UserID       uuid.UUID `json:"user_id"`
		ReferralCode string    `json:"referral_code"`
	}

	referralsRepository interface {
		GetReferralCodeByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(rr referralsRepository) *Service {
	if rr == nil {
		log.Fatalln("referrals repository is not set")
	}

	return &Service{rr: rr}
}

// GetMyReferralCode returns referral code if there is or generate new if not.
// TODO: implement logic. Currently it's mock.
func (s *Service) GetMyReferralCode(ctx context.Context, uid uuid.UUID) (Data, error) {
	return Data{
		UserID:       uid,
		ReferralCode: "ReferralCode",
	}, nil
}

// StoreUserWithValidCode used to validate referral code and store current user.
// TODO: implement logic. Currently it's mock.
func (s *Service) StoreUserWithValidCode(ctx context.Context, uid uuid.UUID, code string) (bool, error) {
	return true, nil
}
