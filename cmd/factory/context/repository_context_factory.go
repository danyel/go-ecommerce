package contextfactory

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

type RepositoryContextFactory interface {
	ProductRepository() Repository.CrudRepository[Product.ProductModel]
	CmsRepository() Repository.CrudRepository[CMS.CmsModel]
	CategoryRepository() Repository.CrudRepository[Category.CategoryModel]
	ShoppingBasketRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel]
	ShoppingBasketItemRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel]
	ReservationRepository() Repository.CrudRepository[Reservation.ReservationModel]
}

type repositoryContextFactory struct {
	databaseConnection           *Database.DB
	productRepository            Repository.CrudRepository[Product.ProductModel]
	cmsRepository                Repository.CrudRepository[CMS.CmsModel]
	categoryRepository           Repository.CrudRepository[Category.CategoryModel]
	shoppingBasketRepository     Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel]
	shoppingBasketItemRepository Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel]
	reservationRepository        Repository.CrudRepository[Reservation.ReservationModel]
}

func (repositoryContextFactory *repositoryContextFactory) ProductRepository() Repository.CrudRepository[Product.ProductModel] {
	return getInstanceOfType(&repositoryContextFactory.productRepository, func() Repository.CrudRepository[Product.ProductModel] {
		return Repository.NewCrudRepository[Product.ProductModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) CmsRepository() Repository.CrudRepository[CMS.CmsModel] {
	return getInstanceOfType(&repositoryContextFactory.cmsRepository, func() Repository.CrudRepository[CMS.CmsModel] {
		return Repository.NewCrudRepository[CMS.CmsModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) CategoryRepository() Repository.CrudRepository[Category.CategoryModel] {
	return getInstanceOfType(&repositoryContextFactory.categoryRepository, func() Repository.CrudRepository[Category.CategoryModel] {
		return Repository.NewCrudRepository[Category.CategoryModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) ShoppingBasketRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel] {
	return getInstanceOfType(&repositoryContextFactory.shoppingBasketRepository, func() Repository.CrudRepository[ShoppingBasket.ShoppingBasketModel] {
		return Repository.NewCrudRepository[ShoppingBasket.ShoppingBasketModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) ShoppingBasketItemRepository() Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel] {
	return getInstanceOfType(&repositoryContextFactory.shoppingBasketItemRepository, func() Repository.CrudRepository[ShoppingBasket.ShoppingBasketItemModel] {
		return Repository.NewCrudRepository[ShoppingBasket.ShoppingBasketItemModel](repositoryContextFactory.databaseConnection)
	})
}
func (repositoryContextFactory *repositoryContextFactory) ReservationRepository() Repository.CrudRepository[Reservation.ReservationModel] {
	return getInstanceOfType(&repositoryContextFactory.reservationRepository, func() Repository.CrudRepository[Reservation.ReservationModel] {
		return Repository.NewCrudRepository[Reservation.ReservationModel](repositoryContextFactory.databaseConnection)
	})
}

func InitializeDatabaseContextFactory() {
	databaseContextFactoryOnce.Do(func() {
		databaseConnection, err := DatabaseConnection.Connect(Configuration.Database())
		if err != nil {
			panic(err)
		}
		repositoryContextFactoryInstance = &repositoryContextFactory{
			databaseConnection: databaseConnection,
		}
	})
}
