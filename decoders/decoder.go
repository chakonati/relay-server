package decoders

type Decoder interface {
	Unmarshal([]byte, interface{}) error
}
