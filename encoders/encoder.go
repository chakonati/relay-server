package encoders

type Encoder interface {
	Marshal(interface{}) ([]byte, error)
}
