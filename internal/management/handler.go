package management

import (
	"encoding/json"
	"net/http"

	"github.com/dnoulet/ecommerce/internal/category"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler struct {
	GetCategories      func(w http.ResponseWriter, _ *http.Request)
	CreateTranslations func(_ http.ResponseWriter, _ *http.Request)
}

func NewHandler(DB *gorm.DB) ManagementHandler {
	categoryService := category.NewCategoryService(DB)
	return ManagementHandler{
		GetCategories: func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(categoryService.GetCategories())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
		CreateTranslations: func(_ http.ResponseWriter, _ *http.Request) {},
	}
}
