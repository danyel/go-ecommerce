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
	UpdateShoppingBasketItem(u uuid.UUID, i UpdateShoppingBasketItem) (ShoppingBasket, error)
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

func (s *shoppingBasketService) UpdateShoppingBasketItem(u uuid.UUID, i UpdateShoppingBasketItem) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	var prd product.Product
	if err != nil {
		return ShoppingBasket{}, err
	}
	if prd, err = s.p.GetProduct(i.ProductId); err != nil {
		return ShoppingBasket{}, err
	}

	item := ShoppingBasketItemModel{ID: uuid.Nil, ShoppingBasketID: id.ID, ProductId: prd.ID, Price: prd.Price, Quantity: i.Quantity}
	for _, it := range id.Items {
		if it.ProductId == item.ProductId {
			item.ID = it.ID
			item.Quantity = i.Quantity
		}
	}
	if item.ID == uuid.Nil {
		err = s.si.Create(&item)
	} else {
		if item.Quantity > 0 {
			err = s.si.Update(&item)
		} else {
			err = s.si.Delete(item.ID)
		}
	}
	return s.GetShoppingBasket(u)
}

func (s *shoppingBasketService) GetShoppingBasket(u uuid.UUID) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	total := 0.0
	if err != nil {
		return ShoppingBasket{}, err
	}
	sm := ShoppingBasket{
		ID: id.ID,
	}
	if len(id.Items) > 0 {
		ps := make([]ShoppingBasketItem, len(id.Items))
		for i, item := range id.Items {
			pr, _ := s.pm.GetProduct(item.ProductId)
			total += float64(pr.Price * item.Quantity)
			ps[i] = ShoppingBasketItem{
				item.ID, pr.Name, item.Price, pr.ID, pr.ImageUrl, item.Quantity,
			}
		}
		sm.Items = ps
	}
	sm.TotalPriceInclusive = float32(total)
	sm.Tax = float32(total - (total / 1.21))
	sm.TotalPriceExclusive = float32(total / 1.21)

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
