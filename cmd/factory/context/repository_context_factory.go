package contextfactory

import (
	Configuration "github.com/danyel/ecommerce/cmd/config"
	DatabaseConnection "github.com/danyel/ecommerce/cmd/database"
	Category "github.com/danyel/ecommerce/internal/persistence"
	Database "gorm.io/gorm"
)

type RepositoryContextFactory interface {
	ProductRepository() Category.CrudRepository[Category.ProductModel]
	CmsRepository() Category.CrudRepository[Category.CmsModel]
	CategoryRepository() Category.CrudRepository[Category.CategoryModel]
	ShoppingBasketRepository() Category.CrudRepository[Category.ShoppingBasketModel]
	ShoppingBasketItemRepository() Category.CrudRepository[Category.ShoppingBasketItemModel]
	ReservationRepository() Category.CrudRepository[Category.ReservationModel]
}

type repositoryContextFactory struct {
	databaseConnection           *Database.DB
	productRepository            Category.CrudRepository[Category.ProductModel]
	cmsRepository                Category.CrudRepository[Category.CmsModel]
	categoryRepository           Category.CrudRepository[Category.CategoryModel]
	shoppingBasketRepository     Category.CrudRepository[Category.ShoppingBasketModel]
	shoppingBasketItemRepository Category.CrudRepository[Category.ShoppingBasketItemModel]
	reservationRepository        Category.CrudRepository[Category.ReservationModel]
}

func (repositoryContextFactory *repositoryContextFactory) ProductRepository() Category.CrudRepository[Category.ProductModel] {
	return getInstanceOfType(&repositoryContextFactory.productRepository, func() Category.CrudRepository[Category.ProductModel] {
		return Category.NewCrudRepository[Category.ProductModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) CmsRepository() Category.CrudRepository[Category.CmsModel] {
	return getInstanceOfType(&repositoryContextFactory.cmsRepository, func() Category.CrudRepository[Category.CmsModel] {
		return Category.NewCrudRepository[Category.CmsModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) CategoryRepository() Category.CrudRepository[Category.CategoryModel] {
	return getInstanceOfType(&repositoryContextFactory.categoryRepository, func() Category.CrudRepository[Category.CategoryModel] {
		return Category.NewCrudRepository[Category.CategoryModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) ShoppingBasketRepository() Category.CrudRepository[Category.ShoppingBasketModel] {
	return getInstanceOfType(&repositoryContextFactory.shoppingBasketRepository, func() Category.CrudRepository[Category.ShoppingBasketModel] {
		return Category.NewCrudRepository[Category.ShoppingBasketModel](repositoryContextFactory.databaseConnection)
	})
}

func (repositoryContextFactory *repositoryContextFactory) ShoppingBasketItemRepository() Category.CrudRepository[Category.ShoppingBasketItemModel] {
	return getInstanceOfType(&repositoryContextFactory.shoppingBasketItemRepository, func() Category.CrudRepository[Category.ShoppingBasketItemModel] {
		return Category.NewCrudRepository[Category.ShoppingBasketItemModel](repositoryContextFactory.databaseConnection)
	})
}
func (repositoryContextFactory *repositoryContextFactory) ReservationRepository() Category.CrudRepository[Category.ReservationModel] {
	return getInstanceOfType(&repositoryContextFactory.reservationRepository, func() Category.CrudRepository[Category.ReservationModel] {
		return Category.NewCrudRepository[Category.ReservationModel](repositoryContextFactory.databaseConnection)
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
