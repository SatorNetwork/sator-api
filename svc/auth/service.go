package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/SatorNetwork/sator-api/svc/auth/repository"

	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Service struct
	Service struct {
		ur         userRepository
		ws         walletService
		jwt        jwtInteractor
		mail       mailer
		ic         invitationsClient
		otpLen     int
		masterCode string
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
		UpdateUserEmail(ctx context.Context, arg repository.UpdateUserEmailParams) error
		UpdateUsername(ctx context.Context, arg repository.UpdateUsernameParams) error
		UpdateUserPassword(ctx context.Context, arg repository.UpdateUserPasswordParams) error
		UpdateUserVerifiedAt(ctx context.Context, arg repository.UpdateUserVerifiedAtParams) error
		DestroyUser(ctx context.Context, id uuid.UUID) error

		// email verification
		CreateUserVerification(ctx context.Context, arg repository.CreateUserVerificationParams) error
		GetUserVerificationByUserID(ctx context.Context, arg repository.GetUserVerificationByUserIDParams) (repository.UserVerification, error)
		GetUserVerificationByEmail(ctx context.Context, arg repository.GetUserVerificationByEmailParams) (repository.UserVerification, error)
		DeleteUserVerificationsByUserID(ctx context.Context, arg repository.DeleteUserVerificationsByUserIDParams) error
	}

	mailer interface {
		SendVerificationCode(ctx context.Context, email, otp string) error
		SendResetPasswordCode(ctx context.Context, email, otp string) error
		SendDestroyAccountCode(ctx context.Context, email, otp string) error
	}

	walletService interface {
		CreateWallet(ctx context.Context, userID uuid.UUID) error
	}

	invitationsClient interface {
		AcceptInvitation(ctx context.Context, inviteeID uuid.UUID, inviteeEmail string) error
		IsEmailInvited(ctx context.Context, inviteeEmail string) (bool, error)
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation.
func NewService(ji jwtInteractor, ur userRepository, ws walletService, ic invitationsClient, opt ...ServiceOption) *Service {
	if ur == nil {
		log.Fatalln("user repository is not set")
	}
	if ji == nil {
		log.Fatalln("jwt interactor is not set")
	}
	if ws == nil {
		log.Fatalln("wallet service is not set")
	}
	if ic == nil {
		log.Fatalln("invitations client is not set")
	}

	s := &Service{jwt: ji, ur: ur, ic: ic, ws: ws, otpLen: 5}

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
	u, err := s.ur.GetUserByID(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("could not refresh access token: %w", err)
	}

	if u.Disabled {
		return "", ErrUserIsDisabled
	}

	// TODO: add JWT id into the revoked tokens list
	_, token, err := s.jwt.NewWithUserData(uid, u.Username)
	if err != nil {
		return "", fmt.Errorf("could not refresh access token: %w", err)
	}

	return token, nil
}

// SignUp registers account with email, password and username.
func (s *Service) SignUp(ctx context.Context, email, password, username string) (string, error) {
	var otpHash []byte
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
	if s.mail == nil {
		otpHash = []byte(s.masterCode)
	} else {
		otpHash, err = bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
		if err != nil {
			return "", fmt.Errorf("could not create a new account: %w", err)
		}
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.VerifyConfirmAccount,
		UserID:           u.ID,
		Email:            email,
		VerificationCode: otpHash,
	}); err != nil {
		return "", fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendVerificationCode(ctx, email, otp); err != nil {
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

	if isInvited, _ := s.ic.IsEmailInvited(ctx, email); isInvited {
		if err := s.ic.AcceptInvitation(ctx, u.ID, email); err != nil {
			log.Printf("could not accept invitation for user id = %s: %v", u.ID, err)
		}
	}

	return token, nil
}

// ForgotPassword requests password reset with email.
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	var otpHash []byte

	u, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return errors.New("please check if you set correct email")
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	if s.mail == nil {
		otpHash = []byte(s.masterCode)
	} else {
		otpHash, err = bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
		if err != nil {
			return fmt.Errorf("could not process forgot password: %w", err)
		}
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.VerifyResetPassword,
		UserID:           u.ID,
		Email:            email,
		VerificationCode: otpHash,
	}); err != nil {
		return fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendResetPasswordCode(ctx, email, otp); err != nil {
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
// it's needed to implement the reset password flow on the client.
func (s *Service) ValidateResetPasswordCode(ctx context.Context, email, otp string) (uuid.UUID, error) {
	v, err := s.ur.GetUserVerificationByEmail(ctx, repository.GetUserVerificationByEmailParams{
		RequestType: repository.VerifyResetPassword,
		Email:       email,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return uuid.Nil, fmt.Errorf("%w user with given email address", ErrNotFound)
		}
		return uuid.Nil, fmt.Errorf("could not get user with given email address: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(v.VerificationCode, []byte(otp))
	if err != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(s.masterCode), []byte(otp)); err != nil {
			return uuid.Nil, ErrOTPCode
		}
	}

	return v.UserID, nil
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

	if err := s.ur.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.VerifyResetPassword,
		UserID:      userID,
	}); err != nil {
		// just log, not any error for user
		log.Printf("could not delete password resets for user with id=%s: %v", userID.String(), err)
	}

	return nil
}

// ChangePassword changing password.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.ur.GetUserByID(ctx, userID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("could not log in: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return fmt.Errorf("could not reset password: %w", err)
	}

	if err := s.ur.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:       userID,
		Password: passwordHash,
	}); err != nil {
		return fmt.Errorf("could not reset password: %w", err)
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

	uv, err := s.ur.GetUserVerificationByUserID(ctx, repository.GetUserVerificationByUserIDParams{
		RequestType: repository.VerifyConfirmAccount,
		UserID:      userID,
	})
	if err != nil {
		return fmt.Errorf("could not verify email address: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(uv.VerificationCode, []byte(otp))
	if err != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(s.masterCode), []byte(otp)); err != nil {
			return ErrOTPCode
		}
	}

	if err := s.ur.UpdateUserVerifiedAt(ctx, repository.UpdateUserVerifiedAtParams{
		UserID:     userID,
		VerifiedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		return fmt.Errorf("could not verify email address: %w", err)
	}

	if err := s.ws.CreateWallet(ctx, userID); err != nil {
		return fmt.Errorf("could not create solana wallet: %w", err)
	}

	if err := s.ur.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.VerifyConfirmAccount,
		UserID:      userID,
	}); err != nil {
		// just log, not any error for user
		log.Printf("could not delete verification code for user with id=%s: %v", userID.String(), err)
	}

	return nil
}

