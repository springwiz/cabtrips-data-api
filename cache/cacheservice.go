package cache

type CacheService interface {
	Refresh() error
}

// Cache is BigCache based service
type CacheServiceImpl struct {
	Cache Repository
}

func NewCacheService(cacheProvided Repository) *CacheServiceImpl {
	return &CacheServiceImpl{
		Cache: cacheProvided,
	}
}

func (c *CacheServiceImpl) Refresh() error {
	if err := c.Cache.Refresh(); err != nil {
		return err
	}
	return nil
}
