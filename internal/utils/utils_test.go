package utils

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassGenerate(t *testing.T) {
	saltHash := "cfc04ad973ebee11d59db2ac3750f20d"

	passHash := PassGenerate("password", saltHash)

	assert.Equal(t, "a19301f767b4a89bf81d0b15f3723ce737ee76dad0f11063", passHash)
}

func TestSaltGenerate(t *testing.T) {
	salt, saltHash := SaltGenerate()

	saltByte, _ := hex.DecodeString(saltHash)
	assert.Equal(t, salt, saltByte)
}
