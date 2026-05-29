package models

import "github.com/mupt-ai/dari-coffee-cli/internal/openapi"

type AddressCheckStatus = openapi.AddressCheckStatus

const (
	AddressCheckStatusEligible            = openapi.AddressCheckStatusEligible
	AddressCheckStatusAddressUnresolved   = openapi.AddressCheckStatusAddressUnresolved
	AddressCheckStatusOutsideDeliveryZone = openapi.AddressCheckStatusOutsideDeliveryZone
	AddressCheckStatusGeocoderUnavailable = openapi.AddressCheckStatusGeocoderUnavailable
)

type AddressCheckRequest = openapi.AddressCheckRequest
type AddressCheck = openapi.AddressCheck
