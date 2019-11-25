package cabs

import (
	"cabtrips-data-api/cache"
	"cabtrips-data-api/model"
)

type CabtripService interface {
	GetCabtripByMedallion(medallion, cacheKey string) ([]model.Cabtrip, error)

	GetCabtripByMedallionAndPickupdate(medallion, pickupDate, cacheKey string) ([]model.Cabtrip, error)
}

type CabtripServiceImpl struct {
	Mysql Repository
	Cache cache.Repository
}

func newCabtripService(cacheProvided cache.Repository, mysqlProvided Repository) *CabtripServiceImpl {
	return &CabtripServiceImpl{
		Mysql: mysqlProvided,
		Cache: cacheProvided,
	}
}

func (c *CabtripServiceImpl) GetCabtripByMedallion(medallion, cacheKey string) ([]model.Cabtrip, error) {
	if len(cacheKey) > 0 {
		return c.Cache.GetCabtripByMedallion(medallion)
	}
	return c.Mysql.GetCabtripByMedallion(medallion)
}

func (c *CabtripServiceImpl) GetCabtripByMedallionAndPickupdate(
	medallion, pickupDate, cacheKey string,
) ([]model.Cabtrip, error) {
	if len(cacheKey) > 0 {
		return c.Cache.GetCabtripByMedallionAndPickupdate(medallion, pickupDate)
	}
	return c.Mysql.GetCabtripByMedallionAndPickupdate(medallion, pickupDate)
}
