package enc

import (
	"gopkg.in/vmihailenco/msgpack.v2"
)

type MessagePackEncoder struct {
}

func (e MessagePackEncoder) Encode(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}

func (e MessagePackEncoder) Decode(d []byte) (interface{}, error) {
	var i interface{}
	err := msgpack.Unmarshal(d, &i)
	return i, err
}
