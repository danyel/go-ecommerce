package category

import (
	"encoding/json"
	"net/http"

	requestUtil "github.com/dnoulet/ecommerce/internal/util/request"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type CategoryHandler struct {
	CreateCategory     func(w http.ResponseWriter, r *http.Request)
	CreateTranslations func(_ http.ResponseWriter, _ *http.Request)
}

func NewHandler(DB *gorm.DB) CategoryHandler {
	categoryService := NewCategoryService(DB)
	return CategoryHandler{
		CreateCategory: func(w http.ResponseWriter, r *http.Request) {
			var createCategory CreateCategory
			var categoryId CategoryId
			var err error
			if err = requestUtil.ValidateRequest(r, &createCategory); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if categoryId, err = categoryService.CreateCategory(createCategory); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(categoryId); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
		CreateTranslations: func(_ http.ResponseWriter, _ *http.Request) {},
	}
}
