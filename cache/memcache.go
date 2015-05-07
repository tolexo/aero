package cache

import (
	"github.com/kklis/gomemcache"
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/enc"
	"github.com/thejackrabbit/aero/panik"
	"time"
)

func NewMemcache(host string, port int, enc enc.Encoder) Cacher {

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

// memcache:
// - host
// - port
// - encoding
func MemcacheFromConfig(container string) Cacher {
	host := conf.String(container+".host", "")
	panik.If(host == "", "memcache host not specified")

	port := conf.Int(container+".port", 0)
	panik.If(port == 0, "memcache port not specified")

	encd := enc.FromConfig(container + ".encoding")
	panik.If(encd == nil, "memcache encoding not specified")

	return NewMemcache(host, port, encd)
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
	} else {
		return c.Encoder.Decode(j)
	}
}
