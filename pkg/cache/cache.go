package cache

import (
	"sync"
	"time"
)

type cleaner struct {
	id      uint64
	isReady <-chan time.Time
}

type Cache struct {
	lifetime   time.Duration
	storage    map[uint64]*objectCache
	rWMutex    sync.RWMutex
	deathQueue chan cleaner
}

func New(opts ...CacheOption) *Cache {
	c := &Cache{
		defaultLifetime,
		make(map[uint64]*objectCache),
		sync.RWMutex{},
		make(chan cleaner),
	}

	for _, opt := range opts {
		opt(c)
	}

	go c.startCleanupDaemon()

	return c
}

func (c *Cache) startCleanupDaemon() {
	for {
		select {
		case toClean := <-c.deathQueue:
			<-toClean.isReady
			c.clean(toClean.id)
		}
	}
}

func (c *Cache) Set(user *UserReadModel) {
	// Calculate expiration time right after method call, without waiting for RW mutex
	expireAt := time.Now().Add(c.lifetime)

	c.rWMutex.Lock()
	c.storage[user.id] = &objectCache{user, expireAt}
	c.rWMutex.Unlock()

	c.deathQueue <- cleaner{user.id, time.After(c.lifetime)}
}

func (c *Cache) Get(id uint64) *UserReadModel {
	c.rWMutex.RLock()
	defer c.rWMutex.RUnlock()

	userCache, exists := c.storage[id]

	if !exists {
		return nil
	}

	if time.Now().After(userCache.expireAt) {
		return nil
	}

	return userCache.object
}

func (c *Cache) clean(id uint64) {
	c.rWMutex.Lock()
	defer c.rWMutex.Unlock()

	object, exists := c.storage[id]

	if !exists {
		return
	}

	if time.Now().After(object.expireAt) {
		delete(c.storage, id)
	}
}
