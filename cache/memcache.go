package cache

import (
	"github.com/kklis/gomemcache"
	"github.com/thejackrabbit/aero/enc"
	"github.com/thejackrabbit/aero/panik"
	"time"
)

func NewMemcacheCache(host string, port int, enc enc.Encoder) Cacher {

	serv, err := gomemcache.Connect(host, port)
	panik.On(err)

	m := memcacheCache{
		mc: serv,
		CacheBase: CacheBase{
			Encoder: enc,
		},
	}
	return m
}

type memcacheCache struct {
	CacheBase
	mc *gomemcache.Memcache
}

func (c memcacheCache) Set(key string, i interface{}, expireIn time.Duration) {
	j, err := c.Encoder.Encode(i)
	panik.On(err)
	c.mc.Set(c.Index(key), j, 0, int64(expireIn.Seconds()))
}

func (c memcacheCache) Get(key string) (interface{}, error) {
	j, _, err := c.mc.Get(c.Index(key))
	if err != nil {
		return nil, err
	}

	return c.Encoder.Decode(j)
}
