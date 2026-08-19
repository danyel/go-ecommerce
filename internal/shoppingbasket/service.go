package shoppingbasket

import (
	Log "log"

	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Port "github.com/danyel/ecommerce/internal/common/port"
	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
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
		ID: r.ID,
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
	if prd, err = s.p.GetProduct(i.ProductID.ID); err != nil {
		return EmptyShoppingBasket(), err
	}

	item := ShoppingBasketItemModel{ID: Uuid.Nil, ShoppingBasketID: id.ID, ProductID: prd.ID.ID, Price: float64(prd.Price.Inclusive), Quantity: i.Quantity}
	for _, it := range id.Items {
		if it.ProductID == item.ProductID {
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
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	Log.Printf("ShoppingBasketItem To Publish: %v", item)
	if err = s.publisher.Publish(ShoppingBasketUpdated.Queue, ShoppingBasketUpdatedEvent{
		ID:        Types.NewID(u),
		Quantity:  i.Quantity,
		ProductID: i.ProductID,
	}); err != nil {
		return EmptyShoppingBasket(), err
	}
	return s.GetShoppingBasket(u)
}

func (s *shoppingBasketService) GetShoppingBasket(u Uuid.UUID) (ShoppingBasket, error) {
	id, err := s.r.FindById(u, "Items")
	totalPrice := float64(0)
	fetchAll := s.r.FetchAll()
	Log.Printf("ShoppingBasket To Publish: %v", fetchAll)
	if err != nil {
		all := s.r.FetchAll()
		Log.Printf("Fetched: %v", all)
		return EmptyShoppingBasket(), err
	}
	sm := ShoppingBasket{
		ID: Types.NewID(id.ID),
	}
	if len(id.Items) > 0 {
		ps := make([]ShoppingBasketItem, len(id.Items))
		for i, item := range id.Items {
			pr, _ := s.pm.GetProduct(Types.NewID(item.ProductID))
			price := float64(pr.Price.Inclusive) * float64(item.Quantity)
			totalPrice += price
			ps[i] = ShoppingBasketItem{
				ID:         Types.NewID(item.ID),
				Name:       pr.Name,
				BasePrice:  Types.NewPrice(item.Price, "EUR"),
				TotalPrice: Types.NewPrice(price, "EUR"),
				ProductID:  pr.ID,
				ImageURL:   pr.ImageURL,
				Quantity:   item.Quantity,
			}
		}
		sm.Items = ps
	}
	sm.TotalPrice = Types.NewPrice(totalPrice, "EUR")

	return sm, nil
}

func NewService(db *Database.DB, publisher Port.EventPublisher) ShoppingBasketService {
	r := CommonRepository.NewCrudRepository[ShoppingBasketModel](db)
	p := Product.NewProductService(db)
	s := ProductManagement.NewProductService(db)
	m := Product.NewProductMapper(Category.NewCategoryService(db), CMS.NewCmsService(db))
	si := CommonRepository.NewCrudRepository[ShoppingBasketItemModel](db)
	return &shoppingBasketService{r, p, s, m, si, publisher}
}
