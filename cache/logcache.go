package cache

import (
	"errors"
	"github.com/thejackrabbit/aero/enc"
	"github.com/thejackrabbit/aero/panik"
	"os"
	"time"
)

func NewLogCache(path string) Cacher {

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	panik.On(err)

	return logCache{
		m: make(map[string]interface{}),
		f: f,
		CacheBase: CacheBase{
			Encoder: enc.NewEncoder("pretty-json"),
		},
	}
}

type logCache struct {
	CacheBase
	m map[string]interface{}
	f *os.File
}

func (c logCache) Set(key string, i interface{}, expireIn time.Duration) {
	j, err := c.Encoder.Encode(i)
	panik.On(err)

	newLine := "\r\n"
	c.f.WriteString(newLine + "~ ~ ~")
	c.f.WriteString(newLine + "SET: " + c.Index(key) + " @ " + time.Now().Format("2006-01-02_15:04:05"))
	c.f.WriteString(newLine + string(j))

	c.m[c.Index(key)] = i
}

func (c logCache) Get(key string) (interface{}, error) {
	if i, ok := c.m[c.Index(key)]; ok {
		return i, nil
	} else {
		return i, errors.New("key not found in cache map")
	}

}
