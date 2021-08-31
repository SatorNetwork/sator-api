package referrals

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/firebase"
	"github.com/SatorNetwork/sator-api/svc/referrals/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		rr     referralsRepository
		fb     *firebase.Interactor
		config firebase.Config
	}

	ReferralCode struct {
		ID           uuid.UUID `json:"id"`
		Title        string    `json:"title"`
		Code         string    `json:"code"`
		ReferralLink string    `json:"referral_link"`
		IsPersonal   bool      `json:"is_personal"`
		UserID       uuid.UUID `json:"user_id"`
		CreatedAt    time.Time `json:"created_at"`
	}

	Referral struct {
		ReferralCodeID uuid.UUID `json:"referral_code_id"`
		UserID         uuid.UUID `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
	}

	referralsRepository interface {
		// Referral codes
		AddReferralCodeData(ctx context.Context, arg repository.AddReferralCodeDataParams) (repository.ReferralCode, error)
		DeleteReferralCodeDataByID(ctx context.Context, id uuid.UUID) error
		GetReferralCodeDataByUserID(ctx context.Context, userID uuid.UUID) ([]repository.ReferralCode, error)
		GetReferralCodesDataList(ctx context.Context) ([]repository.ReferralCode, error)
		UpdateReferralCodeData(ctx context.Context, arg repository.UpdateReferralCodeDataParams) error

		// Referrals
		AddReferral(ctx context.Context, arg repository.AddReferralParams) error
		GetReferralsWithPaginationByUserID(ctx context.Context, arg repository.GetReferralsWithPaginationByUserIDParams) ([]repository.Referral, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(rr referralsRepository, fb *firebase.Interactor, config firebase.Config) *Service {
	if rr == nil {
		log.Fatalln("referrals repository is not set")
	}
	if fb == nil {
		log.Fatalln("firebase client is not set")
	}

	return &Service{rr: rr, fb: fb, config: config}
}

// GetMyReferralCode returns referral code if there is or generate new if not.
func (s *Service) GetMyReferralCode(ctx context.Context, uid uuid.UUID, username string) ([]ReferralCode, error) {
	referralCodeData, err := s.rr.GetReferralCodeDataByUserID(ctx, uid)
	if err != nil {
		return []ReferralCode{}, fmt.Errorf("could not get referral code data by user id: %w", err)
	}
	if len(referralCodeData) < 1 {
		id := uuid.New()
		link, err := s.fb.GenerateDynamicLink(ctx, firebase.DynamicLinkRequest{
			DynamicLinkInfo: firebase.DynamicLinkInfo{
				DomainUriPrefix: s.config.BaseFirebaseURL,
				Link:            s.config.MainSiteLink + "referral/" + id.String(),
				AndroidInfo: firebase.AndroidInfo{
					AndroidPackageName: s.config.AndroidPackageName,
				},
				IosInfo: firebase.IosInfo{
					IosBundleId: s.config.IosBundleId,
				},
				//NavigationInfo: firebase.NavigationInfo{EnableForcedRedirect: true},
			},
			Suffix: firebase.Suffix{
				Option: s.config.SuffixOption,
			},
		})
		if err != nil {
			return []ReferralCode{}, fmt.Errorf("could not generate dynamic link for user = %v: %w", uid, err)
		}

		data, err := s.rr.AddReferralCodeData(ctx, repository.AddReferralCodeDataParams{
			ID: id,
			Title: sql.NullString{
				String: username,
				Valid:  len(username) > 0,
			},
			Code: id.String(),
			ReferralLink: sql.NullString{
				String: link.ShortLink,
				Valid:  len(link.ShortLink) > 0,
			},
			IsPersonal: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			UserID: uid,
		})
		if err != nil {
			return []ReferralCode{}, fmt.Errorf("could not add referral code data for user = %v: %w", uid, err)
		}
		return castReferralCodeToReferralCodes(data), err
	}

	return castToReferralCodesWithUsername(referralCodeData, username), nil
}

// Cast repository.ReferralCode to service ReferralCode structure
func castToReferralCode(source repository.ReferralCode, username string) ReferralCode {
	if source.Title.String == "" {
		source.Title = sql.NullString{
			String: username,
			Valid:  true,
		}
	}

	return ReferralCode{
		ID:           source.ID,
		Title:        source.Title.String,
		Code:         source.Code,
		ReferralLink: source.ReferralLink.String,
		IsPersonal:   source.IsPersonal.Bool,
		UserID:       source.UserID,
		CreatedAt:    source.CreatedAt,
	}
}

// Cast repository.ReferralCode array to service ReferralCode array.
func castToReferralCodesWithUsername(source []repository.ReferralCode, username string) []ReferralCode {
	result := make([]ReferralCode, 0, len(source))
	for _, s := range source {
		if s.Title.String == "" {
			s.Title = sql.NullString{
				String: username,
				Valid:  true,
			}
		}
		result = append(result, ReferralCode{
			ID:           s.ID,
			Title:        s.Title.String,
			Code:         s.Code,
			ReferralLink: s.ReferralLink.String,
			IsPersonal:   s.IsPersonal.Bool,
			UserID:       s.UserID,
			CreatedAt:    s.CreatedAt,
		})
	}

	return result
}

// Cast repository.ReferralCode to service ReferralCode array.
func castReferralCodeToReferralCodes(source repository.ReferralCode) []ReferralCode {
	result := make([]ReferralCode, 0, 1)
	result = append(result, ReferralCode{
		ID:           source.ID,
		Title:        source.Title.String,
		Code:         source.Code,
		ReferralLink: source.ReferralLink.String,
		IsPersonal:   source.IsPersonal.Bool,
		UserID:       source.UserID,
		CreatedAt:    source.CreatedAt,
	})

	return result
}

// AddReferralCodeData used to store new referral code.
func (s *Service) AddReferralCodeData(ctx context.Context, rc ReferralCode) (ReferralCode, error) {
	id := uuid.New()
	link, err := s.fb.GenerateDynamicLink(ctx, firebase.DynamicLinkRequest{
		DynamicLinkInfo: firebase.DynamicLinkInfo{
			DomainUriPrefix: s.config.BaseFirebaseURL,
			Link:            s.config.MainSiteLink + "referral/" + id.String(),
			AndroidInfo: firebase.AndroidInfo{
				AndroidPackageName: s.config.AndroidPackageName,
			},
			IosInfo: firebase.IosInfo{
				IosBundleId: s.config.IosBundleId,
			},
			//NavigationInfo: firebase.NavigationInfo{EnableForcedRedirect: true},
		},
		Suffix: firebase.Suffix{
			Option: s.config.SuffixOption,
		},
	})
	if err != nil {
		return ReferralCode{}, fmt.Errorf("could not generate dynamic link: %w", err)
	}

	referralCode, err := s.rr.AddReferralCodeData(ctx, repository.AddReferralCodeDataParams{
		ID: id,
		Title: sql.NullString{
			String: rc.Title,
			Valid:  len(rc.Title) > 0,
		},
		Code: rc.Code,
		ReferralLink: sql.NullString{
			String: link.ShortLink,
			Valid:  len(link.ShortLink) > 0,
		},
		IsPersonal: sql.NullBool{
			Bool:  rc.IsPersonal,
			Valid: true,
		},
		UserID: rc.UserID,
	})
	if err != nil {
		return ReferralCode{}, fmt.Errorf("could not store referral code: %w", err)
	}

	return castToReferralCode(referralCode, ""), nil
}

// DeleteReferralCodeDataByID used to delete referral code by id.
func (s *Service) DeleteReferralCodeDataByID(ctx context.Context, id uuid.UUID) error {
	if err := s.rr.DeleteReferralCodeDataByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete show referral code id=%s:%w", id, err)
	}

	return nil
}

// GetReferralCodesDataList returns referral codes list.
func (s *Service) GetReferralCodesDataList(ctx context.Context) ([]ReferralCode, error) {
	referralCodes, err := s.rr.GetReferralCodesDataList(ctx)
	if err != nil {
		return []ReferralCode{}, fmt.Errorf("could not get referral codes list: %w", err)
	}

	return castToReferralCodes(referralCodes), nil
}

// Cast repository.ReferralCode to service ReferralCode structure.
func castToReferralCodes(source []repository.ReferralCode) []ReferralCode {
	result := make([]ReferralCode, 0, len(source))
	for _, s := range source {
		result = append(result, ReferralCode{
			ID:         s.ID,
			Title:      s.Title.String,
			Code:       s.Code,
			IsPersonal: s.IsPersonal.Bool,
			UserID:     s.UserID,
			CreatedAt:  s.CreatedAt,
		})
	}

	return result
}

// UpdateReferralCodeData used to update referral code.
func (s *Service) UpdateReferralCodeData(ctx context.Context, rc ReferralCode) error {
	if err := s.rr.UpdateReferralCodeData(ctx, repository.UpdateReferralCodeDataParams{
		Title: sql.NullString{
			String: rc.Title,
			Valid:  len(rc.Title) > 0,
		},
		Code: rc.Code,
		ReferralLink: sql.NullString{
			String: rc.ReferralLink,
			Valid:  len(rc.ReferralLink) > 0,
		},
		IsPersonal: sql.NullBool{
			Bool:  rc.IsPersonal,
			Valid: true,
		},
		UserID: rc.UserID,
		ID:     rc.ID,
	}); err != nil {
		return fmt.Errorf("could not update referral code with id=%s:%w", rc.ID, err)
	}

	return nil
}

// GetReferralsWithPaginationByUserID returns paginated list users referrals.
func (s *Service) GetReferralsWithPaginationByUserID(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]Referral, error) {
	referral, err := s.rr.GetReferralsWithPaginationByUserID(ctx, repository.GetReferralsWithPaginationByUserIDParams{
		UserID: uid,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Referral{}, fmt.Errorf("could not get referrals list: %w", err)
	}

	return castToListReferral(referral), nil
}

// Cast repository.Referral to service Referral structure
func castToListReferral(source []repository.Referral) []Referral {
	result := make([]Referral, 0, len(source))
	for _, s := range source {
		result = append(result, Referral{
			ReferralCodeID: s.ReferralCodeID,
			UserID:         s.UserID,
			CreatedAt:      s.CreatedAt,
		})
	}

	return result
}

// StoreUserWithValidCode used to validate referral code and store current user.
func (s *Service) StoreUserWithValidCode(ctx context.Context, uid uuid.UUID, code string) (bool, error) {
	list, err := s.rr.GetReferralCodesDataList(ctx)
	if err != nil {
		return false, fmt.Errorf("could not get referral codes list: %w", err)
	}
	for _, v := range list {
		if v.Code == code {
			if v.UserID == uid {
				return false, errors.New("could not referral yourself")
			}
			err := s.rr.AddReferral(ctx, repository.AddReferralParams{
				ReferralCodeID: v.ID,
				UserID:         uid,
			})
			if err != nil {
				return false, fmt.Errorf("could not store referral with id = %v, and code = %v: %w", uid, code, err)
			}

			return true, nil
		}
	}

	return false, fmt.Errorf("referral code = %v is not found: %w", code, err)
}
