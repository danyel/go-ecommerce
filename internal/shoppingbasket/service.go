package shoppingbasket

import (
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Port "github.com/danyel/ecommerce/internal/common/port"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Types "github.com/danyel/ecommerce/internal/common/types"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	Uuid "github.com/google/uuid"
)

//goland:noinspection GoNameStartsWithPackageName
type ShoppingBasketService interface {
	CreateShoppingBasket() (ShoppingBasket, error)
	UpdateShoppingBasketItem(ID Uuid.UUID, UpdateShoppingBasketItem UpdateShoppingBasketItem) (ShoppingBasket, error)
	GetShoppingBasket(u Uuid.UUID) (ShoppingBasket, error)
}

type shoppingBasketService struct {
	shoppingBasketRepository     Repository.CrudRepository[ShoppingBasketModel]
	productService               Product.ProductService
	productManagementService     ProductManagement.ProductManagementService
	productMapper                Product.ProductMapper
	shoppingBasketItemRepository Repository.CrudRepository[ShoppingBasketItemModel]
	publisher                    Port.EventPublisher
}

func (shoppingBasketService *shoppingBasketService) CreateShoppingBasket() (ShoppingBasket, error) {
	shoppingBasketModel := ShoppingBasketModel{}
	err := shoppingBasketService.shoppingBasketRepository.Create(&shoppingBasketModel)
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	r := ShoppingBasket{
		ID: Types.NewID(shoppingBasketModel.ID),
	}

	if err = shoppingBasketService.publisher.Publish(ShoppingBasketCreated.Queue, ShoppingBasketCreatedEvent{
		ID: r.ID,
	}); err != nil {
		return r, err
	}

	return r, nil
}

func (shoppingBasketService *shoppingBasketService) UpdateShoppingBasketItem(ID Uuid.UUID, i UpdateShoppingBasketItem) (ShoppingBasket, error) {
	shoppingBasketModel, err := shoppingBasketService.shoppingBasketRepository.FindById(ID, "Items")
	var product Product.Product
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	if product, err = shoppingBasketService.productService.GetProduct(i.ProductID.ID); err != nil {
		return EmptyShoppingBasket(), err
	}

	shoppingBasketItemModel := ShoppingBasketItemModel{ID: Uuid.Nil, ShoppingBasketID: shoppingBasketModel.ID, ProductID: product.ID.ID, Price: float64(product.Price.Inclusive), Quantity: i.Quantity}
	for _, it := range shoppingBasketModel.Items {
		if it.ProductID == shoppingBasketItemModel.ProductID {
			shoppingBasketItemModel.ID = it.ID
			shoppingBasketItemModel.Quantity = i.Quantity
		}
	}
	if shoppingBasketItemModel.ID == Uuid.Nil {
		err = shoppingBasketService.shoppingBasketItemRepository.Create(&shoppingBasketItemModel)
	} else {
		if shoppingBasketItemModel.Quantity > 0 {
			err = shoppingBasketService.shoppingBasketItemRepository.Update(&shoppingBasketItemModel)
		} else {
			err = shoppingBasketService.shoppingBasketItemRepository.Delete(shoppingBasketItemModel.ID)
		}
	}
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	Logger.Log.Debug("ShoppingBasketItem To Publish: %v", shoppingBasketItemModel)
	if err = shoppingBasketService.publisher.Publish(ShoppingBasketUpdated.Queue, ShoppingBasketUpdatedEvent{
		ID:        Types.NewID(ID),
		Quantity:  i.Quantity,
		ProductID: i.ProductID,
	}); err != nil {
		return EmptyShoppingBasket(), err
	}
	return shoppingBasketService.GetShoppingBasket(ID)
}

func (shoppingBasketService *shoppingBasketService) GetShoppingBasket(ID Uuid.UUID) (ShoppingBasket, error) {
	shoppingBasketModel, err := shoppingBasketService.shoppingBasketRepository.FindById(ID, "Items")
	totalPrice := float64(0)
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	shoppingBasket := ShoppingBasket{
		ID: Types.NewID(shoppingBasketModel.ID),
	}
	if len(shoppingBasketModel.Items) > 0 {
		shoppingBasketItems := make([]ShoppingBasketItem, len(shoppingBasketModel.Items))
		for index, shoppingBasketItemModel := range shoppingBasketModel.Items {
			currentProduct, _ := shoppingBasketService.productManagementService.GetProduct(Types.NewID(shoppingBasketItemModel.ProductID))
			calculatedPriceToAdd := calculateTotal(currentProduct.Price.Inclusive, shoppingBasketItemModel.Quantity)
			totalPrice += calculatedPriceToAdd
			shoppingBasketItems[index] = ShoppingBasketItem{
				ID:         Types.NewID(shoppingBasketItemModel.ID),
				Name:       currentProduct.Name,
				BasePrice:  Types.NewPrice(shoppingBasketItemModel.Price, "EUR"),
				TotalPrice: Types.NewPrice(calculatedPriceToAdd, "EUR"),
				ProductID:  currentProduct.ID,
				ImageURL:   currentProduct.ImageURL,
				Quantity:   shoppingBasketItemModel.Quantity,
			}
		}
		shoppingBasket.Items = shoppingBasketItems
	}
	shoppingBasket.TotalPrice = Types.NewPrice(totalPrice, "EUR")

	return shoppingBasket, nil
}

func calculateTotal(price Types.Float64, quantity int) float64 {
	return float64(price) * float64(quantity)
}

func NewService(productService Product.ProductService, productManagementService ProductManagement.ProductManagementService, productMapper Product.ProductMapper, shoppingBasketRepository Repository.CrudRepository[ShoppingBasketModel], shoppingBasketItemRepository Repository.CrudRepository[ShoppingBasketItemModel], publisher Port.EventPublisher) ShoppingBasketService {
	return &shoppingBasketService{
		shoppingBasketRepository,
		productService,
		productManagementService,
		productMapper,
		shoppingBasketItemRepository,
		publisher,
	}
}
