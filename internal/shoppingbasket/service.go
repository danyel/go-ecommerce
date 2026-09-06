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
	Create() (ShoppingBasket, error)
	Update(ID Uuid.UUID, UpdateShoppingBasketItem UpdateShoppingBasketItem) (ShoppingBasket, error)
	FindById(u Uuid.UUID) (ShoppingBasket, error)
}

type shoppingBasketService struct {
	shoppingBasketRepository     Repository.CrudRepository[ShoppingBasketModel]
	productService               Product.ProductService
	productManagementService     ProductManagement.ProductManagementService
	productMapper                Product.ProductMapper
	shoppingBasketItemRepository Repository.CrudRepository[ShoppingBasketItemModel]
	publisher                    Port.EventPublisher
}

func (shoppingBasketService *shoppingBasketService) Create() (ShoppingBasket, error) {
	shoppingBasketModel := ShoppingBasketModel{}
	err := shoppingBasketService.shoppingBasketRepository.Create(&shoppingBasketModel)
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	r := ShoppingBasket{
		ID: Types.NewID(shoppingBasketModel.ID),
	}
	return r, nil
}

func (shoppingBasketService *shoppingBasketService) Update(ID Uuid.UUID, updateShoppingBasketItem UpdateShoppingBasketItem) (ShoppingBasket, error) {
	shoppingBasketModel, err := shoppingBasketService.shoppingBasketRepository.FindById(ID, "Items")
	var product Product.Product
	if err != nil {
		return EmptyShoppingBasket(), err
	}
	if product, err = shoppingBasketService.productService.FindById(updateShoppingBasketItem.ProductID.ID); err != nil {
		return EmptyShoppingBasket(), err
	}

	shoppingBasketItemModel := ShoppingBasketItemModel{ID: Uuid.Nil, ShoppingBasketID: shoppingBasketModel.ID, ProductID: product.ID.ID, Price: float64(product.Price.Inclusive), Quantity: updateShoppingBasketItem.Quantity}
	for _, currentItem := range shoppingBasketModel.Items {
		if currentItem.ProductID == shoppingBasketItemModel.ProductID {
			if updateShoppingBasketItem.Quantity != currentItem.Quantity {
				Logger.Log.Debug("Going to publish %s for product id %s with quantity: %d", Product.UpdateProductStock.Queue, currentItem.ProductID, -(currentItem.Quantity - updateShoppingBasketItem.Quantity))
				e := shoppingBasketService.publisher.Publish(Product.UpdateProductStock.Queue, Product.UpdateProductStockCommand{
					ProductID:        currentItem.ProductID,
					Quantity:         -(currentItem.Quantity - updateShoppingBasketItem.Quantity),
					ShoppingBasketId: shoppingBasketModel.ID,
				})
				if e != nil {
					return EmptyShoppingBasket(), e
				}
			}
			shoppingBasketItemModel.ID = currentItem.ID
			shoppingBasketItemModel.Quantity = updateShoppingBasketItem.Quantity
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
	return shoppingBasketService.FindById(ID)
}

func (shoppingBasketService *shoppingBasketService) FindById(ID Uuid.UUID) (ShoppingBasket, error) {
	shoppingBasketModel, err := shoppingBasketService.shoppingBasketRepository.FindById(ID, "Items")
	Logger.Log.Debug("Shopping Basket By Id: %v", shoppingBasketModel)
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
