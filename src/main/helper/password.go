package helper

import (
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

var salt = "2c+MFslHpRqi1+FOF7eepw=="

func HashPasswordArgon2(password string) (string, string, error) {
	hashed := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	hashBase64 := base64.StdEncoding.EncodeToString(hashed)
	saltBase64 := base64.StdEncoding.EncodeToString([]byte(salt))

	return hashBase64, saltBase64, nil
}

func HashPasswordArgon2UseSalt(password string, salt []byte) (string, string, error) {
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
