package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHashFn func(string) (string, error)
type PasswordCompareFn func(string, string) error

func BcryptHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func BcryptCompare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
