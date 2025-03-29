package helper

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPasswordArgon2(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hashBase64 := base64.StdEncoding.EncodeToString(hashed)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	return hashBase64, saltBase64, nil
}

func VerifyPasswordArgon2(password, hashBase64, saltBase64 string) bool {
	salt, _ := base64.StdEncoding.DecodeString(saltBase64)
	hashed := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hashed) == hashBase64
}
