// Package service contains all the services/aggregates
package service

import (
	"cabtrips-data-api/repository"
)

type cache interface {
	Refresh() error
}

// Cache is BigCache based service
type Cache struct {
	repository.Cache
}

func NewCache(cacheProvided *repository.Cache) *Cache {
	return &Cache{
		Cache: *cacheProvided,
	}
}

func (c *Cache) Refresh() error {
	if err := c.Cache.Refresh(); err != nil {
		return err
	}
	return nil
}
