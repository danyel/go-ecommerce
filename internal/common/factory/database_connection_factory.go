package factory

import (
	Configuration "github.com/danyel/ecommerce/cmd/config"
	DatabaseConnection "github.com/danyel/ecommerce/cmd/database"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Repository "github.com/danyel/ecommerce/internal/common/repository"
	Product "github.com/danyel/ecommerce/internal/product"
	Reservation "github.com/danyel/ecommerce/internal/reservation"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	Database "gorm.io/gorm"
)

type DatabaseConnectionFactory interface {
	ProductRepository() Repository.CrudRepository[Product.ProductModel]
	CmsRepository() Repository.CrudRepository[CMS.CmsModel]
	CategoryRepository() Repository.CrudRepository[Category.CategoryModel]
	ShoppingBasketRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel]
	ShoppingBasketItemRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel]
	ReservationRepository() Repository.CrudRepository[Reservation.ReservationModel]
}

type databaseConnectionFactory struct {
	databaseConnection           *Database.DB
	productRepository            Repository.CrudRepository[Product.ProductModel]
	cmsRepository                Repository.CrudRepository[CMS.CmsModel]
	categoryRepository           Repository.CrudRepository[Category.CategoryModel]
	shoppingBasketRepository     Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel]
	shoppingBasketItemRepository Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel]
	reservationRepository        Repository.CrudRepository[Reservation.ReservationModel]
}

func (databaseConnectionFactory *databaseConnectionFactory) ProductRepository() Repository.CrudRepository[Product.ProductModel] {
	return getInstanceOfType(&databaseConnectionFactory.productRepository, func() Repository.CrudRepository[Product.ProductModel] {
		return Repository.NewCrudRepository[Product.ProductModel](databaseConnectionFactory.databaseConnection)
	})
}

func (databaseConnectionFactory *databaseConnectionFactory) CmsRepository() Repository.CrudRepository[CMS.CmsModel] {
	return getInstanceOfType(&databaseConnectionFactory.cmsRepository, func() Repository.CrudRepository[CMS.CmsModel] {
		return Repository.NewCrudRepository[CMS.CmsModel](databaseConnectionFactory.databaseConnection)
	})
}

func (databaseConnectionFactory *databaseConnectionFactory) CategoryRepository() Repository.CrudRepository[Category.CategoryModel] {
	return getInstanceOfType(&databaseConnectionFactory.categoryRepository, func() Repository.CrudRepository[Category.CategoryModel] {
		return Repository.NewCrudRepository[Category.CategoryModel](databaseConnectionFactory.databaseConnection)
	})
}

func (databaseConnectionFactory *databaseConnectionFactory) ShoppingBasketRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel] {
	return getInstanceOfType(&databaseConnectionFactory.shoppingBasketRepository, func() Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel] {
		return Repository.NewCrudRepository[ShoppingBasket.ShoppingBasketModel](databaseConnectionFactory.databaseConnection)
	})
}

func (databaseConnectionFactory *databaseConnectionFactory) ShoppingBasketItemRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel] {
	return getInstanceOfType(&databaseConnectionFactory.shoppingBasketItemRepository, func() Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel] {
		return Repository.NewCrudRepository[ShoppingBasket.ShoppingBasketItemModel](databaseConnectionFactory.databaseConnection)
	})
}
func (databaseConnectionFactory *databaseConnectionFactory) ReservationRepository() Repository.CrudRepository[Reservation.ReservationModel] {
	return getInstanceOfType(&databaseConnectionFactory.reservationRepository, func() Repository.CrudRepository[Reservation.ReservationModel] {
		return Repository.NewCrudRepository[Reservation.ReservationModel](databaseConnectionFactory.databaseConnection)
	})
}

func NewDatabaseConnectionFactory(databaseConfiguration *Configuration.DatabaseConfiguration) DatabaseConnectionFactory {
	databaseConnectionOnce.Do(func() {
		databaseConnection, err := DatabaseConnection.Connect(databaseConfiguration)
		if err != nil {
			panic(err)
		}
		databaseConnectionFactoryInstance = &databaseConnectionFactory{
			databaseConnection: databaseConnection,
		}
	})
	return databaseConnectionFactoryInstance
}
