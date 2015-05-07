package cache

import (
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/enc"
	"github.com/thejackrabbit/aero/panik"
	"strings"
	"time"
)

type Cacher interface {
	Set(key string, i interface{}, expireIn time.Duration)
	Get(key string) (interface{}, error)
}

type CacheBase struct {
	Encoder enc.Encoder
}

func (c CacheBase) Index(key string) string {
	return strings.Replace(key, " ", "-", -1)
}

func FromConfig(container string) (out Cacher) {

	cType := conf.String(container+".type", "")
	panik.If(cType == "", "cache type is not specified")

	switch cType {
	case "memcache":
		out = MemcacheFromConfig(container)

	case "inmem":
		out = InmemFromConfig(container)

	default:
		panic("unknown cache type specifed")
	}

	return out
}
