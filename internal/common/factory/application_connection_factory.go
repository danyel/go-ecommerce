package factory

import (
	Sync "sync"

	MessageBroker "github.com/danyel/ecommerce/cmd/broker"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	Reservation "github.com/danyel/ecommerce/internal/reservation"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
)

type ApplicationConnectionFactory interface {
	ReservationService() Reservation.ReservationService
	ProductService() Product.ProductService
	CategoryService() Category.CategoryService
	ShoppingBasketService() ShoppingBasket.ShoppingBasketService
	CmsService() CMS.CmsService
	ManagementService() Management.ManagementService
	ProductManagementService() ProductManagement.ProductManagementService
	ProductMapper() Product.ProductMapper
}

type applicationConnectionFactory struct {
	databaseConnectionFactory DatabaseConnectionFactory
	mutex                     Sync.Mutex
	messageBroker             *MessageBroker.MessageBroker

	reservationService       Reservation.ReservationService
	productService           Product.ProductService
	categoryService          Category.CategoryService
	shoppingBasketService    ShoppingBasket.ShoppingBasketService
	managementService        Management.ManagementService
	productManagementService ProductManagement.ProductManagementService
	productMapper            Product.ProductMapper
	cmsService               CMS.CmsService
}

func (applicationConnectionFactory *applicationConnectionFactory) ReservationService() Reservation.ReservationService {
	return getInstanceOfType(&applicationConnectionFactory.reservationService, func() Reservation.ReservationService {
		return Reservation.NewService(applicationConnectionFactory.databaseConnectionFactory.ReservationRepository())
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) CategoryService() Category.CategoryService {
	return getInstanceOfType(&applicationConnectionFactory.categoryService, func() Category.CategoryService {
		return Category.NewService(applicationConnectionFactory.databaseConnectionFactory.CategoryRepository())
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) ProductService() Product.ProductService {
	return getInstanceOfType(&applicationConnectionFactory.productService, func() Product.ProductService {
		return Product.NewService(applicationConnectionFactory.databaseConnectionFactory.ProductRepository(), applicationConnectionFactory.ProductMapper())
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) ShoppingBasketService() ShoppingBasket.ShoppingBasketService {
	return getInstanceOfType(&applicationConnectionFactory.shoppingBasketService, func() ShoppingBasket.ShoppingBasketService {
		return ShoppingBasket.NewService(applicationConnectionFactory.ProductService(), applicationConnectionFactory.ProductManagementService(), applicationConnectionFactory.ProductMapper(), applicationConnectionFactory.databaseConnectionFactory.ShoppingBasketRepository(), applicationConnectionFactory.databaseConnectionFactory.ShoppingBasketItemRepository(), applicationConnectionFactory.messageBroker)
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) CmsService() CMS.CmsService {
	return getInstanceOfType(&applicationConnectionFactory.cmsService, func() CMS.CmsService {
		return CMS.NewService(applicationConnectionFactory.databaseConnectionFactory.CmsRepository())
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) ManagementService() Management.ManagementService {
	return getInstanceOfType(&applicationConnectionFactory.managementService, func() Management.ManagementService {
		return Management.NewService(applicationConnectionFactory.databaseConnectionFactory.CmsRepository())
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) ProductManagementService() ProductManagement.ProductManagementService {
	return getInstanceOfType(&applicationConnectionFactory.productManagementService, func() ProductManagement.ProductManagementService {
		return ProductManagement.NewService(applicationConnectionFactory.databaseConnectionFactory.ProductRepository(), applicationConnectionFactory.ProductService())
	})
}

func (applicationConnectionFactory *applicationConnectionFactory) ProductMapper() Product.ProductMapper {
	return getInstanceOfType(&applicationConnectionFactory.productMapper, func() Product.ProductMapper {
		return Product.NewProductMapper(applicationConnectionFactory.CategoryService(), applicationConnectionFactory.CmsService())
	})
}

func NewApplicationConnectionFactory(databaseConnectionFactory DatabaseConnectionFactory, brokerConnectionFactory MessageBrokerConnectionFactory) ApplicationConnectionFactory {
	applicationConnectionOnce.Do(func() {
		applicationConnectionFactoryInstance = &applicationConnectionFactory{
			databaseConnectionFactory: databaseConnectionFactory,
			messageBroker:             brokerConnectionFactory.MessageBroker(),
		}
	})
	return applicationConnectionFactoryInstance
}
