package auth

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	mail2 "github.com/SatorNetwork/sator-api/internal/mail"
	"github.com/SatorNetwork/sator-api/svc/auth/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		ur   userRepository
		jwt  jwtInteractor
		mail *mail2.SMTPSender
	}

	jwtInteractor interface {
		NewWithUserID(userID uuid.UUID) (uuid.UUID, string, error)
	}

	userRepository interface {
		CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
		CreateUserVerification(ctx context.Context, arg repository.CreateUserVerificationParams) error
		CreatePasswordReset(ctx context.Context, arg repository.CreatePasswordResetParams) error
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		UpdateUserPassword(ctx context.Context, arg repository.UpdateUserPasswordParams) error
		UpdateUserVerifiedAt(ctx context.Context, verifiedAt sql.NullTime) error
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation
func NewService(ji jwtInteractor, ur userRepository) *Service {
	if ur == nil {
		log.Fatalln("user repository is not set")
	}
	if ji == nil {
		log.Fatalln("jwt interactor is not set")
	}
	return &Service{jwt: ji, ur: ur}
}

// Login by email and password, returns token
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil{
		return "", errors.New("wrong password")
	}

	_, token, err := s.jwt.NewWithUserID(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Logout logging out from account.
func (s *Service) Logout(ctx context.Context) error {
	return nil
}

// SignUp registers account with email password and username.
func (s *Service) SignUp(ctx context.Context, email, password, username string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	_, err = s.ur.CreateUser(ctx, repository.CreateUserParams{Email: email, Password: bytes, ID: uuid.New(), Username: username})
	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword requests password reset with email.
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	err := s.mail.SendEmail(&mail2.Message{
		From:      mail2.Address{},
		To:        []mail2.Address{{Address: email}},
		Subject:   "",
		ID:        "",
		Date:      time.Time{},
		ReceiptTo: nil,
		PlainText: "",
		Parts:     nil,
	})
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword changing password.
func (s *Service) ResetPassword(ctx context.Context, email, password, otp string) error {
	token := []byte(otp)

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = s.ur.CreatePasswordReset(ctx, repository.CreatePasswordResetParams{Email: email, UserID: user.ID, Token: token})
	if err != nil {
		return err
	}

	err = s.ur.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{ID: user.ID, Password: bytes})
	if err != nil {
		return err
	}

	return nil
}

// VerifyAccount verifies account.
func (s *Service) VerifyAccount(ctx context.Context) error {
	return s.ur.UpdateUserVerifiedAt(ctx, sql.NullTime{Time: time.Now(), Valid: true})
}
