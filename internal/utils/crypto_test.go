package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymmPassCreate(t *testing.T) {
	key := SymmPassCreate("12345678912345678912345678", "tailpfsecretstring")
	assert.Equal(t, "12345678912345678912345678tailpf", key)

	key = SymmPassCreate("12345678", "tail1234")
	assert.Equal(t, "12345678tail1234tail1234tail1234", key)

}

func TestEncryption(t *testing.T) {
	key := SymmPassCreate("12345678912345678912345678", "tailpfsecretstring")
	str := "string to crypting"

	cipher := Encrypt(str, key)
	decipher := Decrypt(cipher, key)

	assert.Equal(t, str, decipher)

}

func TestBinaryEncryption(t *testing.T) {
	key := SymmPassCreate("12345678912345678912345678", "tailpfsecretstring")
	str := "string to crypting"

	cipher := EncryptBinary([]byte(str), key)
	decipher := DecryptBinary(cipher, key)

	assert.Equal(t, str, string(decipher))

}
