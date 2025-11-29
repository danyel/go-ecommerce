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
}

func (s *shoppingBasketService) CreateShoppingBasket() (ShoppingBasket, error) {
	sb := ShoppingBasketModel{}
	err := s.r.Create(&sb)
	if err != nil {
		return ShoppingBasket{}, err
	}
	r := ShoppingBasket{
		Id: sb.ID,
	}

	return r, nil
}

func (s *shoppingBasketService) AddItemToShoppingBasket(u uuid.UUID, i AddItem) (ShoppingBasket, error) {
	id, err := s.r.FindById(u)
	if err != nil {
		return ShoppingBasket{}, err
	}
	id.Items = append(id.Items, &ShoppingBasketItemModel{ID: i.ProductId})
	if err = s.r.AssocAppend(id, "Items", id.Items); err != nil {
		return ShoppingBasket{}, err
	}
	return s.GetShoppingBasket(u)
}

func (s *shoppingBasketService) GetShoppingBasket(u uuid.UUID) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	if err != nil {
		return ShoppingBasket{}, err
	}
	sm := ShoppingBasket{
		Id: id.ID,
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
	return &shoppingBasketService{r, p, s, m}
}
