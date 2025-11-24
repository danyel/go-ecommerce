package product

import (
	"encoding/json"
	"net/http"

	repository "github.com/dnoulet/ecommerce/internal/common"
	"github.com/dnoulet/ecommerce/internal/management"
	util "github.com/dnoulet/ecommerce/internal/util/request"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
type ProductHandler struct {
	ProductRepository  repository.CrudRepository[ProductModel]
	CategoryRepository repository.CrudRepository[management.CategoryModel]
	CategoryMapper     management.CategoryMapper
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, _ *http.Request) {
	orderBy := "created_at asc"
	products := h.ProductRepository.FindAll(repository.SearchCriteria{OrderBy: &orderBy})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(MapToProduct(products, h))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func MapToProduct(productModels []*ProductModel, h *ProductHandler) []Product {
	result := make([]Product, len(productModels))
	for i, product := range productModels {
		category, _ := h.CategoryRepository.FindById(productModels[i].Category)
		result[i] = Product{
			Code:        product.Code,
			Price:       product.Price,
			Category:    h.CategoryMapper.MapCategory(category),
			ImageUrl:    product.ImageUrl,
			Brand:       product.Brand,
			Description: product.Description,
			Name:        product.Name,
			ID:          product.ID,
			Stock:       product.Stock,
		}
	}
	return result
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	parse, err2 := uuid.Parse(productId)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductRepository.FindById(parse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.StatusText(200)
	category, _ := h.CategoryRepository.FindById(product.Category)
	err = json.NewEncoder(w).Encode(Product{
		Code:        product.Code,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    h.CategoryMapper.MapCategory(category),
		ImageUrl:    product.ImageUrl,
		Brand:       product.Brand,
		Description: product.Description,
		Name:        product.Name,
		ID:          product.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	parse, err2 := uuid.Parse(productId)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.ProductRepository.Delete(parse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	parse, err2 := uuid.Parse(productId)
	var updateProduct UpdateProduct
	requestBodyError := util.ValidateRequest(r, &updateProduct)
	if requestBodyError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductRepository.FindById(parse)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product.Name = updateProduct.Name
	product.Brand = updateProduct.Brand
	product.Description = updateProduct.Description
	product.Stock = updateProduct.Stock
	product.Category = updateProduct.Category
	product.ImageUrl = updateProduct.ImageUrl
	product.Price = updateProduct.Price
	dbError := h.ProductRepository.Update(product)

	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var createProduct CreateProduct
	err := util.ValidateRequest(r, &createProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	product := ProductModel{
		Code:        createProduct.Code,
		Price:       createProduct.Price,
		Category:    createProduct.Category,
		ImageUrl:    createProduct.ImageUrl,
		Brand:       createProduct.Brand,
		Description: createProduct.Description,
		Name:        createProduct.Name,
	}
	err = h.ProductRepository.Create(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ProductId{ID: product.ID})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewProductHandler(DB *gorm.DB) ProductHandler {
	return ProductHandler{
		ProductRepository:  repository.NewCrudRepository[ProductModel](DB),
		CategoryRepository: repository.NewCrudRepository[management.CategoryModel](DB),
		CategoryMapper:     management.NewCategoryMapper(),
	}
}
