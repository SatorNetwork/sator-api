package invitations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SatorNetwork/sator-api/svc/invitations/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		ir invitationsRepository
		m  mailer
	}

	// Invitation struct
	// Fields were rearranged to optimize memory usage.
	Invitation struct {
		ID                     uuid.UUID `json:"id"`
		InviteeEmail           string    `json:"invitee_email"`
		NormalizedInviteeEmail string    `json:"normalized_invitee_email"`
		InvitedAt              time.Time `json:"invited_at"`
		InvitedBy              uuid.UUID `json:"invited_by"`
		AcceptedAt             time.Time `json:"accepted_at"`
		AcceptedBy             uuid.UUID `json:"accepted_by"`
	}

	invitationsRepository interface {
		AcceptInvitationByInviteeEmail(ctx context.Context, arg repository.AcceptInvitationByInviteeEmailParams) error
		CreateInvitation(ctx context.Context, arg repository.CreateInvitationParams) (repository.Invitation, error)
		GetInvitations(ctx context.Context) ([]repository.Invitation, error)
		GetInvitationsPaginated(ctx context.Context, arg repository.GetInvitationsPaginatedParams) ([]repository.Invitation, error)
		GetInvitationByInviteeEmail(ctx context.Context, normalizedInviteeEmail string) (repository.Invitation, error)
		GetInvitationsByInvitedByID(ctx context.Context, invitedBy uuid.UUID) ([]repository.Invitation, error)
	}

	mailer interface {
		SendInvitationCode(ctx context.Context, email, otp string) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(ir invitationsRepository, m mailer) *Service {
	if ir == nil {
		log.Fatalln("invitations repository is not set")
	}
	if m == nil {
		log.Fatalln("mailer client is not set")
	}

	return &Service{ir: ir, m: m}
}

// SendInvitation used to send invitation if person doesn't exist in invitation table.
func (s *Service) SendInvitation(ctx context.Context, invitedByID uuid.UUID, inviteeEmail string) error {
	normalizedEmail := strings.ToUpper(inviteeEmail)

	invitation, err := s.ir.GetInvitationByInviteeEmail(ctx, normalizedEmail)
	switch {
	case err == sql.ErrNoRows:
		invitation, err := s.ir.CreateInvitation(ctx, repository.CreateInvitationParams{
			InviteeEmail:           inviteeEmail,
			NormalizedInviteeEmail: normalizedEmail,
			InvitedBy:              invitedByID,
		})
		if err != nil {
			return fmt.Errorf("could not create invitation: %w", err)
		}

		err = s.m.SendInvitationCode(ctx, invitation.InviteeEmail, "STRING?")
		if err != nil {

			return fmt.Errorf("could not send invitation: %w", err)
		}

		return nil
	case err != nil:
		return fmt.Errorf("could not get invitation by invitee email: %w", err)
	default:
		return fmt.Errorf("user with email = %s, alredy invited: %w", invitation.InviteeEmail, err)
	}
}

// GetInvitations returns invitations.
func (s *Service) GetInvitations(ctx context.Context, limit, offset int32) ([]Invitation, error) {
	invitations, err := s.ir.GetInvitationsPaginated(ctx, repository.GetInvitationsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get invitations list: %w", err)
	}

	return castToListInvitations(invitations), nil
}

// Cast repository.Show to service Show structure
func castToListInvitations(source []repository.Invitation) []Invitation {
	result := make([]Invitation, 0, len(source))
	for _, s := range source {
		result = append(result, Invitation{
			ID:                     s.ID,
			InviteeEmail:           s.InviteeEmail,
			NormalizedInviteeEmail: s.NormalizedInviteeEmail,
			InvitedAt:              s.InvitedAt,
			InvitedBy:              s.InvitedBy,
			AcceptedAt:             s.AcceptedAt.Time,
			AcceptedBy:             s.AcceptedBy,
		})
	}

	return result
}
