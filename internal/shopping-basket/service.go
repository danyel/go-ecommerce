package shopping_basket

import (
	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/danyel/ecommerce/internal/product"
	productmanagement "github.com/danyel/ecommerce/internal/product-management"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShoppingBasketService interface {
	CreateShoppingBasket() (ShoppingBasket, error)
	AddItemToShoppingBasket(u uuid.UUID, i AddItem) (ShoppingBasket, error)
	GetShoppingBasket(u uuid.UUID) (ShoppingBasket, error)
}

type shoppingBasketService struct {
	r  commonRepository.CrudRepository[ShoppingBasketModel]
	p  product.ProductService
	pm productmanagement.ProductService
	m  product.ProductMapper
	si commonRepository.CrudRepository[ShoppingBasketItemModel]
}

func (s *shoppingBasketService) CreateShoppingBasket() (ShoppingBasket, error) {
	sb := ShoppingBasketModel{}
	err := s.r.Create(&sb)
	if err != nil {
		return ShoppingBasket{}, err
	}
	r := ShoppingBasket{
		ID: sb.ID,
	}

	return r, nil
}

func (s *shoppingBasketService) AddItemToShoppingBasket(u uuid.UUID, i AddItem) (ShoppingBasket, error) {
	id, err := s.r.FindById(u)
	var prd product.Product
	if err != nil {
		return ShoppingBasket{}, err
	}
	if prd, err = s.p.GetProduct(i.ProductId); err != nil {
		return ShoppingBasket{}, err
	}

	item := ShoppingBasketItemModel{ID: prd.ID, ShoppingBasketID: id.ID, Name: prd.Name, Price: prd.Price, ImageUrl: prd.ImageUrl, Amount: 1}
	for _, it := range id.Items {
		if it.ID == i.ProductId {
			item.Amount = it.Amount + 1
		}
	}
	err = s.si.Create(&item)
	return s.GetShoppingBasket(u)
}

func (s *shoppingBasketService) GetShoppingBasket(u uuid.UUID) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	if err != nil {
		return ShoppingBasket{}, err
	}
	sm := ShoppingBasket{
		ID: id.ID,
	}
	if len(id.Items) > 0 {
		ps := make([]ShoppingBasketItem, len(id.Items))
		for i, item := range id.Items {
			ps[i] = ShoppingBasketItem{
				item.ID, item.Name, item.Price, item.ImageUrl, item.Amount,
			}
		}
		sm.Items = ps
	}
	return sm, nil
}

func NewService(db *gorm.DB) ShoppingBasketService {
	r := commonRepository.NewCrudRepository[ShoppingBasketModel](db)
	p := product.NewProductService(db)
	s := productmanagement.NewProductService(db)
	m := product.NewProductMapper(category.NewCategoryService(db), cms.NewCmsService(db))
	si := commonRepository.NewCrudRepository[ShoppingBasketItemModel](db)
	return &shoppingBasketService{r, p, s, m, si}
}
