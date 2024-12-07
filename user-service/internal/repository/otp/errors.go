package otp

import "errors"

var (
	ErrFailedToCreateSession = errors.New("failed to create session")
	ErrFailedToRetrieveOTP   = errors.New("failed to retrieve OTP")
	ErrOTPNotFound           = errors.New("OTP not found or expired")
	ErrFailedToStoreOTP      = errors.New("failed to store OTP")
)
