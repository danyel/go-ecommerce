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

type ApplicationContextFactory interface {
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

type applicationContextFactory struct {
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

func (applicationContextFactory *applicationContextFactory) ReservationService() Reservation.ReservationService {
	return getInstanceOfType(&applicationContextFactory.reservationService, func() Reservation.ReservationService {
		return Reservation.NewService(repositoryContextFactoryInstance.ReservationRepository())
	})
}

func (applicationContextFactory *applicationContextFactory) CategoryService() Category.CategoryService {
	return getInstanceOfType(&applicationContextFactory.categoryService, func() Category.CategoryService {
		return Category.NewService(repositoryContextFactoryInstance.CategoryRepository())
	})
}

func (applicationContextFactory *applicationContextFactory) ProductService() Product.ProductService {
	return getInstanceOfType(&applicationContextFactory.productService, func() Product.ProductService {
		return Product.NewService(repositoryContextFactoryInstance.ProductRepository(), applicationContextFactory.ProductMapper())
	})
}

func (applicationContextFactory *applicationContextFactory) ShoppingBasketService() ShoppingBasket.ShoppingBasketService {
	return getInstanceOfType(&applicationContextFactory.shoppingBasketService, func() ShoppingBasket.ShoppingBasketService {
		return ShoppingBasket.NewService(applicationContextFactory.ProductService(), applicationContextFactory.ProductManagementService(), applicationContextFactory.ProductMapper(), repositoryContextFactoryInstance.ShoppingBasketRepository(), repositoryContextFactoryInstance.ShoppingBasketItemRepository(), messageBrokerContextFactoryInstance.MessageBroker())
	})
}

func (applicationContextFactory *applicationContextFactory) CmsService() CMS.CmsService {
	return getInstanceOfType(&applicationContextFactory.cmsService, func() CMS.CmsService {
		return CMS.NewService(repositoryContextFactoryInstance.CmsRepository())
	})
}

func (applicationContextFactory *applicationContextFactory) ManagementService() Management.ManagementService {
	return getInstanceOfType(&applicationContextFactory.managementService, func() Management.ManagementService {
		return Management.NewService(repositoryContextFactoryInstance.CmsRepository())
	})
}

func (applicationContextFactory *applicationContextFactory) ProductManagementService() ProductManagement.ProductManagementService {
	return getInstanceOfType(&applicationContextFactory.productManagementService, func() ProductManagement.ProductManagementService {
		return ProductManagement.NewService(repositoryContextFactoryInstance.ProductRepository(), applicationContextFactory.ProductService())
	})
}

func (applicationContextFactory *applicationContextFactory) ProductMapper() Product.ProductMapper {
	return getInstanceOfType(&applicationContextFactory.productMapper, func() Product.ProductMapper {
		return Product.NewProductMapper(applicationContextFactory.CategoryService(), applicationContextFactory.CmsService())
	})
}

func (applicationContextFactory *applicationContextFactory) MessageBroker() *MessageBroker.MessageBroker {
	return applicationContextFactory.messageBroker
}

func (applicationContextFactory *applicationContextFactory) StartMessageBroker() error {
	if applicationContextFactory.messageBroker != nil {
		return applicationContextFactory.messageBroker.Start()
	}
	return errors.New("message broker is nil")
}

func BuildApplicationContextFactory() ApplicationContextFactory {
	applicationContextFactoryOnce.Do(func() {
		applicationContextFactoryInstance = &applicationContextFactory{
			messageBroker: messageBrokerContextFactoryInstance.MessageBroker(),
		}
	})
	return applicationContextFactoryInstance
}
