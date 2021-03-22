package crypto

import (
	"fmt"
	"log"
	"math/rand"
	"server/decoders"
	"server/encoders"
	"time"

	"github.com/pkg/errors"
	"go.step.sm/crypto/randutil"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const (
	scryptDefaultSaltLen            = 64
	scryptDefaultKeyLen             = 128
	scryptDefaultInnerCPUAndMemCost = 10
	scryptDefaultCPUAndMemCost      = 32768
	scryptParallelization           = 1
)

const (
	scryptMinimumSaltLen   = 16
	scryptMinimumKeyLen    = 32
	scryptMinimumCost      = 6
	scryptMinimumCPUAndMem = 16384
)

type scryptParameters struct {
	CPUAndMemCost      int
	InnerCPUAndMemCost int
	Parallelization    int
	SaltLen            int
	HashLen            int
	KeyLen             int
}

func ScryptHash(byt []byte) ([]byte, error) {
	salt, err := randutil.Salt(scryptDefaultSaltLen)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate salt")
	}
	params := scryptParameters{
		CPUAndMemCost:      scryptDefaultCPUAndMemCost,
		InnerCPUAndMemCost: scryptDefaultInnerCPUAndMemCost,
		Parallelization:    scryptParallelization,
		KeyLen:             scryptDefaultKeyLen,
	}
	key, err := scrypt.Key(byt, salt,
		params.CPUAndMemCost,
		params.InnerCPUAndMemCost,
		params.Parallelization,
		params.KeyLen)
	if err != nil {
		return nil, errors.Wrap(err, "could not derive key")
	}
	hash, err := bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash derived key")
	}

	params.SaltLen = len(salt)
	params.HashLen = len(hash)
	packedParams, err := encoders.Msgpack{}.Pack(&params)
	if err != nil {
		return nil, errors.Wrap(err, "could not pack params")
	}

	combined := make([]byte, 0, 1+len(packedParams)+len(hash)+len(hash))
	combined = append(append(append(append(combined, byte(len(packedParams))), packedParams...), hash...), salt...)
	return combined, err
}

func ScryptCompare(password []byte, hash []byte) (matches bool, err error) {
	if len(password) == 0 {
		err = fmt.Errorf("no password specified")
		return
	}
	defer func() {
		time.Sleep(time.Duration(rand.Intn(1_000))*time.Nanosecond + time.Duration(password[0])*time.Nanosecond)
		if r := recover(); r != nil {
			log.Println(r)
			time.Sleep(time.Duration(rand.Intn(1_000_000)) * time.Nanosecond)
			matches = false
			err = errors.New("comparing scrypt hash failed miserably")
		}
	}()
	packedParamsLen := int(hash[0])
	if packedParamsLen == 0 {
		err = fmt.Errorf("packed params too short")
		return
	}
	packedParams := hash[1 : packedParamsLen+1]
	var params scryptParameters
	err = decoders.Msgpack{}.Unmarshal(packedParams, &params)
	if err != nil {
		err = errors.Wrap(err, "could not unpack packed params")
		return
	}

	saltLen := params.SaltLen
	bcryptLen := params.HashLen
	if saltLen < scryptMinimumSaltLen {
		err = fmt.Errorf("salt len of hash is smaller than minimum of %d", scryptMinimumSaltLen)
		return
	}
	if bcryptLen == 0 {
		err = fmt.Errorf("missing bcrypt hash")
		return
	}

	bcryptHash := hash[1+packedParamsLen : bcryptLen+packedParamsLen+1]
	salt := hash[1+packedParamsLen+bcryptLen:]

	key, err := scrypt.Key(password, salt,
		params.CPUAndMemCost, params.InnerCPUAndMemCost, params.Parallelization, params.KeyLen)
	if err != nil {
		err = errors.Wrap(err, "could not derive key")
	}

	resultErr := bcrypt.CompareHashAndPassword(bcryptHash, key)
	matches = resultErr == nil
	return
}
