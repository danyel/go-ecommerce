package management

import (
	"encoding/json"
	"net/http"

	repository "github.com/dnoulet/ecommerce/internal/common"
	requestUtil "github.com/dnoulet/ecommerce/internal/util/request"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ManagementHandler struct {
	CategoryRepository repository.CrudRepository[CategoryModel]
	CategoryMapper     CategoryMapper
}

func (h *ManagementHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var createCategory CreateCategory
	if err := requestUtil.ValidateRequest(r, &createCategory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 1️⃣ Create parent
	category := &CategoryModel{
		Name: createCategory.Name,
	}
	if err := h.CategoryRepository.Create(category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 2️⃣ Load existing children from DB
	var children []*CategoryModel
	if len(createCategory.Children) > 0 {
		children = h.CategoryRepository.FindAll(repository.SearchCriteria{
			WhereClause: repository.WhereClause{
				Query:  "id IN ?",
				Params: []interface{}{createCategory.Children},
			},
		})
	}

	// 3️⃣ Append children WITHOUT triggering Create on them
	if len(children) > 0 {
		if err := h.CategoryRepository.DB.
			Session(&gorm.Session{FullSaveAssociations: false}). // important: prevents saving children
			Model(category).
			Association("Children").
			Append(children); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// 4️⃣ Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(CategoryId{ID: category.ID}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ManagementHandler) GetCategories(w http.ResponseWriter, _ *http.Request) {
	categoryModels := h.CategoryRepository.FindAll(repository.SearchCriteria{Preloads: []string{"Children"}})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(h.CategoryMapper.MapCategories(categoryModels))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ManagementHandler) CreateTranslations(_ http.ResponseWriter, _ *http.Request) {}

func NewManagementHandler(DB *gorm.DB) ManagementHandler {
	return ManagementHandler{
		CategoryRepository: repository.NewCrudRepository[CategoryModel](DB),
		CategoryMapper:     NewCategoryMapper(),
	}
}
