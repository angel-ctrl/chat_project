package utils

import "golang.org/x/crypto/bcrypt"

func EncriptPass(pass string) (string, error) {

	costo := 8

	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), costo)
	
	return string(bytes), err
}
