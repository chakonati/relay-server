package crypto

import (
	"bytes"
	"testing"

	"go.step.sm/crypto/randutil"
)

func TestScrypt(t *testing.T) {
	password, err := randutil.Alphanumeric(24)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := ScryptHash([]byte(password))
	if err != nil {
		t.Fatal(err)
	}

	matches, err := ScryptCompare([]byte(password), hash)
	if err != nil {
		t.Fatal(err)
	}
	if !matches {
		t.Fatal("password doesn't match hash but it should")
	}

	matches, err = ScryptCompare([]byte("not matching"), hash)
	if err != nil {
		t.Fatal(err)
	}
	if matches {
		t.Fatal("password matches hash but it shouldn't")
	}

	differentHash, err := ScryptHash([]byte("something different"))
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(hash, differentHash) {
		t.Fatal("Generated same hash for different inputs!")
	}
}
