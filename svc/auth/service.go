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
	"github.com/SatorNetwork/sator-api/internal/rbac"
	"github.com/SatorNetwork/sator-api/internal/sumsub"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/SatorNetwork/sator-api/svc/auth/repository"

	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Service struct
	Service struct {
		ur                    userRepository
		ws                    walletService
		jwt                   jwtInteractor
		mail                  mailer
		ic                    invitationsClient
		kyc                   kycClient
		otpLen                int
		masterCode            string
		blacklistEmailDomains []string
	}

	Whitelist struct {
		AllowedType  string `json:"allowed_type"`
		AllowedValue string `json:"allowed_value"`
	}

	Blacklist struct {
		RestrictedType  string `json:"restricted_type"`
		RestrictedValue string `json:"restricted_value"`
	}

	UserStatus struct {
		Email       string `json:"email"`
		Username    string `json:"username"`
		IsDisabled  bool   `json:"is_disabled"`
		BlockReason string `json:"block_reason,omitempty"`
		IsFinal     bool   `json:"is_final"`
		KYCStatus   string `json:"kyc_status,omitempty"`
	}

	// ServiceOption function
	// interface to extend service via options
	ServiceOption func(*Service)

	jwtInteractor interface {
		NewWithRefreshToken(userID uuid.UUID, username, role string) (access, refresh string, err error)
	}

	userRepository interface {
		// user
		CreateUser(ctx context.Context, arg repository.CreateUserParams) (repository.User, error)
		GetUserByEmail(ctx context.Context, email string) (repository.User, error)
		GetUserBySanitizedEmail(ctx context.Context, email string) (repository.User, error)
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

		// Blacklist
		IsEmailBlacklisted(ctx context.Context, email string) (bool, error)
		AddToBlacklist(ctx context.Context, arg repository.AddToBlacklistParams) (repository.Blacklist, error)
		DeleteFromBlacklist(ctx context.Context, arg repository.DeleteFromBlacklistParams) error
		GetBlacklist(ctx context.Context, arg repository.GetBlacklistParams) ([]repository.Blacklist, error)
		GetBlacklistByRestrictedValue(ctx context.Context, arg repository.GetBlacklistByRestrictedValueParams) ([]repository.Blacklist, error)

		// Whitelist
		IsEmailWhitelisted(ctx context.Context, email string) (bool, error)
		AddToWhitelist(ctx context.Context, arg repository.AddToWhitelistParams) (repository.Whitelist, error)
		DeleteFromWhitelist(ctx context.Context, arg repository.DeleteFromWhitelistParams) error
		GetWhitelist(ctx context.Context, arg repository.GetWhitelistParams) ([]repository.Whitelist, error)
		GetWhitelistByAllowedValue(ctx context.Context, arg repository.GetWhitelistByAllowedValueParams) ([]repository.Whitelist, error)

		LinkDeviceToUser(ctx context.Context, arg repository.LinkDeviceToUserParams) error
		DoesUserHaveMoreThanOneAccount(ctx context.Context, userID uuid.UUID) (bool, error)

		// KYC
		UpdateKYCStatus(ctx context.Context, arg repository.UpdateKYCStatusParams) error
		UpdateUserStatus(ctx context.Context, arg repository.UpdateUserStatusParams) error
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

	kycClient interface {
		GetSDKAccessTokenByApplicantID(ctx context.Context, applicantID string) (string, error)
		GetSDKAccessTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error)
		GetByExternalUserID(ctx context.Context, userID uuid.UUID) (*sumsub.Response, error)
	}

	// JWTs
	Token struct {
		AccessToken  string
		RefreshToken string
	}

	User struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation.
func NewService(ji jwtInteractor, ur userRepository, ws walletService, ic invitationsClient, kyc kycClient, opt ...ServiceOption) *Service {
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
	if kyc == nil {
		log.Fatalln("kyc client is not set")
	}

	s := &Service{jwt: ji, ur: ur, ic: ic, kyc: kyc, ws: ws, otpLen: 5}

	// Set up options.
	for _, o := range opt {
		o(s)
	}

	return s
}

