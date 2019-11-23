// Package service contains all the services/aggregates
package service

import (
	"cabtrips-data-api/model"
	"cabtrips-data-api/repository"
)

type Cabtrip interface {
	GetCabtripByMedallion(medallion, cacheKey string) ([]model.Cabtrip, error)

	GetCabtripByMedallionAndPickupdate(medallion, pickupDate, cacheKey string) ([]model.Cabtrip, error)
}

type CabtripImpl struct {
	repository.Mysql
	repository.Cache
}

func NewCabtrip(cacheProvided *repository.Cache, mysqlProvided *repository.Mysql) *CabtripImpl {
	return &CabtripImpl{
		Mysql: *mysqlProvided,
		Cache: *cacheProvided,
	}
}

func (c *CabtripImpl) GetCabtripByMedallion(medallion, cacheKey string) ([]model.Cabtrip, error) {
	if len(cacheKey) > 0 {
		return c.Cache.GetCabtripByMedallion(medallion)
	}
	return c.Mysql.GetCabtripByMedallion(medallion)
}

func (c *CabtripImpl) GetCabtripByMedallionAndPickupdate(
	medallion, pickupDate, cacheKey string,
) ([]model.Cabtrip, error) {
	if len(cacheKey) > 0 {
		return c.Cache.GetCabtripByMedallionAndPickupdate(medallion, pickupDate)
	}
	return c.Mysql.GetCabtripByMedallionAndPickupdate(medallion, pickupDate)
}
