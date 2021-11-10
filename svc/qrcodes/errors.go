package qrcodes

import "errors"

var (
	// ErrInvalidParameter indicates that passed invalid parameter.
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrQRCodeExpired indicates that QR code is expired.
	ErrQRCodeExpired = errors.New("QR code is expired")
	ErrQRCodeInvalid = errors.New("QR code is not valid yet")
	ErrQRCodeScanned = errors.New("QR code can only be scanned once")
)
