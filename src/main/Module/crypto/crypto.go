package crypto

import (
	"crypto"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

/// StackOverflow : REF -> https://stackoverflow.com/users/1705598/icza
func RandStringBytes(n int) string {

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)

}

func Hash(s string) []byte {

	hash := crypto.SHA1.New()

	hash.Write([]byte(s))

	bs := hash.Sum(nil)

	return bs

}

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err

}

func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}
