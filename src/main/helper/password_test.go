package helper

import (
	"encoding/base64"
	"log"
	"testing"
)

func TestPasswordGenerate(t *testing.T) {
	password := "qwertyuiop"
	salt := "2c+MFslHpRqi1+FOF7eepw=="
	byt, _ := base64.StdEncoding.DecodeString(salt)
	hash, salt, err := HashPasswordArgon2UseSalt(password, byt)
	if err != nil {
		t.Error(err)
	}
	//jITF9qoQ/MRbHqvMQXVY9rvOtRGdO7k+lpRcWWcOKag=
	log.Println(hash, salt)
}
