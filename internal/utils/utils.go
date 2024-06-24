package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/scrypt"

	"github.com/dnsoftware/gophkeeper/internal/constants"
)

// PassGenerate Генерация пароля, возвращает хеш пароля в строковом виде
func PassGenerate(password string, salt string) string {
	saltByte, _ := hex.DecodeString(salt)
	hash, err := scrypt.Key([]byte(password), saltByte, 1<<14, 8, 1, constants.PW_HASH_BYTES)
	if err != nil {
		log.Fatal(err)
	}

	p := fmt.Sprintf("%x", hash)

	return p
}

// SaltGenerate генерация соли для пароля. Возвращает соль в бинарном и строковом виде
func SaltGenerate() ([]byte, string) {
	salt := make([]byte, constants.PW_SALT_BYTES)
	io.ReadFull(rand.Reader, salt)

	return salt, fmt.Sprintf("%x", salt)
}
