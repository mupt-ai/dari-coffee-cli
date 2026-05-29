package models

import "github.com/mupt-ai/dari-coffee-cli/internal/openapi"

type ErrorCode = openapi.ErrorCode

const (
	ErrorCodeInvalidOrder        = openapi.ErrorCodeInvalidOrder
	ErrorCodeMenuItemUnavailable = openapi.ErrorCodeMenuItemUnavailable
	ErrorCodeAddressUnresolved   = openapi.ErrorCodeAddressUnresolved
	ErrorCodeGeocoderUnavailable = openapi.ErrorCodeGeocoderUnavailable
	ErrorCodeCheckoutUnavailable = openapi.ErrorCodeCheckoutUnavailable
	ErrorCodePaymentUnavailable  = openapi.ErrorCodePaymentUnavailable
	ErrorCodeServiceOff          = openapi.ErrorCodeServiceOff
	ErrorCodeOutsideHours        = openapi.ErrorCodeOutsideHours
)

type ErrorResponse = openapi.ErrorResponse
