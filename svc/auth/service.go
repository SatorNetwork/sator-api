package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/SatorNetwork/sator-api/svc/auth/repository"
	repository2 "github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Service struct
	Service struct {
		ur     userRepository
		sc     solanaClient
		ws     walletService
		jwt    jwtInteractor
		mail   mailer
		otpLen int
	}

	// ServiceOption function
	// interface to extend service via options
	ServiceOption func(*Service)

	jwtInteractor interface {
		NewWithUserData(userID uuid.UUID, username string) (uuid.UUID, string, error)
	}

	userRepository interface {
		// user
		CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		GetUserByUsername(ctx context.Context, username string) (repository.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (repository.User, error)
		UpdateUserPassword(ctx context.Context, arg repository.UpdateUserPasswordParams) error
		UpdateUserVerifiedAt(ctx context.Context, verifiedAt sql.NullTime) error

		// email verification
		CreateUserVerification(ctx context.Context, arg repository.CreateUserVerificationParams) error
		GetUserVerificationByUserID(ctx context.Context, userID uuid.UUID) (repository.UserVerification, error)
		DeleteUserVerificationsByUserID(ctx context.Context, userID uuid.UUID) error

		// reset password
		CreatePasswordReset(ctx context.Context, arg repository.CreatePasswordResetParams) error
		GetPasswordResetByEmail(ctx context.Context, email string) (repository.PasswordReset, error)
		DeletePasswordResetsByUserID(ctx context.Context, userID uuid.UUID) error
	}

	mailer interface {
		SendVerificationEmail(ctx context.Context, email, otp string) error
		SendResetPasswordEmail(ctx context.Context, email, otp string) error
	}

	solanaClient interface {
		CreateAccount(ctx context.Context) (string, []byte, error)
	}

	walletService interface {
		CreateWallet(ctx context.Context, userID uuid.UUID, publicKey string, privateKey []byte) (repository2.Wallet, error)
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation.
func NewService(ji jwtInteractor, ur userRepository, sc solanaClient, ws walletService, opt ...ServiceOption) *Service {
	if ur == nil {
		log.Fatalln("user repository is not set")
	}
	if ji == nil {
		log.Fatalln("jwt interactor is not set")
	}
	if sc == nil {
		log.Fatalln("solana client is not set")
	}
	if ws == nil {
		log.Fatalln("wallet service is not set")
	}

	s := &Service{jwt: ji, ur: ur, sc: sc, ws: ws, otpLen: 5}

	// Set up options.
	for _, o := range opt {
		o(s)
	}

	return s
}

// Login by email and password, returns token.
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("could not log in: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	_, token, err := s.jwt.NewWithUserData(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("could not generate new access token: %w", err)
	}

	return token, nil
}

// Logout revokes JWT token.
func (s *Service) Logout(ctx context.Context, tid string) error {
	// TODO: add JWT id into the revoked tokens list
	return nil
}

// RefreshToken returns new jwt string.
func (s *Service) RefreshToken(ctx context.Context, uid uuid.UUID, username, tid string) (string, error) {
	// TODO: add JWT id into the revoked tokens list
	_, token, err := s.jwt.NewWithUserData(uid, username)
	if err != nil {
		return "", fmt.Errorf("could not refresh access token: %w", err)
	}
	return token, nil
}

// SignUp registers account with email, password and username.
func (s *Service) SignUp(ctx context.Context, email, password, username string) (string, error) {
	// Make email address case-insensitive
	email = strings.ToLower(email)

	// Check if the passed email address is not taken yet
	if _, err := s.ur.GetUserByEmail(ctx, email); err == nil {
		return "", validator.NewValidationError(url.Values{
			"email": []string{"email is already taken"},
		})
	} else if !db.IsNotFoundError(err) {
		return "", fmt.Errorf("could not create a new account: %w", err)
	}

	// Check if the passed username is not taken yet
	if _, err := s.ur.GetUserByUsername(ctx, username); err == nil {
		return "", validator.NewValidationError(url.Values{
			"username": []string{"username is already taken"},
		})
	} else if !db.IsNotFoundError(err) {
		return "", fmt.Errorf("could not create a new account: %w", err)
	}

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not create a new account: %w", err)
	}

	u, err := s.ur.CreateUser(ctx, repository.CreateUserParams{
		Email:    email,
		Password: passwdHash,
		Username: username,
	})
	if err != nil {
		return "", fmt.Errorf("could not create a new account: %w", err)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("could not create a new account: %w", err)
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		UserID:           u.ID,
		Email:            email,
		VerificationCode: otpHash,
	}); err != nil {
		return "", fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendVerificationEmail(ctx, email, otp); err != nil {
			return "", fmt.Errorf("could not send verification code: %w", err)
		}
	} else {
		// log data for debug mode
		log.Println("mail service is not set")
		log.Printf("[email verification] email: %s, otp: %s", email, otp)
	}

	_, token, err := s.jwt.NewWithUserData(u.ID, u.Username)
	if err != nil {
		return "", fmt.Errorf("could not generate new access token: %w", err)
	}

	return token, nil
}

