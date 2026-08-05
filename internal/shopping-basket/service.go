package shopping_basket

import (
	"log"

	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	"github.com/danyel/ecommerce/internal/common/port"
	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
	"github.com/danyel/ecommerce/internal/common/types"
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
	r         commonRepository.CrudRepository[ShoppingBasketModel]
	p         product.ProductService
	pm        productmanagement.ProductService
	m         product.ProductMapper
	si        commonRepository.CrudRepository[ShoppingBasketItemModel]
	publisher port.EventPublisher
}

func (s *shoppingBasketService) CreateShoppingBasket() (ShoppingBasket, error) {
	shoppingBasketModel := ShoppingBasketModel{}
	err := s.r.Create(&shoppingBasketModel)
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	r := ShoppingBasket{
		ID: types.NewID(shoppingBasketModel.ID),
	}

	if err = s.publisher.Publish(ShoppingBasketCreatedQueue, ShoppingBasketCreatedEvent{
		Id: r.ID,
	}); err != nil {
		return r, err
	}

	return r, nil
}

func (s *shoppingBasketService) UpdateShoppingBasketItem(u uuid.UUID, i UpdateShoppingBasketItem) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	var prd product.Product
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	if prd, err = s.p.GetProduct(i.ProductId.ID); err != nil {
		return EmptyShoppingBasket(), err
	}

	item := ShoppingBasketItemModel{ID: uuid.Nil, ShoppingBasketID: id.ID, ProductId: prd.ID.ID, Price: prd.Price, Quantity: i.Quantity}
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

	if err = s.publisher.Publish(ShoppingBasketUpdatedQueue, ShoppingBasketUpdatedEvent{
		Id:        types.NewID(u),
		Quantity:  i.Quantity,
		ProductId: i.ProductId,
	}); err != nil {
		return EmptyShoppingBasket(), err
	}
	return s.GetShoppingBasket(u)
}

func (s *shoppingBasketService) GetShoppingBasket(u uuid.UUID) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	total := types.Float64(0)
	if err != nil {
		all := s.r.FetchAll()
		log.Printf("Fetched: %v", all)
		return EmptyShoppingBasket(), err
	}
	sm := ShoppingBasket{
		ID: NewShoppingBasketId(id.ID),
	}
	if len(id.Items) > 0 {
		ps := make([]ShoppingBasketItem, len(id.Items))
		for i, item := range id.Items {
			pr, _ := s.pm.GetProduct(types.NewID(item.ProductId))
			calculatedPrice := pr.Price * float64(item.Quantity)
			total += types.Float64(calculatedPrice)
			ps[i] = ShoppingBasketItem{
				types.NewID(item.ID), pr.Name, types.Float64(item.Price), types.Float64(item.Price / 1.21), types.Float64(item.Price - (item.Price / 1.21)), types.Float64(calculatedPrice), types.Float64(calculatedPrice / 1.21), types.Float64(calculatedPrice - (calculatedPrice / 1.21)), pr.ID, pr.ImageUrl, item.Quantity,
			}
		}
		sm.Items = ps
	}
	sm.TotalPriceInclusive = total
	sm.Tax = total - (total / 1.21)
	sm.TotalPriceExclusive = total / 1.21

	return sm, nil
}

func NewService(db *gorm.DB, publisher port.EventPublisher) ShoppingBasketService {
	r := commonRepository.NewCrudRepository[ShoppingBasketModel](db)
	p := product.NewProductService(db)
	s := productmanagement.NewProductService(db)
	m := product.NewProductMapper(category.NewCategoryService(db), cms.NewCmsService(db))
	si := commonRepository.NewCrudRepository[ShoppingBasketItemModel](db)
	return &shoppingBasketService{r, p, s, m, si, publisher}
}
