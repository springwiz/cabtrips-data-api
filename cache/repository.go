package cache

import "cabtrips-data-api/model"

// Repository defines the repository specification
type Repository interface {
	GetCabtripByMedallion(medallion string) ([]model.Cabtrip, error)

	GetCabtripByMedallionAndPickupdate(medallion string, pickupDate string) ([]model.Cabtrip, error)

	Refresh() error
}
