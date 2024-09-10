package cache

import "time"

type CacheOption func(*Cache)

func SetLifetime(interval time.Duration) CacheOption {
	return func(c *Cache) {
		c.lifetime = interval
	}
}
