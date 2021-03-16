package decoders

import (
	"encoding/json"
)

type JSON struct{}

func (J JSON) Unmarshal(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}
