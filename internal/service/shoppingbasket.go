package service

import (
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Domain "github.com/danyel/ecommerce/internal/domain"
	Mapper "github.com/danyel/ecommerce/internal/mapper"
	Persistence "github.com/danyel/ecommerce/internal/persistence"
	Port "github.com/danyel/ecommerce/internal/port"
	Types "github.com/danyel/ecommerce/internal/types"
	Validator "github.com/danyel/ecommerce/internal/validator"
	Uuid "github.com/google/uuid"
)

type IShoppingBasketService interface {
	Create() (Domain.ShoppingBasket, error)
	Update(ID Uuid.UUID, update Domain.ShoppingBasketItemUpdate) error
	FindByID(u Uuid.UUID) (Domain.ShoppingBasket, error)
}

type shoppingBasketService struct {
	shoppingBasketRepository     Persistence.CrudRepository[Persistence.ShoppingBasketModel]
	productService               IProductService
	productManagementService     IProductManagementService
	productMapper                Mapper.ProductMapper
	shoppingBasketItemRepository Persistence.CrudRepository[Persistence.ShoppingBasketItemModel]
	publisher                    Port.EventPublisher
	shoppingBasketValidator      Validator.Validator
}

func (shoppingBasketService *shoppingBasketService) Create() (Domain.ShoppingBasket, error) {
	shoppingBasketModel := Persistence.ShoppingBasketModel{}
	err := shoppingBasketService.shoppingBasketRepository.Create(&shoppingBasketModel)
	if err != nil {
		return Domain.EmptyShoppingBasket(), err
	}
	r := Domain.ShoppingBasket{
		ID: Types.NewID(shoppingBasketModel.ID),
	}
	return r, nil
}

func (shoppingBasketService *shoppingBasketService) Update(ID Uuid.UUID, update Domain.ShoppingBasketItemUpdate) error {
	shoppingBasketModel, err := shoppingBasketService.shoppingBasketRepository.FindByID(ID, "Items")
	if err != nil {
		return err
	}
	product, err := shoppingBasketService.productService.FindByID(update.ProductID)
	if err != nil {
		return err
	}

	currentItem := Persistence.ShoppingBasketItemModel{}
	for _, item := range shoppingBasketModel.Items {
		if item.ProductID == product.ID.ID {
			currentItem = item
			break
		}
	}
	shoppingBasketItemModel := Persistence.ShoppingBasketItemModel{
		ID:               currentItem.ID,
		ShoppingBasketID: shoppingBasketModel.ID,
		ProductID:        product.ID.ID,
		Price:            float64(product.Price.Inclusive),
		Quantity:         update.Quantity,
	}
	product.Stock += Domain.StockAdjustment(currentItem.Quantity, update.Quantity)

	if err := shoppingBasketService.productService.Update(product); err != nil {
		return err
	}

	if shoppingBasketItemModel.ID == Uuid.Nil {
		err = shoppingBasketService.shoppingBasketItemRepository.Create(&shoppingBasketItemModel)
	} else if shoppingBasketItemModel.ID != Uuid.Nil {
		if shoppingBasketItemModel.Quantity > 0 {
			err = shoppingBasketService.shoppingBasketItemRepository.Update(&shoppingBasketItemModel)
		} else {
			err = shoppingBasketService.shoppingBasketItemRepository.Delete(shoppingBasketItemModel.ID)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (shoppingBasketService *shoppingBasketService) FindByID(ID Uuid.UUID) (Domain.ShoppingBasket, error) {
	shoppingBasketModel, err := shoppingBasketService.shoppingBasketRepository.FindByID(ID, "Items")
	Logger.Log.Debug("Shopping Basket By Id: %v", shoppingBasketModel)
	if err != nil {
		return Domain.EmptyShoppingBasket(), err
	}
	shoppingBasket := Domain.ShoppingBasket{
		ID: Types.NewID(shoppingBasketModel.ID),
	}
	if len(shoppingBasketModel.Items) > 0 {
		shoppingBasketItems := make([]Domain.ShoppingBasketItem, len(shoppingBasketModel.Items))
		for index, shoppingBasketItemModel := range shoppingBasketModel.Items {
			currentProduct, err := shoppingBasketService.productService.FindByID(shoppingBasketItemModel.ProductID)
			if err != nil {
				return Domain.EmptyShoppingBasket(), err
			}
			shoppingBasketItems[index] = Domain.ShoppingBasketItem{
				Product:  currentProduct,
				Quantity: shoppingBasketItemModel.Quantity,
			}
		}
		shoppingBasket.Items = shoppingBasketItems
	}
	return shoppingBasket, nil
}

func calculateTotal(price Types.Float64, quantity int) float64 {
	return float64(price) * float64(quantity)
}

func ShoppingBasketService(productService IProductService, productManagementService IProductManagementService, productMapper Mapper.ProductMapper, shoppingBasketRepository Persistence.CrudRepository[Persistence.ShoppingBasketModel], shoppingBasketItemRepository Persistence.CrudRepository[Persistence.ShoppingBasketItemModel], publisher Port.EventPublisher, shoppingBasketValidator Validator.Validator) IShoppingBasketService {
	return &shoppingBasketService{
		shoppingBasketRepository,
		productService,
		productManagementService,
		productMapper,
		shoppingBasketItemRepository,
		publisher,
		shoppingBasketValidator,
	}
}
