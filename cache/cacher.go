package cache

import (
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/panik"
	"strings"
	"time"
)

type Cacher interface {
	Set(key string, data []byte, expireIn time.Duration)
	Get(key string) ([]byte, error)
}

func getIndex(key string) string {
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

	case "debug":
		out = DebugFromConfig(container)

	default:
		panic("unknown cache type specifed: " + cType)
	}

	return out
}