// RequestChangeEmail requests change email for authorized user. Checks if email already exists in db, creates user verification and sends code to new email.
func (s *Service) RequestChangeEmail(ctx context.Context, userID uuid.UUID, email string) error {
	var otpHash []byte

	u, err := s.ur.GetUserByID(ctx, userID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("user %w", ErrNotFound)
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	if _, err := s.ur.GetUserByEmail(ctx, email); err == nil {
		return fmt.Errorf("could not update email: %w", ErrEmailAlreadyTaken)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	if s.mail == nil {
		otpHash = []byte(s.masterCode)
	} else {
		otpHash, err = bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
		if err != nil {
			return fmt.Errorf("could not request email change: %w", err)
		}
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.VerifyChangeEmail,
		UserID:           u.ID,
		Email:            email,
		VerificationCode: otpHash,
	}); err != nil {
		return fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendVerificationCode(ctx, email, otp); err != nil {
			return fmt.Errorf("could not send verification code: %w", err)
		}
	} else {
		// log data for debug mode
		log.Println("mail service is not set")
		log.Printf("[reset password] email: %s, otp: %s", email, otp)
	}

	return nil
}

// ValidateChangeEmailCode validates change email code,
// it's needed to implement the reset password flow on the client.
func (s *Service) ValidateChangeEmailCode(ctx context.Context, userID uuid.UUID, email, otp string) error {
	v, err := s.ur.GetUserVerificationByEmail(ctx, repository.GetUserVerificationByEmailParams{
		RequestType: repository.VerifyChangeEmail,
		Email:       email,
	})
	if err != nil || v.UserID != userID {
		return ErrOTPCode
	}

	err = bcrypt.CompareHashAndPassword(v.VerificationCode, []byte(otp))
	if err != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(s.masterCode), []byte(otp)); err != nil {
			return ErrOTPCode
		}
	}

	return nil
}

