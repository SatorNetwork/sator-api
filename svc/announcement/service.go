package announcement

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	announcement_repository "github.com/SatorNetwork/sator-api/svc/announcement/repository"
)

type (
	Service struct {
		ar announcementRepository
	}

	announcementRepository interface {
		CreateAnnouncement(
			ctx context.Context,
			arg announcement_repository.CreateAnnouncementParams,
		) (announcement_repository.Announcement, error)
		GetAnnouncementByID(ctx context.Context, id uuid.UUID) (announcement_repository.Announcement, error)
		UpdateAnnouncementByID(ctx context.Context, arg announcement_repository.UpdateAnnouncementByIDParams) error
		DeleteAnnouncementByID(ctx context.Context, id uuid.UUID) error
		ListAnnouncements(ctx context.Context) ([]announcement_repository.Announcement, error)
		ListUnreadAnnouncements(ctx context.Context, userID uuid.UUID) ([]announcement_repository.Announcement, error)
		ListActiveAnnouncements(ctx context.Context) ([]announcement_repository.Announcement, error)

		MarkAsRead(ctx context.Context, arg announcement_repository.MarkAsReadParams) error
		IsRead(ctx context.Context, arg announcement_repository.IsReadParams) (bool, error)
		IsNotRead(ctx context.Context, arg announcement_repository.IsNotReadParams) (bool, error)
	}

	Empty struct{}

	CreateAnnouncementRequest struct {
		Title              string            `json:"title"`
		Description        string            `json:"description"`
		ActionUrl          string            `json:"action_url"`
		StartsAt           int64             `json:"starts_at"`
		EndsAt             int64             `json:"ends_at"`
		Type               string            `json:"type"`
		TypeSpecificParams map[string]string `json:"type_specific_params"`
	}

	CreateAnnouncementResponse struct {
		ID string `json:"id"`
	}

	GetAnnouncementByIDRequest struct {
		ID string `json:"id"`
	}

	Announcement struct {
		ID                 string            `json:"id"`
		Title              string            `json:"title"`
		Description        string            `json:"description"`
		ActionUrl          string            `json:"action_url"`
		StartsAt           int64             `json:"starts_at"`
		EndsAt             int64             `json:"ends_at"`
		Type               string            `json:"type"`
		TypeSpecificParams map[string]string `json:"type_specific_params"`
	}

	UpdateAnnouncementRequest struct {
		ID                 string            `json:"id"`
		Title              string            `json:"title"`
		Description        string            `json:"description"`
		ActionUrl          string            `json:"action_url"`
		StartsAt           int64             `json:"starts_at"`
		EndsAt             int64             `json:"ends_at"`
		Type               string            `json:"type"`
		TypeSpecificParams map[string]string `json:"type_specific_params"`
	}

	DeleteAnnouncementRequest struct {
		ID string `json:"id"`
	}

	MarkAsReadRequest struct {
		AnnouncementID string `json:"announcement_id"`
	}

	GetAnnouncementTypesResponse struct {
		Types []string `json:"types"`
	}
)

func NewService(
	ar announcementRepository,
) *Service {
	s := &Service{
		ar: ar,
	}

	return s
}

func (s *Service) CreateAnnouncement(ctx context.Context, req *CreateAnnouncementRequest) (*CreateAnnouncementResponse, error) {
	startsAt := time.Unix(req.StartsAt, 0).UTC()
	endsAt := time.Unix(req.EndsAt, 0).UTC()

	typeSpecificParamsInJSON, err := json.Marshal(req.TypeSpecificParams)
	if err != nil {
		return nil, err
	}
	announcement, err := s.ar.CreateAnnouncement(ctx, announcement_repository.CreateAnnouncementParams{
		Title:              req.Title,
		Description:        req.Description,
		ActionUrl:          req.ActionUrl,
		StartsAt:           startsAt,
		EndsAt:             endsAt,
		Type:               req.Type,
		TypeSpecificParams: string(typeSpecificParamsInJSON),
	})
	if err != nil {
		return nil, err
	}

	return &CreateAnnouncementResponse{
		ID: announcement.ID.String(),
	}, nil
}

