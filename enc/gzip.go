package enc

// import (
// 	"bytes"
// 	_ "compress/gzip"
// 	"encoding/gob"
// 	"encoding/json"
// )

// type GzipEncoder struct {
// }

// func (e GzipEncoder) Encode(i interface{}) ([]byte, error) {
// 	var b bytes.Buffer
// 	w := gob.NewEncoder(&b)
// 	err := w.Encode(i)
// 	return b, err
// }

// func (e GzipEncoder) Decode(d []byte) (interface{}, error) {

// 	var i interface{}
// 	err := json.Unmarshal(d, &i)
// 	return i, err
// }
