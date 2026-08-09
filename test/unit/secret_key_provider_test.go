package unit

import (
	Testing "testing"

	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	Assert "github.com/stretchr/testify/assert"
)

func TestSecretKeyProvider(t *Testing.T) {
	ut := ApplicationMiddleware.NewSecretKeyProvider()

	t.Run("Generating key", func(t *Testing.T) {
		key, err := ut.GenerateKey()
		Assert.NotNil(t, "Error generating key", err)
		Assert.NotEmpty(t, "Key cab not be empty", key)
	})
}