// UpdateEmail updates user's email to provided new one in case of correct otp provided.
func (s *Service) UpdateEmail(ctx context.Context, userID uuid.UUID, email, otp string) error {
	if err := s.ValidateChangeEmailCode(ctx, userID, email, otp); err != nil {
		return fmt.Errorf("could not update email: %w", err)
	}

	if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
		ID:    userID,
		Email: email,
	}); err != nil {
		return fmt.Errorf("could not update email: %w", err)
	}

	if err := s.ur.UpdateUserVerifiedAt(ctx, repository.UpdateUserVerifiedAtParams{
		UserID:     userID,
		VerifiedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}); err != nil {
		return fmt.Errorf("could not verify email address: %w", err)
	}

	if err := s.ur.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.VerifyConfirmAccount,
		UserID:      userID,
	}); err != nil {
		// just log, not any error for user
		log.Printf("could not delete verification code for user with id=%s: %v", userID.String(), err)
	}

	return nil
}

// UpdateUsername ...
func (s *Service) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	if _, err := s.ur.GetUserByUsername(ctx, username); err == nil {
		return fmt.Errorf("user with username %s already exists", username)
	}

	if err := s.ur.UpdateUsername(ctx, repository.UpdateUsernameParams{ID: userID, Username: username}); err != nil {
		return fmt.Errorf("could not update username: %w", err)
	}

	return nil
}

// IsVerified returns if account is being verified.
func (s *Service) IsVerified(ctx context.Context, userID uuid.UUID) (bool, error) {
	u, err := s.ur.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user by provided id: %w", err)
	}

	return u.VerifiedAt.Valid, nil
}

// ResendOTP resends OTP to user by provided ID.
func (s *Service) ResendOTP(ctx context.Context, userID uuid.UUID) error {
	var otpHash []byte

	u, err := s.ur.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user by provided id: %w", err)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	if s.mail == nil {
		otpHash = []byte(s.masterCode)
	} else {
		otpHash, err = bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
		if err != nil {
			return fmt.Errorf("could not request resend otp: %w", err)
		}
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.VerifyConfirmAccount, // FIXME: it should be received from request
		UserID:           u.ID,
		Email:            u.Email,
		VerificationCode: otpHash,
	}); err != nil {
		return fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendVerificationCode(ctx, u.Email, otp); err != nil {
			return fmt.Errorf("could not send verification code: %w", err)
		}
	} else {
		// log data for debug mode
		log.Println("mail service is not set")
		log.Printf("[email verification] email: %s, otp: %s", u.Email, otp)
	}

	return nil
}

// RequestDestroyAccount requests password reset with email.
func (s *Service) RequestDestroyAccount(ctx context.Context, uid uuid.UUID) error {
	var otpHash []byte

	u, err := s.ur.GetUserByID(ctx, uid)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("user %w", ErrNotFound)
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	if s.mail == nil {
		otpHash = []byte(s.masterCode)
	} else {
		otpHash, err = bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
		if err != nil {
			return fmt.Errorf("could not request destroy account: %w", err)
		}
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.VerifyDestroyAccount,
		UserID:           u.ID,
		VerificationCode: otpHash,
	}); err != nil {
		return fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendDestroyAccountCode(ctx, u.Email, otp); err != nil {
			return fmt.Errorf("could not send reset password code: %w", err)
		}
	} else {
		// log data for debug mode
		log.Println("mail service is not set")
		log.Printf("[destroy account] email: %s, otp: %s", u.Email, otp)
	}

	return nil
}

// ValidateDestroyAccountCode validates destroy account code,
// it's needed to implement the destroy account flow on the client.
func (s *Service) ValidateDestroyAccountCode(ctx context.Context, uid uuid.UUID, otp string) error {
	v, err := s.ur.GetUserVerificationByUserID(ctx, repository.GetUserVerificationByUserIDParams{
		RequestType: repository.VerifyResetPassword,
		UserID:      uid,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("%w user with given email address", ErrNotFound)
		}
		return fmt.Errorf("could not get user with given email address: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(v.VerificationCode, []byte(otp))
	if err != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(s.masterCode), []byte(otp)); err != nil {
			return ErrOTPCode
		}
	}

	return nil
}

// DestroyAccount destroys account.
func (s *Service) DestroyAccount(ctx context.Context, uid uuid.UUID, otp string) error {
	if err := s.ValidateDestroyAccountCode(ctx, uid, otp); err != nil {
		return err
	}

	if err := s.ur.DestroyUser(ctx, uid); err != nil {
		return fmt.Errorf("could not destroy account: %w", err)
	}

	if err := s.ur.DeleteUserVerificationsByUserID(ctx, repository.DeleteUserVerificationsByUserIDParams{
		RequestType: repository.VerifyDestroyAccount,
		UserID:      uid,
	}); err != nil {
		// just log, not any error for user
		log.Printf("could not delete verification code for user with id=%s: %v", uid.String(), err)
	}

	return nil
}
