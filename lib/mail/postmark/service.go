//go:build !mock_postmark

package postmark

import (
	"context"
	"fmt"

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
)

// New mail service
func New(pc postmarkClient, cnf Config) *Service {
	return &Service{client: pc, config: cnf}
}

// SendVerificationCode ...
func (s *Service) SendVerificationCode(_ context.Context, email, otp string) error {
	if err := s.send(VerificationCodeTmpl, "verification", email, map[string]interface{}{
		"otp": otp,
	}); err != nil {
		return fmt.Errorf("could not send verification code: %w", err)
	}
	return nil
}

// SendResetPasswordCode ...
func (s *Service) SendResetPasswordCode(_ context.Context, email, otp string) error {
	if err := s.send(PasswordResetTmpl, "reset_password", email, map[string]interface{}{
		"otp": otp,
	}); err != nil {
		return fmt.Errorf("could not send reset password code: %w", err)
	}
	return nil
}

// SendDestroyAccountCode ...
func (s *Service) SendDestroyAccountCode(_ context.Context, email, otp string) error {
	if err := s.send(DestroyAccountCodeTmpl, "destroy_account", email, map[string]interface{}{
		"otp": otp,
	}); err != nil {
		return fmt.Errorf("could not send verification code: %w", err)
	}
	return nil
}

// SendInvitation ...
func (s *Service) SendInvitation(_ context.Context, email, invitedBy string) error {
	if err := s.send(InvitationCodeTmpl, "invitation", email, map[string]interface{}{
		"invited_by": invitedBy,
	}); err != nil {
		return fmt.Errorf("could not send invitation to email %s: %w", email, err)
	}
	return nil
}

// send email
func (s *Service) send(tpl, tag, email string, data map[string]interface{}) error {
	// Default model data
	payload := map[string]interface{}{
		"product_url":     s.config.ProductURL,
		"product_name":    s.config.ProductName,
		"support_url":     s.config.SupportURL,
		"company_name":    s.config.CompanyName,
		"company_address": s.config.CompanyAddress,
		"email":           email,
	}

	// Merge custom data with default fields
	for k, v := range data {
		payload[k] = v
	}

	if _, err := s.client.SendTemplatedEmail(postmark.TemplatedEmail{
		TemplateAlias: tpl,
		InlineCss:     true,
		TrackOpens:    true,
		From:          s.config.FromEmail,
		To:            email,
		Tag:           tag,
		ReplyTo:       s.config.SupportEmail,
		TemplateModel: payload,
	}); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}
