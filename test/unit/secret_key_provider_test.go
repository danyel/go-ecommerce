package unit

import (
	Fmt "fmt"
	Testing "testing"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	Assert "github.com/stretchr/testify/assert"
)

func TestSecretKeyProvider(unitTest *Testing.T) {
	secretKeyProvider := ApplicationMiddleware.NewSecretKeyProvider()

	unitTest.Run("Generating key", func(unitTest *Testing.T) {
		key, err := secretKeyProvider.GenerateKey()
		Assert.NotNil(unitTest, "Error generating key", err)
		Assert.NotEmpty(unitTest, "Key can not be empty", key)
	})

	unitTest.Run("Encrypt User Claims", func(unitTest *Testing.T) {
		secret := "b87d91cbfabed9283fdb389c032f158730b2c9810200ef24eaabb82fb49a2a9c"
		Fmt.Printf("Secret size: %d", len(secret))
		Assert.NotEmpty(unitTest, "Key can not be empty", secret)
		token, err := ApplicationMiddleware.EncryptClaims(&ApplicationMiddleware.UserClaims{
			UserID: "userId",
			Roles:  []string{"ADMIN", "USER"},
		}, secret)
		Assert.NotNil(unitTest, "Error generating key", err)
		Fmt.Printf("Token: %s", token)
		decryptedUserClaims, err := ApplicationMiddleware.DecryptClaims(token, secret)
		Assert.NotNil(unitTest, "Error generating key", err)
		Assert.Equal(unitTest, "userId", decryptedUserClaims.UserID)
	})
}
