package mail

import (
	"context"

	"github.com/keighl/postmark"
)

// Predefined email templates
var (
	VerificationCodeTmpl   = "verification_code"
	PasswordResetTmpl      = "password_reset"
	DestroyAccountCodeTmpl = "destroy_account"
	InvitationCodeTmpl     = "invitation"
)

type (
	// Service struct
	Service struct {
		client postmarkClient
		config Config
	}

	// Config struct
	Config struct {
		ProductName    string
		ProductURL     string
		SupportURL     string
		SupportEmail   string
		CompanyName    string
		CompanyAddress string
		FromEmail      string
		FromName       string
	}

	postmarkClient interface {
		SendTemplatedEmail(email postmark.TemplatedEmail) (postmark.EmailResponse, error)
	}
)

// New mail service
func New(pc postmarkClient, cnf Config) *Service {
	return &Service{client: pc, config: cnf}
}

// SendVerificationCode ...
func (s *Service) SendVerificationCode(_ context.Context, email, otp string) error {
	return nil
}

// SendResetPasswordCode ...
func (s *Service) SendResetPasswordCode(_ context.Context, email, otp string) error {
	return nil
}

// SendDestroyAccountCode ...
func (s *Service) SendDestroyAccountCode(_ context.Context, email, otp string) error {
	return nil
}

// SendInvitation ...
func (s *Service) SendInvitation(_ context.Context, email, invitedBy string) error {
	return nil
}
