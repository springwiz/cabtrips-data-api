package model

import (
	"database/sql"
	"time"
)

// The Cabtrip entity
type Cabtrip struct {
	Medallion string `json:"medallion"`

	HackLicense string `json:"hackLicense"`

	VendorID string `json:"vendorID"`

	RateCode uint32 `json:"rateCode"`

	StoreAndFwdFlag string `json:"storeAndFwdFlag"`

	PickupDatetime time.Time `json:"pickupDatetime"`

	DropoffDatetime time.Time `json:"dropoffDatetime"`

	PassengerCount uint32 `json:"passengerCount"`

	TripTimeInSecs uint32 `json:"tripTimeInSecs"`

	TripDistance float32 `json:"tripDistance"`

	PickupLongitude sql.NullFloat64 `json:"pickupLongitude"`

	PickupLatitude sql.NullFloat64 `json:"pickupLatitude"`

	DropoffLongitude sql.NullFloat64 `json:"dropoffLongitude"`

	DropoffLatitude sql.NullFloat64 `json:"dropoffLatitude"`
}
