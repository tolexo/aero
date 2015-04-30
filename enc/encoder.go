package enc

import (
	"encoding/json"
	"github.com/thejackrabbit/aero/conf"
	"github.com/thejackrabbit/aero/panik"
	"github.com/vmihailenco/msgpack"
)

type Encoder interface {
	Encode(i interface{}) ([]byte, error)
	Decode(d []byte) (interface{}, error)
}

func List() []string {
	return []string{"message-pack", "json", "pretty-json"}
}

func FromConfig(key string) Encoder {
	eType := conf.String(key, "")
	panik.If(eType == "", "encoding not specified @ "+key)

	return NewEncoder(eType)
}

func NewEncoder(eType string) Encoder {
	switch eType {
	case "message-pack":
		return messagePackEncoder{}
	case "json":
		return jsonEncoder{}
	case "pretty-json":
		return prettyJsonEncoder{}
	}

	return nil
}

// MessagePack Encoder
type messagePackEncoder struct {
}

func (e messagePackEncoder) Encode(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func (e messagePackEncoder) Decode(d []byte) (interface{}, error) {
	var i interface{}
	err := msgpack.Unmarshal(d, &i)
	return i, err
}

// Json Encoder (standard)
type jsonEncoder struct {
}

func (e jsonEncoder) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (e jsonEncoder) Decode(d []byte) (interface{}, error) {
	var i interface{}
	err := json.Unmarshal(d, &i)
	return i, err
}

// Pretty Json Encoder
type prettyJsonEncoder struct {
}

func (e prettyJsonEncoder) Encode(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "\t")
}

func (e prettyJsonEncoder) Decode(d []byte) (interface{}, error) {
	var i interface{}
	err := json.Unmarshal(d, &i)
	return i, err
}
