package management

import (
	"encoding/json"
	"net/http"

	"github.com/dnoulet/ecommerce/internal/category"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler interface {
	GetCategories(w http.ResponseWriter, _ *http.Request)
	CreateTranslations(w http.ResponseWriter, _ *http.Request)
}

type hiddenManagementHandler struct {
	categoryService category.CategoryService
}

func (h *hiddenManagementHandler) GetCategories(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(h.categoryService.GetCategories())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *hiddenManagementHandler) CreateTranslations(_ http.ResponseWriter, _ *http.Request) {}

func NewHandler(DB *gorm.DB) ManagementHandler {
	h := &hiddenManagementHandler{
		categoryService: category.NewCategoryService(DB),
	}
	return h
}
