package cache

import (
	"errors"
	goc "github.com/pmylund/go-cache"
	"time"
)

// Stores data in memory. Major use-case is for use on development machines
type inmemCache struct {
	ram *goc.Cache
}

func NewInmemCache() Cacher {
	return inmemCache{
		ram: goc.New(60*time.Minute, 1*time.Minute),
	}
}

// inmem
func InmemFromConfig(container string) Cacher {
	return NewInmemCache()
}

func (c inmemCache) Set(key string, data []byte, expireIn time.Duration) {
	key = getIndex(key)
	c.ram.Set(key, data, expireIn)
}

func (c inmemCache) Get(key string) ([]byte, error) {
	key = getIndex(key)
	if i, ok := c.ram.Get(key); ok {
		return i.([]byte), nil
	} else {
		return nil, errors.New("key not found in inmem cache: " + key)
	}
}
