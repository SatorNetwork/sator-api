package invitations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/invitations/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		ir invitationsRepository
		m  mailer
		rc rewardsClient
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
		GetInvitationByInviteeID(ctx context.Context, acceptedBy uuid.UUID) (repository.Invitation, error)
		GetInvitationsByInvitedByID(ctx context.Context, invitedBy uuid.UUID) ([]repository.Invitation, error)
	}

	mailer interface {
		SendInvitationCode(ctx context.Context, email, otp string) error
	}

	rewardsClient interface {
		AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(ir invitationsRepository, m mailer, rc rewardsClient) *Service {
	if ir == nil {
		log.Fatalln("invitations repository is not set")
	}
	if m == nil {
		log.Fatalln("mailer client is not set")
	}
	if rc == nil {
		log.Fatalln("rewards client is not set")
	}

	return &Service{ir: ir, m: m, rc: rc}
}

// SendReward ...
func (s *Service) SendReward(sendRewards func(ctx context.Context, uid, relationID uuid.UUID, relationType string, amount float64, trType int32) error) func(userID, quizID uuid.UUID) {
	return func(userID, quizID uuid.UUID) {
		// invited?
		//_, err := s.ir.GetInvitationByInviteeID(ctx, userID)
		//if err != nil {
		//	if !db.IsNotFoundError(err) {
		//
		//	}
		//}

		// TODO: add check is get reward per referral

		// sendRewards

	}
}

// SendInvitation used to send invitation if person doesn't exist in invitation table.
func (s *Service) SendInvitation(ctx context.Context, invitedByID uuid.UUID, inviteeEmail string) error {
	normalizedEmail := strings.ToUpper(inviteeEmail)

	invitation, err := s.ir.GetInvitationByInviteeEmail(ctx, normalizedEmail)
	if err != nil {
		if !db.IsNotFoundError(err) {
			inv, err := s.ir.CreateInvitation(ctx, repository.CreateInvitationParams{
				InviteeEmail:           inviteeEmail,
				NormalizedInviteeEmail: normalizedEmail,
				InvitedBy:              invitedByID,
			})
			if err != nil {
				return fmt.Errorf("could not create invitation: %w", err)
			}

			err = s.m.SendInvitationCode(ctx, inv.InviteeEmail, "STRING?")
			if err != nil {

				return fmt.Errorf("could not send invitation: %w", err)
			}
		}

		return fmt.Errorf("could not get invitation by invitee email: %w", err)
	}

	return fmt.Errorf("user with email = %s, alredy invited: %w", invitation.InviteeEmail, err)
}

// GetInvitationsPaginated returns list invitations with pagination.
func (s *Service) GetInvitationsPaginated(ctx context.Context, limit, offset int32) ([]Invitation, error) {
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

// GetInvitations returns list invitations.
func (s *Service) GetInvitations(ctx context.Context) ([]Invitation, error) {
	invitations, err := s.ir.GetInvitations(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get invitations list: %w", err)
	}

	return castToListInvitations(invitations), nil
}

// AcceptInvitation used to accept invitation and store invitee ID and email.
func (s *Service) AcceptInvitation(ctx context.Context, inviteeID uuid.UUID, inviteeEmail string) error {
	err := s.ir.AcceptInvitationByInviteeEmail(ctx, repository.AcceptInvitationByInviteeEmailParams{
		AcceptedBy: inviteeID,
		AcceptedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		InviteeEmail: inviteeEmail,
	})
	if err != nil {
		return fmt.Errorf("could not accept invitation for user id = %s: %w", inviteeID, err)
	}

	return nil
}

// IsEmailInvited returns true if email invited, false if not.
func (s *Service) IsEmailInvited(ctx context.Context, inviteeEmail string) (bool, error) {
	_, err := s.ir.GetInvitationByInviteeEmail(ctx, inviteeEmail)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("could not get invitation data for user with email = %s: %w", inviteeEmail, err)
	}

	return true, nil
}
