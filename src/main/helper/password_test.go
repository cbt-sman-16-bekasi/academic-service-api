package helper

import (
	"log"
	"testing"
)

func TestPasswordGenerate(t *testing.T) {
	password := "123456"
	hash, salt, err := HashPasswordArgon2(password)
	if err != nil {
		t.Error(err)
	}

	log.Println(hash, salt)
}
