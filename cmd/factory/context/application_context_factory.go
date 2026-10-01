package contextfactory

import (
	Errors "errors"

	MessageBroker "github.com/danyel/ecommerce/cmd/broker"
	Mapper "github.com/danyel/ecommerce/internal/mapper"
	Category "github.com/danyel/ecommerce/internal/service"
	ShoppingBasket "github.com/danyel/ecommerce/internal/validator"
)

type ServiceContextFactory interface {
	ReservationService() Category.IReservationService
	ProductService() Category.IProductService
	CategoryService() Category.ICategoryService
	ShoppingBasketService() Category.IShoppingBasketService
	CmsService() Category.ICmsService
	ManagementService() Category.IManagementService
	ProductManagementService() Category.IProductManagementService
	ProductMapper() Mapper.ProductMapper
	CategoryMapper() Mapper.ICategoryMapper
	MessageBroker() *MessageBroker.MessageBroker
	ShoppingBasketValidator() ShoppingBasket.Validator
	StartMessageBroker() error
}

type serviceContextFactory struct {
	reservationService       Category.IReservationService
	productService           Category.IProductService
	categoryService          Category.ICategoryService
	shoppingBasketService    Category.IShoppingBasketService
	managementService        Category.IManagementService
	productManagementService Category.IProductManagementService
	productMapper            Mapper.ProductMapper
	cmsService               Category.ICmsService
	messageBroker            *MessageBroker.MessageBroker
	shoppingBasketValidator  ShoppingBasket.Validator
	categoryMapper           Mapper.ICategoryMapper
}

func (serviceContextFactory *serviceContextFactory) ReservationService() Category.IReservationService {
	return getInstanceOfType(&serviceContextFactory.reservationService, func() Category.IReservationService {
		return Category.ReservationService(repositoryContextFactoryInstance.ReservationRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) CategoryService() Category.ICategoryService {
	return getInstanceOfType(&serviceContextFactory.categoryService, func() Category.ICategoryService {
		return Category.CategoryService(repositoryContextFactoryInstance.CategoryRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) ProductService() Category.IProductService {
	return getInstanceOfType(&serviceContextFactory.productService, func() Category.IProductService {
		return Category.ProductService(repositoryContextFactoryInstance.ProductRepository(), serviceContextFactory.ProductMapper(), applicationContextFactoryInstance.MessageBroker())
	})
}

func (serviceContextFactory *serviceContextFactory) ShoppingBasketService() Category.IShoppingBasketService {
	return getInstanceOfType(&serviceContextFactory.shoppingBasketService, func() Category.IShoppingBasketService {
		return Category.ShoppingBasketService(serviceContextFactory.ProductService(), serviceContextFactory.ProductManagementService(), serviceContextFactory.ProductMapper(), repositoryContextFactoryInstance.ShoppingBasketRepository(), repositoryContextFactoryInstance.ShoppingBasketItemRepository(), messageBrokerContextFactoryInstance.MessageBroker(), serviceContextFactory.ShoppingBasketValidator())
	})
}

func (serviceContextFactory *serviceContextFactory) ShoppingBasketValidator() ShoppingBasket.Validator {
	return getInstanceOfType(&serviceContextFactory.shoppingBasketValidator, func() ShoppingBasket.Validator {
		return ShoppingBasket.NewValidator(repositoryContextFactoryInstance.ProductRepository(), repositoryContextFactoryInstance.ShoppingBasketRepository(), serviceContextFactory.ProductService())
	})
}

func (serviceContextFactory *serviceContextFactory) CmsService() Category.ICmsService {
	return getInstanceOfType(&serviceContextFactory.cmsService, func() Category.ICmsService {
		return Category.CmsService(repositoryContextFactoryInstance.CmsRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) ManagementService() Category.IManagementService {
	return getInstanceOfType(&serviceContextFactory.managementService, func() Category.IManagementService {
		return Category.ManagementService(repositoryContextFactoryInstance.CmsRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) ProductManagementService() Category.IProductManagementService {
	return getInstanceOfType(&serviceContextFactory.productManagementService, func() Category.IProductManagementService {
		return Category.ProductManagementService(repositoryContextFactoryInstance.ProductRepository(), serviceContextFactory.ProductService())
	})
}

func (serviceContextFactory *serviceContextFactory) ProductMapper() Mapper.ProductMapper {
	return getInstanceOfType(&serviceContextFactory.productMapper, func() Mapper.ProductMapper {
		return Mapper.NewProductMapper(serviceContextFactory.CategoryService(), serviceContextFactory.CmsService())
	})
}

func (serviceContextFactory *serviceContextFactory) CategoryMapper() Mapper.ICategoryMapper {
	return getInstanceOfType(&serviceContextFactory.categoryMapper, func() Mapper.ICategoryMapper {
		return Mapper.CategoryMapper()
	})
}

func (serviceContextFactory *serviceContextFactory) MessageBroker() *MessageBroker.MessageBroker {
	return serviceContextFactory.messageBroker
}

func (serviceContextFactory *serviceContextFactory) StartMessageBroker() error {
	if serviceContextFactory.messageBroker != nil {
		return serviceContextFactory.messageBroker.Start()
	}
	return Errors.New("message broker is nil")
}

func BuildServiceContextFactory() ServiceContextFactory {
	applicationContextFactoryOnce.Do(func() {
		applicationContextFactoryInstance = &serviceContextFactory{
			messageBroker: messageBrokerContextFactoryInstance.MessageBroker(),
		}
	})
	return applicationContextFactoryInstance
}
