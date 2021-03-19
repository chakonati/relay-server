package defs

type PreKeyBundle struct {
	RegistrationID        int
	DeviceID              int
	PreKeyID              int
	PublicPreKey          []byte
	SignedPreKeyID        int
	PublicSignedPreKey    []byte
	SignedPreKeySignature []byte
	IdentityKey           []byte
}
