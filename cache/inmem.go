package cache

import (
	"errors"
	goc "github.com/pmylund/go-cache"
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/enc"
	"github.com/thejackrabbit/aero/panik"
	"os"
	"time"
)

func NewInmemCache(enc enc.Encoder, r string, w string) Cacher {

	var err error
	var rf, wf *os.File

	if r != "" {
		rf, err = os.OpenFile(r, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		panik.On(err)
	}
	if w != "" {
		wf, err = os.OpenFile(w, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		panik.On(err)
	}

	return inmemCache{
		c: goc.New(3*time.Minute, 1*time.Minute),
		CacheBase: CacheBase{
			Encoder: enc,
		},
		r: rf,
		w: wf,
	}
}

// inmem
// - encoding
// - read-file
// - write-file
func InmemFromConfig(container string) Cacher {
	encd := enc.FromConfig(container + ".encoding")
	panik.If(encd == nil, "inmem encoding not specified")

	return NewInmemCache(encd,
		conf.String(container+".read-file", ""),
		conf.String(container+".write-file", ""),
	)
}

// Stores encoded data in memory (map). Key purpose is to be an aid for developers
// during development/debugging phase. All reads, writes can be optionally written
// to specified files.
type inmemCache struct {
	CacheBase
	c *goc.Cache
	r *os.File
	w *os.File
}

func (c inmemCache) Set(key string, i interface{}, expireIn time.Duration) {
	j, err := c.Encoder.Encode(i)
	panik.On(err)

	k := c.Index(key)
	c.c.Set(k, j, expireIn)

	go c.logSet(k, j)
}

func (c inmemCache) Get(key string) (interface{}, error) {
	k := c.Index(key)
	if j, ok := c.c.Get(k); ok {
		b := j.([]byte)
		go c.logGet(k, b)
		return c.CacheBase.Encoder.Decode(b)
	} else {
		go c.logGet(k, nil)
		return nil, errors.New("key not found in inmem cache")
	}
}

func (c inmemCache) logGet(k string, b []byte) {
	if c.r == nil {
		return
	}

	newLine := "\r\n"
	c.r.WriteString(newLine + "~ ~ ~")
	c.r.WriteString(newLine + "GET: " + k + " @ " + time.Now().Format("2006-01-02_15:04:05"))

	if b != nil {
		c.r.WriteString(newLine + string(b))
	} else {
		c.r.WriteString(newLine + "<not-found>")
	}
}

func (c inmemCache) logSet(k string, b []byte) {
	if c.w == nil {
		return
	}

	newLine := "\r\n"
	c.w.WriteString(newLine + "~ ~ ~")
	c.w.WriteString(newLine + "SET: " + k + " @ " + time.Now().Format("2006-01-02_15:04:05"))

	if b != nil {
		c.w.WriteString(newLine + string(b))
	} else {
		c.w.WriteString(newLine + "<nil>")
	}
}
