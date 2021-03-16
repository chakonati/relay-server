package encoders

import "encoding/json"

type JSON struct{}

func (J JSON) Marshal(i interface{}) ([]byte, error) {

	return json.Marshal(i)
}
