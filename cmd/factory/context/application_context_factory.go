package contextfactory

import (
	"errors"

	MessageBroker "github.com/danyel/ecommerce/cmd/broker"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	Reservation "github.com/danyel/ecommerce/internal/reservation"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
)

type ServiceContextFactory interface {
	ReservationService() Reservation.ReservationService
	ProductService() Product.ProductService
	CategoryService() Category.CategoryService
	ShoppingBasketService() ShoppingBasket.ShoppingBasketService
	CmsService() CMS.CmsService
	ManagementService() Management.ManagementService
	ProductManagementService() ProductManagement.ProductManagementService
	ProductMapper() Product.ProductMapper
	MessageBroker() *MessageBroker.MessageBroker
	StartMessageBroker() error
}

type serviceContextFactory struct {
	reservationService       Reservation.ReservationService
	productService           Product.ProductService
	categoryService          Category.CategoryService
	shoppingBasketService    ShoppingBasket.ShoppingBasketService
	managementService        Management.ManagementService
	productManagementService ProductManagement.ProductManagementService
	productMapper            Product.ProductMapper
	cmsService               CMS.CmsService
	messageBroker            *MessageBroker.MessageBroker
}

func (serviceContextFactory *serviceContextFactory) ReservationService() Reservation.ReservationService {
	return getInstanceOfType(&serviceContextFactory.reservationService, func() Reservation.ReservationService {
		return Reservation.NewService(repositoryContextFactoryInstance.ReservationRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) CategoryService() Category.CategoryService {
	return getInstanceOfType(&serviceContextFactory.categoryService, func() Category.CategoryService {
		return Category.NewService(repositoryContextFactoryInstance.CategoryRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) ProductService() Product.ProductService {
	return getInstanceOfType(&serviceContextFactory.productService, func() Product.ProductService {
		return Product.NewService(repositoryContextFactoryInstance.ProductRepository(), serviceContextFactory.ProductMapper(), applicationContextFactoryInstance.MessageBroker())
	})
}

func (serviceContextFactory *serviceContextFactory) ShoppingBasketService() ShoppingBasket.ShoppingBasketService {
	return getInstanceOfType(&serviceContextFactory.shoppingBasketService, func() ShoppingBasket.ShoppingBasketService {
		return ShoppingBasket.NewService(serviceContextFactory.ProductService(), serviceContextFactory.ProductManagementService(), serviceContextFactory.ProductMapper(), repositoryContextFactoryInstance.ShoppingBasketRepository(), repositoryContextFactoryInstance.ShoppingBasketItemRepository(), messageBrokerContextFactoryInstance.MessageBroker())
	})
}

func (serviceContextFactory *serviceContextFactory) CmsService() CMS.CmsService {
	return getInstanceOfType(&serviceContextFactory.cmsService, func() CMS.CmsService {
		return CMS.NewService(repositoryContextFactoryInstance.CmsRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) ManagementService() Management.ManagementService {
	return getInstanceOfType(&serviceContextFactory.managementService, func() Management.ManagementService {
		return Management.NewService(repositoryContextFactoryInstance.CmsRepository())
	})
}

func (serviceContextFactory *serviceContextFactory) ProductManagementService() ProductManagement.ProductManagementService {
	return getInstanceOfType(&serviceContextFactory.productManagementService, func() ProductManagement.ProductManagementService {
		return ProductManagement.NewService(repositoryContextFactoryInstance.ProductRepository(), serviceContextFactory.ProductService())
	})
}

func (serviceContextFactory *serviceContextFactory) ProductMapper() Product.ProductMapper {
	return getInstanceOfType(&serviceContextFactory.productMapper, func() Product.ProductMapper {
		return Product.NewProductMapper(serviceContextFactory.CategoryService(), serviceContextFactory.CmsService())
	})
}

func (serviceContextFactory *serviceContextFactory) MessageBroker() *MessageBroker.MessageBroker {
	return serviceContextFactory.messageBroker
}

func (serviceContextFactory *serviceContextFactory) StartMessageBroker() error {
	if serviceContextFactory.messageBroker != nil {
		return serviceContextFactory.messageBroker.Start()
	}
	return errors.New("message broker is nil")
}

func BuildServiceContextFactory() ServiceContextFactory {
	applicationContextFactoryOnce.Do(func() {
		applicationContextFactoryInstance = &serviceContextFactory{
			messageBroker: messageBrokerContextFactoryInstance.MessageBroker(),
		}
	})
	return applicationContextFactoryInstance
}
