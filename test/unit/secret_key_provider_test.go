package unit

import (
	Fmt "fmt"
	Testing "testing"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	Assert "github.com/stretchr/testify/assert"
)

func TestSecretKeyProvider(t *Testing.T) {
	ut := ApplicationMiddleware.NewSecretKeyProvider()

	t.Run("Generating key", func(t *Testing.T) {
		key, err := ut.GenerateKey()
		Assert.NotNil(t, "Error generating key", err)
		Assert.NotEmpty(t, "Key can not be empty", key)
	})

	t.Run("Encrypt User Claims", func(t *Testing.T) {
		secret := "b87d91cbfabed9283fdb389c032f158730b2c9810200ef24eaabb82fb49a2a9c"
		Fmt.Printf("Secret size: %d", len(secret))
		Assert.NotEmpty(t, "Key can not be empty", secret)
		token, err := ApplicationMiddleware.EncryptClaims(&ApplicationMiddleware.UserClaims{
			UserID: "userId",
			Roles:  []string{"ADMIN", "USER"},
		}, secret)
		Assert.NotNil(t, "Error generating key", err)
		Fmt.Printf("Token: %s", token)
		decryptedUserClaims, err := ApplicationMiddleware.DecryptClaims(token, secret)
		Assert.NotNil(t, "Error generating key", err)
		Assert.Equal(t, "userId", decryptedUserClaims.UserID)
	})
}
