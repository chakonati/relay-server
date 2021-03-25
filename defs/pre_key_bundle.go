package defs

type PreKeyBundle struct {
	RegistrationID        int
	DeviceID              int
	SignedPreKeyID        int
	PublicSignedPreKey    []byte
	SignedPreKeySignature []byte
	IdentityKey           []byte
}
