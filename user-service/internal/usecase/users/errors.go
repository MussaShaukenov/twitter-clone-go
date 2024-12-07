package users

import "errors"

var (
	ErrFailedToRetrieveOTP = errors.New("failed to retrieve OTP")
	ErrInvalidOTP          = errors.New("invalid OTP")
	ErrAgeRestrict         = errors.New("You must be 14 years or older to use this service")
)
