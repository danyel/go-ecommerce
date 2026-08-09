package middleware

import (
	Rand "crypto/rand"
	HEX "encoding/hex"
	"errors"
	Fmt "fmt"
)

type SecretKeyProvider interface {
	GenerateKey() (string, error)
}

type secretKeyProvider struct{}

func NewSecretKeyProvider() SecretKeyProvider {
	return &secretKeyProvider{}
}
func (s *secretKeyProvider) GenerateKey() (string, error) {
	bytes := make([]byte, 64)
	if _, err := Rand.Read(bytes); err != nil {
		return "", err
	}

	secretKey := HEX.EncodeToString(bytes)
	Fmt.Printf("Generated Secure Key: %s", secretKey)
	return secretKey, errors.New("eroro")
}
