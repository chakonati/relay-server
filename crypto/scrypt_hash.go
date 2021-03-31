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

type scryptHash struct {
	CPUAndMemCost      int
	InnerCPUAndMemCost int
	Parallelization    int
	KeyLen             int
	Salt               []byte
	BcryptHash         []byte
}

func ScryptHash(byt []byte) ([]byte, error) {
	salt, err := randutil.Salt(scryptDefaultSaltLen)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate salt")
	}
	scryptHash := scryptHash{
		CPUAndMemCost:      scryptDefaultCPUAndMemCost,
		InnerCPUAndMemCost: scryptDefaultInnerCPUAndMemCost,
		Parallelization:    scryptParallelization,
		KeyLen:             scryptDefaultKeyLen,
	}
	key, err := scrypt.Key(byt, salt,
		scryptHash.CPUAndMemCost,
		scryptHash.InnerCPUAndMemCost,
		scryptHash.Parallelization,
		scryptHash.KeyLen)
	if err != nil {
		return nil, errors.Wrap(err, "could not derive key")
	}
	scryptHash.BcryptHash, err = bcrypt.GenerateFromPassword(key, bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash derived key")
	}
	scryptHash.Salt = salt

	packed, err := encoders.Msgpack{}.Pack(&scryptHash)
	if err != nil {
		return nil, errors.Wrap(err, "could not pack params")
	}

	return packed, err
}

func ScryptCompare(password []byte, hash []byte) (matches bool, err error) {
	if len(password) == 0 {
		err = fmt.Errorf("no password specified")
		return
	}
	defer func() {
		time.Sleep(time.Duration(rand.Intn(10_000)) * time.Nanosecond)
		if r := recover(); r != nil {
			log.Println(r)
			time.Sleep(time.Duration(rand.Intn(1_000_000)) * time.Nanosecond)
			matches = false
			err = errors.New("comparing scrypt hash failed miserably")
		}
	}()
	var scryptHash scryptHash
	err = decoders.Msgpack{}.Unmarshal(hash, &scryptHash)
	if err != nil {
		err = errors.Wrap(err, "could not unmarshal scrypt hash")
		return
	}

	saltLen := len(scryptHash.Salt)
	bcryptLen := len(scryptHash.BcryptHash)
	if saltLen < scryptMinimumSaltLen {
		err = fmt.Errorf("salt len of hash is smaller than minimum of %d", scryptMinimumSaltLen)
		return
	}
	if bcryptLen == 0 {
		err = fmt.Errorf("missing bcrypt hash")
		return
	}

	key, err := scrypt.Key(password, scryptHash.Salt,
		scryptHash.CPUAndMemCost,
		scryptHash.InnerCPUAndMemCost,
		scryptHash.Parallelization,
		scryptHash.KeyLen,
	)
	if err != nil {
		err = errors.Wrap(err, "could not derive key")
	}

	resultErr := bcrypt.CompareHashAndPassword(scryptHash.BcryptHash, key)
	matches = resultErr == nil
	return
}
