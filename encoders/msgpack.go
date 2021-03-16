package encoders

import "github.com/vmihailenco/msgpack/v5"

type Msgpack struct{}

func (m Msgpack) Marshal(i interface{}) ([]byte, error) {
	return msgpack.Marshal(i)
}
