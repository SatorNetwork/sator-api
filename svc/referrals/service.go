package referrals

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
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
		ID           uuid.UUID  `json:"id"`
		Title        string     `json:"title"`
		Code         string     `json:"code"`
		ReferralLink string     `json:"referral_link"`
		IsPersonal   bool       `json:"is_personal"`
		UserID       *uuid.UUID `json:"user_id"`
		CreatedAt    time.Time  `json:"created_at"`
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
		GetReferralCodeDataByUserID(ctx context.Context, userID uuid.NullUUID) (repository.ReferralCode, error)
		GetReferralCodeDataByCode(ctx context.Context, code string) (repository.ReferralCode, error)
		GetReferralCodesDataList(ctx context.Context, arg repository.GetReferralCodesDataListParams) ([]repository.ReferralCode, error)
		UpdateReferralCodeData(ctx context.Context, arg repository.UpdateReferralCodeDataParams) error
		GetNumberOfReferralCodes(ctx context.Context) (int64, error)

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
func (s *Service) GetMyReferralCode(ctx context.Context, uid uuid.UUID, username string) (ReferralCode, error) {
	referralCodeData, err := s.rr.GetReferralCodeDataByUserID(ctx, uuid.NullUUID{
		UUID:  uid,
		Valid: true,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
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
				return ReferralCode{}, fmt.Errorf("could not generate dynamic link for user = %v: %w", uid, err)
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
				UserID: uuid.NullUUID{
					UUID:  uid,
					Valid: true,
				},
			})
			if err != nil {
				return ReferralCode{}, fmt.Errorf("could not add referral code data for user = %v: %w", uid, err)
			}

			return castToReferralCode(data, username), err
		}

		return ReferralCode{}, fmt.Errorf("could not get referral code data by user id: %w", err)
	}

	return castToReferralCode(referralCodeData, username), nil
}

// Cast repository.ReferralCode to service ReferralCode structure
func castToReferralCode(source repository.ReferralCode, username string) ReferralCode {
	if source.Title.String == "" {
		source.Title = sql.NullString{
			String: username,
			Valid:  true,
		}
	}

	data := ReferralCode{
		ID:           source.ID,
		Title:        source.Title.String,
		Code:         source.Code,
		ReferralLink: source.ReferralLink.String,
		IsPersonal:   source.IsPersonal.Bool,
		CreatedAt:    source.CreatedAt,
	}
	if source.UserID.Valid && source.UserID.UUID != uuid.Nil {
		data.UserID = &source.UserID.UUID
	}

	return data
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

	params := repository.AddReferralCodeDataParams{
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
	}

	if rc.UserID != nil && *rc.UserID != uuid.Nil {
		params.UserID = uuid.NullUUID{UUID: *rc.UserID, Valid: true}
	}

	referralCode, err := s.rr.AddReferralCodeData(ctx, params)
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
func (s *Service) GetReferralCodesDataList(ctx context.Context, limit, offset int32) ([]ReferralCode, int64, error) {
	referralCodes, err := s.rr.GetReferralCodesDataList(ctx, repository.GetReferralCodesDataListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("could not get referral codes list: %w", err)
	}

	numberOfReferralCodes, err := s.rr.GetNumberOfReferralCodes(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("could not get number of referral codes list: %w", err)
	}

	return castToReferralCodes(referralCodes), numberOfReferralCodes, nil
}

// Cast repository.ReferralCode to service ReferralCode structure.
func castToReferralCodes(source []repository.ReferralCode) []ReferralCode {
	result := make([]ReferralCode, 0, len(source))
	for _, s := range source {

		data := ReferralCode{
			ID:         s.ID,
			Title:      s.Title.String,
			Code:       s.Code,
			IsPersonal: s.IsPersonal.Bool,
			CreatedAt:  s.CreatedAt,
		}

		if s.UserID.Valid && s.UserID.UUID != uuid.Nil {
			data.UserID = &s.UserID.UUID
		}

		result = append(result, data)
	}

	return result
}

// UpdateReferralCodeData used to update referral code.
func (s *Service) UpdateReferralCodeData(ctx context.Context, rc ReferralCode) error {
	params := repository.UpdateReferralCodeDataParams{
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
		ID: rc.ID,
	}

	if rc.UserID != nil && *rc.UserID != uuid.Nil {
		params.UserID = uuid.NullUUID{UUID: *rc.UserID, Valid: true}
	}

	if err := s.rr.UpdateReferralCodeData(ctx, params); err != nil {
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
	rc, err := s.rr.GetReferralCodeDataByCode(ctx, code)
	if err != nil {
		if db.IsNotFoundError(err) {
			return false, fmt.Errorf("could not found referral code %s", code)
		}
		return false, fmt.Errorf("could not get referral code %s: %w", code, err)
	}

	if rc.UserID.UUID == uid {
		return false, errors.New("it's your own referral code")
	}

	if err := s.rr.AddReferral(ctx, repository.AddReferralParams{
		ReferralCodeID: rc.ID,
		UserID:         uid,
	}); err != nil {
		return false, fmt.Errorf("could not store referral with id = %v, and code = %s: %w", uid, code, err)
	}

	return true, nil
}
