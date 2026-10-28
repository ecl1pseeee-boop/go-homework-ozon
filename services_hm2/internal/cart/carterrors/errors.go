package carterrors

import "errors"

var (
	ErrInvalidData        = errors.New("invalid data")
	ErrCartNotFound       = errors.New("cart not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductUnavailable = errors.New("product unavailable")
	ErrCartEmpty          = errors.New("cart empty")
)
