package enc

import (
	"encoding/json"
)

type PrettyJsonEncoder struct {
}

func (e PrettyJsonEncoder) Encode(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "\t")
}

func (e PrettyJsonEncoder) Decode(d []byte) (interface{}, error) {
	var i interface{}
	err := json.Unmarshal(d, &i)
	return i, err
}
