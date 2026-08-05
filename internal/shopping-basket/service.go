package shopping_basket

import (
	Log "log"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Port "github.com/danyel/ecommerce/internal/common/port"
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/product-management"
	Uuid "github.com/google/uuid"
	Gorm "gorm.io/gorm"
)

type ShoppingBasketService interface {
	CreateShoppingBasket() (ShoppingBasket, error)
	UpdateShoppingBasketItem(u Uuid.UUID, i UpdateShoppingBasketItem) (ShoppingBasket, error)
	GetShoppingBasket(u Uuid.UUID) (ShoppingBasket, error)
}

type shoppingBasketService struct {
	r         CommonRepository.CrudRepository[ShoppingBasketModel]
	p         Product.ProductService
	pm        ProductManagement.ProductService
	m         Product.ProductMapper
	si        CommonRepository.CrudRepository[ShoppingBasketItemModel]
	publisher Port.EventPublisher
}

func (s *shoppingBasketService) CreateShoppingBasket() (ShoppingBasket, error) {
	shoppingBasketModel := ShoppingBasketModel{}
	err := s.r.Create(&shoppingBasketModel)
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	r := ShoppingBasket{
		ID: Types.NewID(shoppingBasketModel.ID),
	}

	if err = s.publisher.Publish(ShoppingBasketCreated.Queue, ShoppingBasketCreatedEvent{
		Id: r.ID,
	}); err != nil {
		return r, err
	}

	return r, nil
}

func (s *shoppingBasketService) UpdateShoppingBasketItem(u Uuid.UUID, i UpdateShoppingBasketItem) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	var prd Product.Product
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	if prd, err = s.p.GetProduct(i.ProductId.ID); err != nil {
		return EmptyShoppingBasket(), err
	}

	item := ShoppingBasketItemModel{ID: Uuid.Nil, ShoppingBasketID: id.ID, ProductId: prd.ID.ID, Price: prd.Price, Quantity: i.Quantity}
	for _, it := range id.Items {
		if it.ProductId == item.ProductId {
			item.ID = it.ID
			item.Quantity = i.Quantity
		}
	}
	if item.ID == Uuid.Nil {
		err = s.si.Create(&item)
	} else {
		if item.Quantity > 0 {
			err = s.si.Update(&item)
		} else {
			err = s.si.Delete(item.ID)
		}
	}

	Log.Printf("ShoppingBasketItem To Publish: %v", item)
	if err = s.publisher.Publish(ShoppingBasketUpdated.Queue, ShoppingBasketUpdatedEvent{
		Id:        Types.NewID(u),
		Quantity:  i.Quantity,
		ProductId: i.ProductId,
	}); err != nil {
		return EmptyShoppingBasket(), err
	}
	return s.GetShoppingBasket(u)
}

func (s *shoppingBasketService) GetShoppingBasket(u Uuid.UUID) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	totalPriceInclusive := Types.Float64(0)
	totalPriceExclusive := Types.Float64(0)
	totalTaxes := Types.Float64(0)
	if err != nil {
		all := s.r.FetchAll()
		Log.Printf("Fetched: %v", all)
		return EmptyShoppingBasket(), err
	}
	sm := ShoppingBasket{
		ID: NewShoppingBasketId(id.ID),
	}
	if len(id.Items) > 0 {
		ps := make([]ShoppingBasketItem, len(id.Items))
		for i, item := range id.Items {
			pr, _ := s.pm.GetProduct(Types.NewID(item.ProductId))
			innerCalculation := pr.Price * float64(item.Quantity)
			innerTotalPriceInclusive := Types.Float64(innerCalculation)
			innerTotalPriceExclusive := Types.Float64(innerCalculation / 1.21)
			innerTotalTaxes := Types.Float64(innerCalculation - (innerCalculation / 1.21))
			totalPriceInclusive += innerTotalPriceInclusive
			totalPriceExclusive += innerTotalPriceExclusive
			totalTaxes += innerTotalTaxes
			ps[i] = ShoppingBasketItem{
				ID:   Types.NewID(item.ID),
				Name: pr.Name,
				BasePrice: Types.Price{
					Inclusive: Types.Float64(item.Price),
					Exclusive: Types.Float64(item.Price / 1.21),
					Tax:       Types.Float64(item.Price - (item.Price / 1.21)),
					Currency:  "EUR",
				},
				TotalPrice: Types.Price{
					Inclusive: innerTotalPriceInclusive,
					Exclusive: innerTotalPriceExclusive,
					Tax:       innerTotalTaxes,
					Currency:  "EUR",
				},
				ProductId: pr.ID,
				ImageUrl:  pr.ImageUrl,
				Quantity:  item.Quantity,
			}
		}
		sm.Items = ps
	}
	sm.TotalPrice = Types.Price{
		Inclusive: totalPriceInclusive,
		Exclusive: totalPriceExclusive,
		Tax:       totalTaxes,
		Currency:  "EUR",
	}

	return sm, nil
}

func NewService(db *Gorm.DB, publisher Port.EventPublisher) ShoppingBasketService {
	r := CommonRepository.NewCrudRepository[ShoppingBasketModel](db)
	p := Product.NewProductService(db)
	s := ProductManagement.NewProductService(db)
	m := Product.NewProductMapper(Category.NewCategoryService(db), CMS.NewCmsService(db))
	si := CommonRepository.NewCrudRepository[ShoppingBasketItemModel](db)
	return &shoppingBasketService{r, p, s, m, si, publisher}
}
