package encoders

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type Msgpack struct{}

func (m Msgpack) Marshal(i interface{}) ([]byte, error) {
	var byt bytes.Buffer
	enc := msgpack.NewEncoder(&byt)
	enc.SetCustomStructTag("key")
	err := enc.Encode(i)
	if err != nil {
		return nil, err
	}
	return byt.Bytes(), nil
}
