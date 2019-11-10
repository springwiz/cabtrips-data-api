// Package handler provides with handler functions for handling the various HTTP Requests
package handler

import (
	"cabtrips-data-api/model"
)

// Repository defines the repository specification
type Repository interface {
	GetCabtripByMedallion(medallion string) (*model.Cabtrip, error)

	GetCabtripByMedallionAndPickupdate(medallion string, pickupDate string) (*model.Cabtrip, error)

	Refresh()
}
