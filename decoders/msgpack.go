package decoders

import (
	bytes2 "bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type Msgpack struct{}

func (m Msgpack) Unmarshal(bytes []byte, i interface{}) error {
	dec := msgpack.NewDecoder(bytes2.NewReader(bytes))
	dec.SetCustomStructTag("key")
	return dec.Decode(i)
}
