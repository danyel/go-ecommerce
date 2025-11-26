package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dnoulet/ecommerce/internal/cms"
	commonRepository "github.com/dnoulet/ecommerce/internal/common/repository"
	"github.com/stretchr/testify/assert"
)

func TestCmsHandler(t *testing.T) {
	app := SetupTestApp(t)
	f := NewFixture(commonRepository.NewCrudRepository[cms.CmsModel](app.DB))

	t.Run("Initial Data", func(t *testing.T) {
		f.SaveModel(&cms.CmsModel{Code: "code", Language: "nl_be", Value: "Value_nl"})
		f.SaveModel(&cms.CmsModel{Code: "code", Language: "nl_fr", Value: "Value_fr"})
		f.SaveModel(&cms.CmsModel{Code: "another_code", Language: "nl_be", Value: "AnotherValue_nl"})
		f.SaveModel(&cms.CmsModel{Code: "another_code", Language: "nl_fr", Value: "AnotherValue_fr"})
		f.SaveModel(&cms.CmsModel{Code: "yet_another_code", Language: "nl_be", Value: "YetAnotherValue_nl"})
		f.SaveModel(&cms.CmsModel{Code: "yet_another_code", Language: "nl_fr", Value: "YetAnotherValue_fr"})
	})

	t.Run("TestCmsHandler", func(t *testing.T) {
		t.Run("CmsHandler retrieve dutch", func(t *testing.T) {
			resp, err := http.Get(app.Server.URL + "/api/cms/v1/translations?language=nl_be")

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var translations []cms.Translation
			err = json.NewDecoder(resp.Body).Decode(&translations)
			assert.Nil(t, err)
			assert.Equal(t, 3, len(translations))
			assert.Equal(t, "code", translations[0].Code)
			assert.Equal(t, "Value_nl", translations[0].Value)
			assert.Equal(t, "another_code", translations[1].Code)
			assert.Equal(t, "AnotherValue_nl", translations[1].Value)
			assert.Equal(t, "yet_another_code", translations[2].Code)
			assert.Equal(t, "YetAnotherValue_nl", translations[2].Value)
		})

		t.Run("CmsHandler retrieve french", func(t *testing.T) {
			resp, err := http.Get(app.Server.URL + "/api/cms/v1/translations?language=nl_fr")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var translations []cms.Translation
			err = json.NewDecoder(resp.Body).Decode(&translations)
			assert.Nil(t, err)
			assert.Equal(t, 3, len(translations))
			assert.Equal(t, "code", translations[0].Code)
			assert.Equal(t, "Value_fr", translations[0].Value)
			assert.Equal(t, "another_code", translations[1].Code)
			assert.Equal(t, "AnotherValue_fr", translations[1].Value)
			assert.Equal(t, "yet_another_code", translations[2].Code)
			assert.Equal(t, "YetAnotherValue_fr", translations[2].Value)
		})

		t.Run("CmsHandler retrieve all", func(t *testing.T) {
			resp, err := http.Get(app.Server.URL + "/api/cms/v1/translations")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var translations []cms.Translation
			err = json.NewDecoder(resp.Body).Decode(&translations)
			assert.Nil(t, err)
			assert.Equal(t, 6, len(translations))
			assert.Equal(t, "code", translations[0].Code)
			assert.Equal(t, "Value_nl", translations[0].Value)
			assert.Equal(t, "code", translations[1].Code)
			assert.Equal(t, "Value_fr", translations[1].Value)
			assert.Equal(t, "another_code", translations[2].Code)
			assert.Equal(t, "AnotherValue_nl", translations[2].Value)
			assert.Equal(t, "another_code", translations[3].Code)
			assert.Equal(t, "AnotherValue_fr", translations[3].Value)
			assert.Equal(t, "yet_another_code", translations[4].Code)
			assert.Equal(t, "YetAnotherValue_nl", translations[4].Value)
			assert.Equal(t, "yet_another_code", translations[5].Code)
			assert.Equal(t, "YetAnotherValue_fr", translations[5].Value)
		})

		t.Run("CmsHandler retrieve none because of invalid language", func(t *testing.T) {
			resp, err := http.Get(app.Server.URL + "/api/cms/v1/translations?language=nl_de")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			var translations []cms.Translation
			err = json.NewDecoder(resp.Body).Decode(&translations)
			assert.Nil(t, err)
			assert.Equal(t, 0, len(translations))
		})
	})
}
