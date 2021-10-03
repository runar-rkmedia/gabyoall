package bboltStorage

import (
	"bytes"
	"encoding/gob"
)

func (g Gob) Marshal(j interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(j)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (g Gob) Unmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(v)
}

type Gob struct{}

type Marshaller interface {
	Unmarshal(data []byte, v interface{}) error
	Marshal(j interface{}) ([]byte, error)
}
