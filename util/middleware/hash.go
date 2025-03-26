package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	hashingTime uint32 = 1         // Jumlah iterasi untuk hashing
	memory      uint32 = 64 * 1024 // Jumlah memori yang digunakan (dalam KB)
	threads     uint8  = 4         // Jumlah parallelism
	keyLength   uint32 = 32        // Panjang hasil hash
	saltSize    int    = 16        // Ukuran salt dalam byte
)

func GenerateSalt(size int) (string, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func HashPassword(password, salt string) (string, error) {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), saltBytes, hashingTime, memory, threads, keyLength)
	return fmt.Sprintf("%s:%s", salt, base64.RawStdEncoding.EncodeToString(hash)), nil
}

func VerifyPassword(password, hashWithSalt string) (bool, error) {
	parts := strings.Split(hashWithSalt, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("hash is not in the correct format")
	}
	salt, _ := parts[0], parts[1]
	generatedHash, err := HashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return generatedHash == hashWithSalt, nil
}
