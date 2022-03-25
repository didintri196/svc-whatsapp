package utils

import "golang.org/x/crypto/bcrypt"

type HashHelper struct{}

func NewHashHelper() HashHelper {
	return HashHelper{}
}

func (helper HashHelper) HashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hash), err
}

func (helper HashHelper) CheckHashString(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
