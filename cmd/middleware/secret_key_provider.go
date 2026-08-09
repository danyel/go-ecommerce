package middleware

import (
	AES "crypto/aes"
	Cipher "crypto/cipher"
	Rand "crypto/rand"
	Base64 "encoding/base64"
	HEX "encoding/hex"
	JSON "encoding/json"
	Errors "errors"
	Fmt "fmt"
	IO "io"
)

const (
	SECRET_SIZE = 32
)

type SecretKeyProvider interface {
	GenerateKey() (string, error)
}

type secretKeyProvider struct{}

func NewSecretKeyProvider() SecretKeyProvider {
	return &secretKeyProvider{}
}
func (s *secretKeyProvider) GenerateKey() (string, error) {
	bytes := make([]byte, SECRET_SIZE)
	if _, err := Rand.Read(bytes); err != nil {
		return "", err
	}

	secretKey := HEX.EncodeToString(bytes)
	Fmt.Printf("Generated Secure Key: %s", secretKey)
	return secretKey, nil
}

// EncryptClaims serializes and encrypts UserClaims into a URL-safe Base64 token string
func EncryptClaims(claims *UserClaims, secretKey string) (string, error) {
	plainText, err := JSON.Marshal(claims)
	if err != nil {
		return "", Fmt.Errorf("failed to marshal claims: %w", err)
	}

	rawKey, err := HEX.DecodeString(secretKey)

	if err != nil {
		return "", Fmt.Errorf("failed to decode hex secret key: %w", err)
	}

	block, err := AES.NewCipher([]byte(rawKey))
	if err != nil {
		return "", Fmt.Errorf("invalid secret key length: %w", err)
	}

	gcm, err := Cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := IO.ReadFull(Rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	c := Base64.URLEncoding.EncodeToString(cipherText)

	return c, nil
}

// DecryptClaims reverses the process, authenticates the data, and returns the UserClaims
func DecryptClaims(tokenStr string, secretKey string) (*UserClaims, error) {
	cipherText, err := Base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, Errors.New("malformed encryption payload structure")
	}

	rawKey, err := HEX.DecodeString(secretKey)
	if err != nil {
		return nil, Fmt.Errorf("failed to decode hex secret key: %w", err)
	}

	block, err := AES.NewCipher([]byte(rawKey))
	if err != nil {
		return nil, err
	}

	gcm, err := Cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, Errors.New("encrypted payload too short")
	}

	nonce, actualEncryptedData := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, actualEncryptedData, nil)
	if err != nil {
		return nil, Errors.New("tampered or invalid encrypted token signature")
	}

	var claims UserClaims
	if err := JSON.Unmarshal(plainText, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