// ForgotPassword requests password reset with email.
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	u, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("user %w", ErrNotFound)
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("could not generate a new reset password code: %w", err)
	}

	if err := s.ur.CreatePasswordReset(ctx, repository.CreatePasswordResetParams{
		UserID: u.ID,
		Email:  email,
		Token:  otpHash,
	}); err != nil {
		return fmt.Errorf("could not store reset password code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendResetPasswordEmail(ctx, email, otp); err != nil {
			return fmt.Errorf("could not send reset password code: %w", err)
		}
	} else {
		// log data for debug mode
		log.Println("mail service is not set")
		log.Printf("[reset password] email: %s, otp: %s", email, otp)
	}

	return nil
}

// ValidateResetPasswordCode validates reset password code,
// it's needed to implement the reset password flow on the mobile client.
func (s *Service) ValidateResetPasswordCode(ctx context.Context, email, otp string) (uuid.UUID, error) {
	pr, err := s.ur.GetPasswordResetByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return uuid.Nil, fmt.Errorf("%w user with given email address", ErrNotFound)
		}
		return uuid.Nil, fmt.Errorf("could not get user with given email address: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(pr.Token, []byte(otp)); err != nil {
		return uuid.Nil, ErrOTPCode
	}

	return pr.UserID, nil
}

// ResetPassword changing password.
func (s *Service) ResetPassword(ctx context.Context, email, password, otp string) error {
	userID, err := s.ValidateResetPasswordCode(ctx, email, otp)
	if err != nil {
		return fmt.Errorf("could not reset password: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return fmt.Errorf("could not reset password: %w", err)
	}

	if err := s.ur.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:       userID,
		Password: passwordHash,
	}); err != nil {
		return fmt.Errorf("could not reset password: %w", err)
	}

	if err := s.ur.DeletePasswordResetsByUserID(ctx, userID); err != nil {
		// just log, not any error for user
		log.Printf("could not delete password resets for user with id=%s: %v", userID.String(), err)
	}

	return nil
}

// VerifyAccount verifies account.
func (s *Service) VerifyAccount(ctx context.Context, userID uuid.UUID, otp string) error {
	u, err := s.ur.GetUserByID(ctx, userID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("user %w", ErrNotFound)
		}
		return fmt.Errorf("could not get user: %w", err)
	}
	if u.VerifiedAt.Valid {
		return ErrEmailAlreadyVerified
	}

	uv, err := s.ur.GetUserVerificationByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not verify email address: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(uv.VerificationCode, []byte(otp)); err != nil {
		return ErrOTPCode
	}

	if err := s.ur.UpdateUserVerifiedAt(ctx, sql.NullTime{Time: time.Now(), Valid: true}); err != nil {
		return fmt.Errorf("could not verify email address: %w", err)
	}

	publicKey, privateKey, err := s.sc.CreateAccount(ctx)
	if err != nil {
		return fmt.Errorf("counld not create solana account: %w", err)
	}

	_, err = s.ws.CreateWallet(ctx, u.ID, publicKey, privateKey)
	if err != nil {
		return fmt.Errorf("counld not create solana wallet: %w", err)
	}

	if err := s.ur.DeleteUserVerificationsByUserID(ctx, userID); err != nil {
		// just log, not any error for user
		log.Printf("could not delete verification code for user with id=%s: %v", userID.String(), err)
	}

	return nil
}
