package invitations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/invitations/repository"
	"github.com/SatorNetwork/sator-api/svc/rewards"

	"github.com/google/uuid"
)

const (
	// RelationTypeInvitation indicates that relation type is "invitation".
	RelationTypeInvitation = "invitation"
)

type (
	// Service struct
	Service struct {
		ir     invitationsRepository
		m      mailer
		rc     rewardsClient
		config Config
	}

	// Config struct
	Config struct {
		InvitationReward float64
		InvitationURL    string
	}

	// Invitation struct
	// Fields were rearranged to optimize memory usage.
	Invitation struct {
		ID             uuid.UUID `json:"id"`
		Email          string    `json:"email"`
		InvitedBy      uuid.UUID `json:"invited_by"`
		InvitedAt      time.Time `json:"invited_at"`
		AcceptedBy     uuid.UUID `json:"accepted_by"`
		AcceptedAt     time.Time `json:"accepted_at"`
		RewardReceived bool      `json:"reward_received"`
	}

	invitationsRepository interface {
		AcceptInvitationByInviteeEmail(ctx context.Context, arg repository.AcceptInvitationByInviteeEmailParams) error
		CreateInvitation(ctx context.Context, arg repository.CreateInvitationParams) (repository.Invitation, error)
		GetInvitations(ctx context.Context) ([]repository.Invitation, error)
		GetInvitationsPaginated(ctx context.Context, arg repository.GetInvitationsPaginatedParams) ([]repository.Invitation, error)
		GetInvitationByInviteeEmail(ctx context.Context, normalizedInviteeEmail string) (repository.Invitation, error)
		GetInvitationByInviteeID(ctx context.Context, acceptedBy uuid.UUID) (repository.Invitation, error)
		GetInvitationsByInviterID(ctx context.Context, invitedBy uuid.UUID) ([]repository.Invitation, error)
		SetRewardReceived(ctx context.Context, arg repository.SetRewardReceivedParams) error
	}

	mailer interface {
		SendInvitation(_ context.Context, email, invitedBy string) error
	}

	rewardsClient interface {
		AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(ir invitationsRepository, m mailer, rc rewardsClient, config Config) *Service {
	if ir == nil {
		log.Fatalln("invitations repository is not set")
	}
	if m == nil {
		log.Fatalln("mailer client is not set")
	}
	if rc == nil {
		log.Fatalln("rewards client is not set")
	}

	return &Service{ir: ir, m: m, rc: rc, config: config}
}

// SendReward ...
func (s *Service) SendReward(sendRewards func(ctx context.Context, uid, relationID uuid.UUID, relationType string, amount float64, trType int32) error) func(userID, quizID uuid.UUID) {
	return func(userID, quizID uuid.UUID) {
		var ctx context.Context

		invitation, err := s.ir.GetInvitationByInviteeID(ctx, userID)
		if err != nil {
			if !db.IsNotFoundError(err) {
				return // user isn't invited.
			}
		}
		if invitation.RewardReceived.Bool {
			return // reward received.
		}

		// sendRewards
		if err := sendRewards(
			ctx,
			invitation.InvitedBy,
			quizID,
			RelationTypeInvitation,
			s.config.InvitationReward,
			rewards.TransactionTypeDeposit,
		); err != nil {
			log.Printf("could not send invitation reward: %v", err)
			return
		}

		if err := s.ir.SetRewardReceived(ctx, repository.SetRewardReceivedParams{
			ID: invitation.ID,
			RewardReceived: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		}); err != nil {
			log.Printf("could not set invitation reward received: %v", err)
			return
		}
	}
}

// SendInvitation used to send invitation if person doesn't exist in invitation table.
func (s *Service) SendInvitation(ctx context.Context, invitedByID uuid.UUID, invitedByUsername, inviteeEmail string) error {
	if invitation, err := s.ir.GetInvitationByInviteeEmail(ctx, inviteeEmail); err != nil {
		if db.IsNotFoundError(err) {
			invitation, err = s.ir.CreateInvitation(ctx, repository.CreateInvitationParams{
				Email:     inviteeEmail,
				InvitedBy: invitedByID,
			})
			if err != nil {
				return fmt.Errorf("could not create invitation: %w", err)
			}

			err = s.m.SendInvitation(ctx, invitation.Email, invitedByUsername)
			if err != nil {

				return fmt.Errorf("could not send invitation: %w", err)
			}
		}

		return fmt.Errorf("could not get invitation by invitee email: %w", err)
	}

	return fmt.Errorf("user with email %s, alredy invited", inviteeEmail)
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
			ID:             s.ID,
			Email:          s.Email,
			InvitedAt:      s.InvitedAt,
			InvitedBy:      s.InvitedBy,
			AcceptedAt:     s.AcceptedAt.Time,
			AcceptedBy:     s.AcceptedBy,
			RewardReceived: s.RewardReceived.Bool,
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
	invitation, err := s.ir.GetInvitationByInviteeEmail(ctx, inviteeEmail)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("could not find invitation by email %s", inviteeEmail)
		}
		return fmt.Errorf("could not get invitation by email %s: %w", inviteeEmail, err)
	}

	if invitation.AcceptedAt.Valid {
		return fmt.Errorf("invitation is already accepted")
	}

	if err := s.ir.AcceptInvitationByInviteeEmail(ctx, repository.AcceptInvitationByInviteeEmailParams{
		ID:         invitation.ID,
		AcceptedBy: inviteeID,
		AcceptedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("could not accept invitation for user id = %s: %w", inviteeID, err)
	}

	return nil
}

// IsEmailInvited returns true if email invited, false if not.
func (s *Service) IsEmailInvited(ctx context.Context, inviteeEmail string) (bool, error) {
	if _, err := s.ir.GetInvitationByInviteeEmail(ctx, inviteeEmail); err != nil {
		if !db.IsNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("could not get invitation data for user with email = %s: %w", inviteeEmail, err)
	}

	return true, nil
}
