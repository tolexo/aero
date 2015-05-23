package enc

import (
	"encoding/json"
)

type JsonEncoder struct {
}

func (e JsonEncoder) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (e JsonEncoder) Decode(d []byte) (interface{}, error) {
	var i interface{}
	err := json.Unmarshal(d, &i)
	return i, err
}