func (s *Service) GetAnnouncementByID(ctx context.Context, req *GetAnnouncementByIDRequest) (*Announcement, error) {
	uid, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}
	a, err := s.ar.GetAnnouncementByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return NewAnnouncementFromSQLC(&a)
}

func (s *Service) UpdateAnnouncementByID(ctx context.Context, req *UpdateAnnouncementRequest) error {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}
	startsAt := time.Unix(req.StartsAt, 0).UTC()
	endsAt := time.Unix(req.EndsAt, 0).UTC()

	typeSpecificParamsInJSON, err := json.Marshal(req.TypeSpecificParams)
	if err != nil {
		return err
	}
	err = s.ar.UpdateAnnouncementByID(ctx, announcement_repository.UpdateAnnouncementByIDParams{
		ID:                 id,
		Title:              req.Title,
		Description:        req.Description,
		ActionUrl:          req.ActionUrl,
		StartsAt:           startsAt,
		EndsAt:             endsAt,
		Type:               req.Type,
		TypeSpecificParams: string(typeSpecificParamsInJSON),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteAnnouncementByID(ctx context.Context, req *DeleteAnnouncementRequest) error {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return err
	}

	err = s.ar.DeleteAnnouncementByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func NewAnnouncementFromSQLC(a *announcement_repository.Announcement) (*Announcement, error) {
	var typeSpecificParams map[string]string
	err := json.Unmarshal([]byte(a.TypeSpecificParams), &typeSpecificParams)
	if err != nil {
		return nil, err
	}
	return &Announcement{
		ID:                 a.ID.String(),
		Title:              a.Title,
		Description:        a.Description,
		ActionUrl:          a.ActionUrl,
		StartsAt:           a.StartsAt.Unix(),
		EndsAt:             a.EndsAt.Unix(),
		Type:               a.Type,
		TypeSpecificParams: typeSpecificParams,
	}, nil
}

func NewAnnouncementsFromSQLC(sqlcAnnouncements []announcement_repository.Announcement) ([]*Announcement, error) {
	announcements := make([]*Announcement, 0, len(sqlcAnnouncements))
	for _, sqlcA := range sqlcAnnouncements {
		a, err := NewAnnouncementFromSQLC(&sqlcA)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}

	return announcements, nil
}

func (s *Service) ListAnnouncements(ctx context.Context) ([]*Announcement, error) {
	announcements, err := s.ar.ListAnnouncements(ctx)
	if err != nil {
		return nil, err
	}

	return NewAnnouncementsFromSQLC(announcements)
}

func (s *Service) ListUnreadAnnouncements(ctx context.Context, userID uuid.UUID) ([]*Announcement, error) {
	announcements, err := s.ar.ListUnreadAnnouncements(ctx, userID)
	if err != nil {
		return nil, err
	}

	return NewAnnouncementsFromSQLC(announcements)
}

func (s *Service) ListActiveAnnouncements(ctx context.Context) ([]*Announcement, error) {
	announcements, err := s.ar.ListActiveAnnouncements(ctx)
	if err != nil {
		return nil, err
	}

	return NewAnnouncementsFromSQLC(announcements)
}

func (s *Service) MarkAsRead(ctx context.Context, userID uuid.UUID, req *MarkAsReadRequest) error {
	announcementID, err := uuid.Parse(req.AnnouncementID)
	if err != nil {
		return errors.Wrap(err, "can't parse announcement ID")
	}

	err = s.ar.MarkAsRead(ctx, announcement_repository.MarkAsReadParams{
		AnnouncementID: announcementID,
		UserID:         userID,
	})
	if err != nil {
		return errors.Wrap(err, "can't mark announcement as read")
	}

	return nil
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	announcements, err := s.ar.ListUnreadAnnouncements(ctx, userID)
	if err != nil {
		return err
	}

	for _, a := range announcements {
		err = s.ar.MarkAsRead(ctx, announcement_repository.MarkAsReadParams{
			AnnouncementID: a.ID,
			UserID:         userID,
		})
		if err != nil {
			return errors.Wrap(err, "can't mark announcement as read")
		}
	}

	return nil
}

func (s *Service) GetAnnouncementTypes(ctx context.Context) (*GetAnnouncementTypesResponse, error) {
	return &GetAnnouncementTypesResponse{
		Types: []string{"show", "episode", "link"},
	}, nil
}
