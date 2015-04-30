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

	// read encoding if specified
	var encd enc.Encoder
	if conf.Exists(container + ".encoding") {
		encd = enc.FromConfig(container + ".encoding")
	}

	switch cType {
	case "memcache":
		// memcache:
		// - host
		// - port
		// - encoding
		host := conf.String(container+".host", "")
		port := conf.Int(container+".port", 0)
		panik.If(host == "", "memcache host not specified")
		panik.If(port == 0, "memcache port not specified")
		panik.If(encd == nil, "memcache encoding not specified")
		out = NewMemcacheCache(host, port, encd)

	case "log-cache":
		// log-cache
		// - path
		path := conf.String(container+".path", "")
		panik.If(path == "", "log-cache path not specified")
		out = NewLogCache(path)

	default:
		panic("unknown cache type specifed")
	}

	return out
}
