package integration_test

import (
	"testing"

	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/dnoulet/ecommerce/test/integration"
	_ "github.com/lib/pq" // ← REQUIRED for Goose + sql.Open("postgres")
	"github.com/stretchr/testify/assert"
)

func TestCmsBackend(t *testing.T) {
	initializer := integration.NewBackendInitializer()
	initializer.Run()
	cmsRepository := commonRepository.NewCrudRepository[cms.CmsModel](initializer.DB)

	t.Run("Create translation", func(t *testing.T) {
		c := cms.CmsModel{
			Code:     "Code",
			Language: "nl_be",
			Value:    "Value",
		}

		e := cmsRepository.Create(&c)
		assert.Nil(t, e)
	})
}
