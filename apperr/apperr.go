package apperr

import "errors"

// NewBooking Errors
var (
	InvalidStatusErr        = errors.New("invalid status")
	InvalidDateErr          = errors.New("invalid date")
	InvalidUserIDErr        = errors.New("invalid user id")
	InvalidListingID        = errors.New("invalid listing id")
	BookingAlreadyExistsErr = errors.New("booking already exists")
)

// date error
var (
	InvalidPeriodErr = errors.New("invalid period")
	IntersectionErr  = errors.New("intersection date")
)

var (
	InvalidPriceErr = errors.New("invalid price")
)

// status Error
var (
	StatusCancel      = errors.New("status cancel")
	StatusConfirmed   = errors.New("status confirm")
	StatusExpire      = errors.New("status expire")
	UnsupportedStatus = errors.New("unsupported status")
	NotExpiredErr     = errors.New("not expired")
)

var (
	NotFoundErr = errors.New("not found")
)

type AppErr struct {
	FrontError   error
	ServiceError error
	Code         int
}

func NewAppErr(frontErr error, serviceErr error, code int) *AppErr {
	return &AppErr{
		FrontError:   frontErr,
		ServiceError: serviceErr,
		Code:         code,
	}
}

func (err *AppErr) Error() string {
	return err.ServiceError.Error()
}

func (err *AppErr) GetMapError() map[string]any {
	return map[string]any{
		"Message": err.FrontError.Error(),
		"Service": err.ServiceError.Error(),
		"Code":    err.Code,
	}
}
