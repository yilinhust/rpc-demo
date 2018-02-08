package encrypt

import (
	"crypto"
	"crypto/hmac"
	"errors"
)


var HmacKey []byte =  []byte("&#ea<!")

type Algorithm interface {
	Encrypt(data string) (string, error)
	Verify(data, signature string) error
}

type HMACAlgorithm struct {
	hash crypto.Hash
	key  []byte
}

func NewHMACAlgorithm(hash crypto.Hash, key []byte) (*HMACAlgorithm, error) {
	algorithm := &HMACAlgorithm{
		hash: hash,
		key:  key,
	}

	if !algorithm.hash.Available() {
		return nil, errors.New("The requested hash function is unavailable")
	}

	return algorithm, nil
}


// Implements the Encrypt method from Algorithm.
func (h *HMACAlgorithm) Encrypt(data string) (string, error) {
	hasher := hmac.New(h.hash.New, h.key)
	hasher.Write([]byte(data))
	return Encode(hasher.Sum(nil)), nil
}

// Implements the Verify method from Algorithm.
func (h *HMACAlgorithm) Verify(data, signature string) error {
	sig, err := Decode(signature)
	if err != nil {
		return err
	}
	hasher := hmac.New(h.hash.New, h.key)
	hasher.Write([]byte(data))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return errors.New("Signature is invalid")
	}
	return nil
}