// Login by email and password, returns token.
func (s *Service) Login(ctx context.Context, email, password, deviceID string) (Token, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return Token{}, ErrInvalidCredentials
		}
		return Token{}, fmt.Errorf("could not log in: %w", err)
	}

	if deviceID == "" {
		return Token{}, ErrEmptyDeviceID
	}

	if err := s.ur.LinkDeviceToUser(ctx, repository.LinkDeviceToUserParams{
		UserID:   user.ID,
		DeviceID: deviceID,
	}); err != nil {
		log.Printf("could not link device to user: %v", err)
	}

	if !user.SanitizedEmail.Valid || len(user.SanitizedEmail.String) < 5 {
		// Sanitize email address
		sanitizedEmail, err := utils.SanitizeEmail(user.Email)
		if err != nil {
			return Token{}, validator.NewValidationError(url.Values{
				"email": []string{ErrInvalidEmailFormat.Error()},
			})
		}

		if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
			ID:             user.ID,
			Email:          user.Email,
			SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
		}); err != nil {
			log.Printf("could not add sanitizesd email for user with id=%s and email=%s: %v", user.ID, user.Email, err)
		}
	}

	if user.Disabled {
		return Token{}, ErrUserIsDisabled
	}

	if !strings.Contains(email, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, user.Email); yes {
			return Token{}, ErrUserIsDisabled
		}
		if user.SanitizedEmail.Valid {
			if yes, _ := s.ur.IsEmailBlacklisted(ctx, user.SanitizedEmail.String); yes {
				return Token{}, ErrUserIsDisabled
			}
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, email); !yes {
			return Token{}, ErrUserIsDisabled
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return Token{}, ErrInvalidCredentials
	}

	token, refreshToken, err := s.jwt.NewWithRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return Token{}, fmt.Errorf("could not generate new access token: %w", err)
	}

	return Token{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

// Logout revokes JWT token.
func (s *Service) Logout(ctx context.Context, tid string) error {
	// TODO: add JWT id into the revoked tokens list
	return nil
}

// RefreshToken returns new jwt string.
func (s *Service) RefreshToken(ctx context.Context, uid uuid.UUID, username, role, deviceID string) (Token, error) {
	if deviceID == "" {
		return Token{}, ErrEmptyDeviceID
	}

	if err := s.ur.LinkDeviceToUser(ctx, repository.LinkDeviceToUserParams{
		UserID:   uid,
		DeviceID: deviceID,
	}); err != nil {
		log.Printf("could not link device to user: %v", err)
	}

	u, err := s.ur.GetUserByID(ctx, uid)
	if err != nil {
		return Token{}, fmt.Errorf("could not refresh access token: %w", err)
	}

	if !u.SanitizedEmail.Valid || len(u.SanitizedEmail.String) < 5 {
		// Sanitize email address
		sanitizedEmail, err := utils.SanitizeEmail(u.Email)
		if err != nil {
			return Token{}, ErrInvalidEmailFormat
		}

		if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
			ID:             u.ID,
			Email:          u.Email,
			SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
		}); err != nil {
			log.Printf("could not add sanitizesd email for user with id=%s and email=%s: %v", u.ID, u.Email, err)
		}
	}

	if u.Disabled {
		return Token{}, ErrUserIsDisabled
	}

	if !strings.HasSuffix(u.Email, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.Email); yes {
			return Token{}, ErrUserIsDisabled
		}
		if u.SanitizedEmail.Valid {
			if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.SanitizedEmail.String); yes {
				return Token{}, ErrUserIsDisabled
			}
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, u.Email); !yes {
			return Token{}, ErrUserIsDisabled
		}
	}

	// TODO: add JWT id into the revoked tokens list
	token, refreshToken, err := s.jwt.NewWithRefreshToken(u.ID, u.Username, u.Role)
	if err != nil {
		return Token{}, fmt.Errorf("could not generate new access token: %w", err)
	}

	return Token{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

// SignUp registers account with email, password and username.
func (s *Service) SignUp(ctx context.Context, email, password, username, deviceID string) (Token, error) {
	var otpHash []byte

	email = strings.ToLower(strings.TrimSpace(email))

	// Sanitize email address
	sanitizedEmail, err := utils.SanitizeEmail(email)
	if err != nil {
		return Token{}, validator.NewValidationError(url.Values{
			"email": []string{ErrInvalidEmailFormat.Error()},
		})
	}

	if !strings.Contains(email, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, email); yes {
			return Token{}, validator.NewValidationError(url.Values{
				"email": []string{ErrRestrictedEmailDomain.Error()},
			})
		}
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, sanitizedEmail); yes {
			return Token{}, validator.NewValidationError(url.Values{
				"email": []string{ErrRestrictedEmailDomain.Error()},
			})
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, email); !yes {
			return Token{}, validator.NewValidationError(url.Values{
				"email": []string{ErrRestrictedEmailDomain.Error()},
			})
		}
	}

	// Check if the passed email address is not taken yet
	if _, err := s.ur.GetUserByEmail(ctx, email); err == nil {
		return Token{}, validator.NewValidationError(url.Values{
			"email": []string{"email is already taken"},
		})
	} else if !db.IsNotFoundError(err) {
		return Token{}, fmt.Errorf("could not create a new account: %w", err)
	}

	if _, err := s.ur.GetUserBySanitizedEmail(ctx, sanitizedEmail); err == nil {
		return Token{}, validator.NewValidationError(url.Values{
			"email": []string{"email is already taken"},
		})
	} else if !db.IsNotFoundError(err) {
		return Token{}, fmt.Errorf("could not create a new account: %w", err)
	}

	// Check if the passed username is not taken yet
	if _, err := s.ur.GetUserByUsername(ctx, username); err == nil {
		return Token{}, validator.NewValidationError(url.Values{
			"username": []string{"username is already taken"},
		})
	} else if !db.IsNotFoundError(err) {
		return Token{}, fmt.Errorf("could not create a new account: %w", err)
	}

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Token{}, fmt.Errorf("could not create a new account: %w", err)
	}

	u, err := s.ur.CreateUser(ctx, repository.CreateUserParams{
		Email:          email,
		SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
		Password:       passwdHash,
		Username:       username,
		Role:           rbac.RoleUser.String(),
	})
	if err != nil {
		return Token{}, fmt.Errorf("could not create a new account: %w", err)
	}

	if deviceID == "" {
		return Token{}, ErrEmptyDeviceID
	}

	if err := s.ur.LinkDeviceToUser(ctx, repository.LinkDeviceToUserParams{
		UserID:   u.ID,
		DeviceID: deviceID,
	}); err != nil {
		log.Printf("could not link device to user: %v", err)
	}

	otp := random.String(uint8(s.otpLen), random.Numeric)
	if s.mail == nil {
		otpHash = []byte(s.masterCode)
	} else {
		otpHash, err = bcrypt.GenerateFromPassword([]byte(otp), bcrypt.MinCost)
		if err != nil {
			return Token{}, fmt.Errorf("could not create a new account: %w", err)
		}
	}

	if err := s.ur.CreateUserVerification(ctx, repository.CreateUserVerificationParams{
		RequestType:      repository.VerifyConfirmAccount,
		UserID:           u.ID,
		Email:            email,
		VerificationCode: otpHash,
	}); err != nil {
		return Token{}, fmt.Errorf("could not generate verification code: %w", err)
	}

	if s.mail != nil {
		if err := s.mail.SendVerificationCode(ctx, email, otp); err != nil {
			return Token{}, fmt.Errorf("could not send verification code: %w", err)
		}
	} else {
		// log data for debug mode
		log.Println("mail service is not set")
		log.Printf("[email verification] email: %s, otp: %s", email, otp)
	}

	token, refreshToken, err := s.jwt.NewWithRefreshToken(u.ID, u.Username, u.Role)
	if err != nil {
		return Token{}, fmt.Errorf("could not generate new access token: %w", err)
	}

	if isInvited, _ := s.ic.IsEmailInvited(ctx, email); isInvited {
		if err := s.ic.AcceptInvitation(ctx, u.ID, email); err != nil {
			log.Printf("could not accept invitation for user id = %s: %v", u.ID, err)
		}
	}

	return Token{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

// ForgotPassword requests password reset with email.
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	var otpHash []byte
	email = strings.ToLower(strings.TrimSpace(email))

	u, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return errors.New("please check if you set correct email")
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	if !u.SanitizedEmail.Valid || len(u.SanitizedEmail.String) < 5 {
		// Sanitize email address
		sanitizedEmail, err := utils.SanitizeEmail(u.Email)
		if err != nil {
			return ErrInvalidEmailFormat
		}

		if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
			ID:             u.ID,
			Email:          u.Email,
			SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
		}); err != nil {
			log.Printf("could not add sanitizesd email for user with id=%s and email=%s: %v", u.ID, u.Email, err)
		}
	}

	if u.Disabled {
		return ErrUserIsDisabled
	}

	if !strings.Contains(u.Email, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.Email); yes {
			return ErrUserIsDisabled
		}
		if u.SanitizedEmail.Valid {
			if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.SanitizedEmail.String); yes {
				return ErrUserIsDisabled
			}
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, u.Email); !yes {
			return ErrUserIsDisabled
		}
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
	email = strings.ToLower(strings.TrimSpace(email))

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
	email = strings.ToLower(strings.TrimSpace(email))

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
		return fmt.Errorf("could not get user bu id: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(oldPassword)); err != nil {
		return validator.NewValidationError(url.Values{
			"old_password": []string{"invalid current password"},
		})
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

	if !u.SanitizedEmail.Valid || len(u.SanitizedEmail.String) < 5 {
		// Sanitize email address
		sanitizedEmail, err := utils.SanitizeEmail(u.Email)
		if err != nil {
			return ErrInvalidEmailFormat
		}

		if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
			ID:             u.ID,
			Email:          u.Email,
			SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
		}); err != nil {
			log.Printf("could not add sanitizesd email for user with id=%s and email=%s: %v", u.ID, u.Email, err)
		}
	}

	if u.VerifiedAt.Valid {
		return ErrEmailAlreadyVerified
	}

	if u.Disabled {
		return ErrUserIsDisabled
	}

	if !strings.Contains(u.Email, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.Email); yes {
			return ErrUserIsDisabled
		}
		if u.SanitizedEmail.Valid {
			if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.SanitizedEmail.String); yes {
				return ErrUserIsDisabled
			}
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, u.Email); !yes {
			return ErrUserIsDisabled
		}
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
	email = strings.ToLower(strings.TrimSpace(email))

	// Sanitize email address
	sanitizedEmail, err := utils.SanitizeEmail(email)
	if err != nil {
		return validator.NewValidationError(url.Values{
			"email": []string{ErrInvalidEmailFormat.Error()},
		})
	}

	u, err := s.ur.GetUserByID(ctx, userID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return fmt.Errorf("user %w", ErrNotFound)
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	if !u.SanitizedEmail.Valid || len(u.SanitizedEmail.String) < 5 {
		// Sanitize email address
		sanitizedEmail, err := utils.SanitizeEmail(u.Email)
		if err != nil {
			return ErrInvalidEmailFormat
		}

		if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
			ID:             u.ID,
			Email:          u.Email,
			SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
		}); err != nil {
			log.Printf("could not add sanitizesd email for user with id=%s and email=%s: %v", u.ID, u.Email, err)
		}
	}

	if u.Disabled {
		return ErrUserIsDisabled
	}

	if !strings.Contains(sanitizedEmail, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.Email); yes {
			return ErrUserIsDisabled
		}
		if u.SanitizedEmail.Valid {
			if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.SanitizedEmail.String); yes {
				return ErrUserIsDisabled
			}
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, sanitizedEmail); !yes {
			return ErrUserIsDisabled
		}
	}

	if _, err := s.ur.GetUserBySanitizedEmail(ctx, sanitizedEmail); err == nil {
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
	email = strings.ToLower(strings.TrimSpace(email))

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
	email = strings.ToLower(strings.TrimSpace(email))

	// Sanitize email address
	sanitizedEmail, err := utils.SanitizeEmail(email)
	if err != nil {
		return validator.NewValidationError(url.Values{
			"email": []string{ErrInvalidEmailFormat.Error()},
		})
	}

	if err := s.ValidateChangeEmailCode(ctx, userID, email, otp); err != nil {
		return fmt.Errorf("could not update email: %w", err)
	}

	if err := s.ur.UpdateUserEmail(ctx, repository.UpdateUserEmailParams{
		ID:             userID,
		Email:          email,
		SanitizedEmail: sql.NullString{String: sanitizedEmail, Valid: true},
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

	if u.Disabled {
		return false, ErrUserIsDisabled
	}

	if !strings.Contains(u.Email, "@sator.io") {
		if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.Email); yes {
			return false, ErrUserIsDisabled
		}
		if u.SanitizedEmail.Valid {
			if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.SanitizedEmail.String); yes {
				return false, ErrUserIsDisabled
			}
		}
		if yes, _ := s.ur.IsEmailWhitelisted(ctx, u.Email); !yes {
			return false, ErrUserIsDisabled
		}
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

// AddToWhitelist used for add allowed type and value to whitelist.
func (s *Service) AddToWhitelist(ctx context.Context, allowedType, allowedValue string) error {
	allowedValue = strings.ToLower(strings.TrimSpace(allowedValue))
	if _, err := s.ur.AddToWhitelist(ctx, repository.AddToWhitelistParams{
		AllowedType:  allowedType,
		AllowedValue: allowedValue,
	}); err != nil {
		return fmt.Errorf("could not add llowed type and value to whitelist: %w", err)
	}

	return nil
}

// GetWhitelist returns whitelist with pagination.
func (s *Service) GetWhitelist(ctx context.Context, limit, offset int32) ([]Whitelist, error) {
	whitelist, err := s.ur.GetWhitelist(ctx, repository.GetWhitelistParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get whitelist: %w", err)
	}

	result := make([]Whitelist, 0, len(whitelist))
	for _, w := range whitelist {
		result = append(result, Whitelist{
			AllowedType:  w.AllowedType,
			AllowedValue: w.AllowedValue,
		})
	}

	return result, nil
}

// SearchInWhitelist returns whitelist with pagination.
func (s *Service) SearchInWhitelist(ctx context.Context, limit, offset int32, query string) ([]Whitelist, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	whitelist, err := s.ur.GetWhitelistByAllowedValue(ctx, repository.GetWhitelistByAllowedValueParams{
		Query:     query,
		OffsetVal: offset,
		LimitVal:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get whitelist by query %v: %w", query, err)
	}

	result := make([]Whitelist, 0, len(whitelist))
	for _, w := range whitelist {
		result = append(result, Whitelist{
			AllowedType:  w.AllowedType,
			AllowedValue: w.AllowedValue,
		})
	}

	return result, nil
}

// DeleteFromWhitelist used for delete allowed type and value from whitelist.
func (s *Service) DeleteFromWhitelist(ctx context.Context, allowedType, allowedValue string) error {
	allowedValue = strings.ToLower(strings.TrimSpace(allowedValue))
	if err := s.ur.DeleteFromWhitelist(ctx, repository.DeleteFromWhitelistParams{
		AllowedType:  allowedType,
		AllowedValue: allowedValue,
	}); err != nil {
		return fmt.Errorf("could not delete record from whitelist: %w", err)
	}

	return nil
}

// AddToBlacklist used for add restricted type and value to blacklist.
func (s *Service) AddToBlacklist(ctx context.Context, restrictedType, restrictedValue string) error {
	restrictedValue = strings.ToLower(strings.TrimSpace(restrictedValue))
	if _, err := s.ur.AddToBlacklist(ctx, repository.AddToBlacklistParams{
		RestrictedType:  restrictedType,
		RestrictedValue: restrictedValue,
	}); err != nil {
		return fmt.Errorf("could not add llowed type and value to blacklist: %w", err)
	}

	return nil
}

// GetBlacklist returns blacklist with pagination.
func (s *Service) GetBlacklist(ctx context.Context, limit, offset int32) ([]Blacklist, error) {
	blacklist, err := s.ur.GetBlacklist(ctx, repository.GetBlacklistParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get blacklist: %w", err)
	}

	result := make([]Blacklist, 0, len(blacklist))
	for _, b := range blacklist {
		result = append(result, Blacklist{
			RestrictedType:  b.RestrictedType,
			RestrictedValue: b.RestrictedValue,
		})
	}

	return result, nil
}

// SearchInBlacklist returns blacklist with pagination.
func (s *Service) SearchInBlacklist(ctx context.Context, limit, offset int32, query string) ([]Blacklist, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	blacklist, err := s.ur.GetBlacklistByRestrictedValue(ctx, repository.GetBlacklistByRestrictedValueParams{
		Query:     query,
		OffsetVal: offset,
		LimitVal:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get blacklist by query %v: %w", query, err)
	}

	result := make([]Blacklist, 0, len(blacklist))
	for _, w := range blacklist {
		result = append(result, Blacklist{
			RestrictedType:  w.RestrictedType,
			RestrictedValue: w.RestrictedValue,
		})
	}

	return result, nil
}

// DeleteFromBlacklist used for delete restricted type and value from blacklist.
func (s *Service) DeleteFromBlacklist(ctx context.Context, restrictedType, restrictedValue string) error {
	restrictedValue = strings.ToLower(strings.TrimSpace(restrictedValue))
	if err := s.ur.DeleteFromBlacklist(ctx, repository.DeleteFromBlacklistParams{
		RestrictedType:  restrictedType,
		RestrictedValue: restrictedValue,
	}); err != nil {
		return fmt.Errorf("could not delete record from blacklist: %w", err)
	}

	return nil
}

// GetAccessTokenByUserID returns access token for web or mobile SDKs by user id.
func (s *Service) GetAccessTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := s.kyc.GetSDKAccessTokenByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("could not get access token: %w", err)
	}

	return token, nil
}

// GetUserStatus returns user status by user email.
func (s *Service) GetUserStatus(ctx context.Context, email string) (UserStatus, error) {
	u, err := s.ur.GetUserByEmail(ctx, email)
	if err != nil {
		if db.IsNotFoundError(err) {
			return UserStatus{}, fmt.Errorf("could not found user with email %s", email)
		}
		return UserStatus{}, err
	}

	var (
		kycStatus string
		reason    string
		isFinal   bool
	)

	if u.Disabled {
		if strings.Contains(u.BlockReason.String, "multiple accounts") {
			reason = "User has multiple accounts"
			isFinal = true
		} else if strings.Contains(u.BlockReason.String, "invalid email") {
			reason = "Invalid email address"
			isFinal = false
		} else if strings.Contains(u.BlockReason.String, "frequent rewards") {
			reason = "Suspicion of fraud"
			isFinal = false
		}

		if !isFinal {
			if yes, _ := s.ur.DoesUserHaveMoreThanOneAccount(ctx, u.ID); yes {
				reason = fmt.Sprintf("%s. Also, the user has multiple accounts", reason)
				isFinal = true
			}
		}
	} else if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.Email); yes {
		reason = "The email address has been found on the blacklist"
		isFinal = true
	} else if yes, _ := s.ur.IsEmailBlacklisted(ctx, u.SanitizedEmail.String); u.SanitizedEmail.Valid && yes {
		reason = "The email address has been found on the blacklist"
		isFinal = true
	} else {
		if u.KycStatus.Valid {
			switch u.KycStatus.String {
			case sumsub.KYCStatusNotVerified:
				kycStatus = "Not verified"
			case sumsub.KYCStatusRetry:
				kycStatus = "Invalid documents or bad quality of selfie/docs. The user should upload valid documents or/and retake a selfie."
			case sumsub.KYCStatusApproved:
				kycStatus = "Verified"
			case sumsub.KYCStatusRejected:
				kycStatus = "The user was rejected. It's the final decision and cannot be changed."
			case sumsub.KYCStatusInProgress:
				kycStatus = "Verification has not been completed yet."
			case sumsub.KYCStatusInit:
				kycStatus = "The user has not uploaded all documents yet."
			default:
				kycStatus = "N/A"
			}
		}
	}

	return UserStatus{
		Email:       u.Email,
		Username:    u.Username,
		IsDisabled:  u.Disabled,
		BlockReason: reason,
		IsFinal:     isFinal,
		KYCStatus:   kycStatus,
	}, nil
}

// VerificationCallback endpoint for kyc service webhook. And used for store user status.
func (s *Service) VerificationCallback(ctx context.Context, userID uuid.UUID) error {
	resp, err := s.kyc.GetByExternalUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not get external user by id: %w", err)
	}

	if resp.Review.ReviewStatus == sumsub.KYCProviderStatusInit {
		err = s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
			KycStatus: sumsub.KYCStatusInit,
			ID:        userID,
		})
		if err != nil {
			return fmt.Errorf("could not update kyc status for user: %v: %w", userID, err)
		}

		return nil
	}

	if resp.Review.ReviewStatus == sumsub.KYCProviderStatusPending {
		err = s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
			KycStatus: sumsub.KYCStatusInProgress,
			ID:        userID,
		})
		if err != nil {
			return fmt.Errorf("could not update kyc status for user: %v: %w", userID, err)
		}
	}

	if resp.Review.ReviewResult.ReviewAnswer == sumsub.KYCProviderStatusGreen {
		if resp.Review.ReviewResult.ReviewRejectType == sumsub.KYCProviderStatusRetry {
			err = s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
				KycStatus: sumsub.KYCStatusRetry,
				ID:        userID,
			})
			if err != nil {
				return fmt.Errorf("could not update kyc status for user: %v: %w", userID, err)
			}
		}

		err = s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
			KycStatus: sumsub.KYCStatusApproved,
			ID:        userID,
		})
		if err != nil {
			return fmt.Errorf("could not update kyc status for user: %v: %w", userID, err)
		}
	}

	if resp.Review.ReviewResult.ReviewAnswer == sumsub.KYCProviderStatusRed && resp.Review.ReviewResult.ReviewRejectType == sumsub.KYCProviderStatusFinal {
		err := s.ur.UpdateUserStatus(ctx, repository.UpdateUserStatusParams{
			ID:       userID,
			Disabled: true,
			BlockReason: sql.NullString{
				String: sumsub.KYCRejectLabelsMap()(resp.Review.ReviewResult.RejectLabels),
				Valid:  true,
			},
		})
		if err != nil {
			return fmt.Errorf("could not block user: %v: %w", userID, err)
		}

		err = s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
			KycStatus: sumsub.KYCStatusRejected,
			ID:        userID,
		})
		if err != nil {
			return fmt.Errorf("could not update kyc status for user: %v: %w", userID, err)
		}
	}

	if resp.Review.ReviewResult.ReviewAnswer == sumsub.KYCProviderStatusRed && resp.Review.ReviewResult.ReviewRejectType == sumsub.KYCProviderStatusRetry {
		err = s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
			KycStatus: sumsub.KYCStatusRetry,
			ID:        userID,
		})
		if err != nil {
			return fmt.Errorf("could not update kyc status for user: %v: %w", userID, err)
		}
	}

	return nil
}

func (s *Service) UpdateKYCStatus(ctx context.Context, uid uuid.UUID) error {
	if err := s.ur.UpdateKYCStatus(ctx, repository.UpdateKYCStatusParams{
		KycStatus: sumsub.KYCStatusApproved,
		ID:        uid,
	}); err != nil {
		return fmt.Errorf("could not update kyc status for user: %v: %w", uid, err)
	}

	return nil
}

func (s *Service) GetUsernameByID(ctx context.Context, uid uuid.UUID) (string, error) {
	user, err := s.ur.GetUserByID(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("could not get user by id: %v: %w", uid, err)
	}

	return user.Username, nil
}
