package bcrypt

import (
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type CustomBcryptInterface interface {
	Hash(plain string) (string, error)
	Compare(password, hashed string) bool
}

type CustomBcryptStruct struct{}

var Bcrypt = getBcrypt()

func getBcrypt() CustomBcryptInterface {
	return &CustomBcryptStruct{}
}

func (b *CustomBcryptStruct) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[BCRYPT][Hash] failed to hash password")

		return "", err
	}

	return string(bytes), nil
}

func (b *CustomBcryptStruct) Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))

	return err == nil
}
