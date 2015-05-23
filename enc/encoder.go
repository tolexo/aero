package enc

import (
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/panik"
)

type Encoder interface {
	Encode(i interface{}) ([]byte, error)
	Decode(d []byte) (interface{}, error)
}

func FromConfig(key string) Encoder {
	eType := conf.String(key, "")
	panik.If(eType == "", "encoding not specified @ "+key)

	return NewEncoder(eType)
}

func NewEncoder(eType string) Encoder {
	switch eType {
	case "message-pack":
		return MessagePackEncoder{}
	case "json":
		return JsonEncoder{}
	case "pretty-json":
		return PrettyJsonEncoder{}
	}
	return nil
}
