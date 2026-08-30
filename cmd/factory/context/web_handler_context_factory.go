package contextfactory

import (
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	Shoppingbasket "github.com/danyel/ecommerce/internal/shoppingbasket"
)

type WebHandlerContextFactory interface {
	ShoppingBasketWebHandler() Shoppingbasket.ShoppingBasketWebHandler
	ProductManagementWebHandler() ProductManagement.ProductManagementWebHandler
	CategoryWebHandler() Category.CategoryWebHandler
	CmsWebHandler() CMS.CmsWebHandler
	ManagementWebHandler() Management.ManagementWebHandler
	ProductWebHandler() Product.ProductWebHandler
}

type webHandlerContextFactory struct {
	shoppingBasketWebHandler    Shoppingbasket.ShoppingBasketWebHandler
	productManagementWebHandler ProductManagement.ProductManagementWebHandler
	productWebHandler           Product.ProductWebHandler
	categoryWebHandler          Category.CategoryWebHandler
	cmsWebHandler               CMS.CmsWebHandler
	managementWebHandler        Management.ManagementWebHandler
}

func (webHandlerContextFactory *webHandlerContextFactory) ShoppingBasketWebHandler() Shoppingbasket.ShoppingBasketWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.shoppingBasketWebHandler, func() Shoppingbasket.ShoppingBasketWebHandler {
		return Shoppingbasket.NewWebHandler(applicationContextFactoryInstance.ShoppingBasketService(), applicationContextFactoryInstance.ShoppingBasketValidator())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) ProductManagementWebHandler() ProductManagement.ProductManagementWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.productManagementWebHandler, func() ProductManagement.ProductManagementWebHandler {
		return ProductManagement.NewWebHandler(applicationContextFactoryInstance.CategoryService(), applicationContextFactoryInstance.CmsService(), applicationContextFactoryInstance.ProductManagementService())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) ProductWebHandler() Product.ProductWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.productWebHandler, func() Product.ProductWebHandler {
		return Product.NewWebHandler(applicationContextFactoryInstance.ProductService())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) CategoryWebHandler() Category.CategoryWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.categoryWebHandler, func() Category.CategoryWebHandler {
		return Category.NewWebHandler(applicationContextFactoryInstance.CategoryService())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) CmsWebHandler() CMS.CmsWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.cmsWebHandler, func() CMS.CmsWebHandler {
		return CMS.NewWebHandler(applicationContextFactoryInstance.CmsService())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) ManagementWebHandler() Management.ManagementWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.managementWebHandler, func() Management.ManagementWebHandler {
		return Management.NewWebHandler(applicationContextFactoryInstance.CategoryService(), applicationContextFactoryInstance.ManagementService(), applicationContextFactoryInstance.CmsService())
	})
}

func BuildAWebHandlerContextFactory() WebHandlerContextFactory {
	webHandlerContextFactoryOnce.Do(func() {
		webHandlerContextFactoryInstance = &webHandlerContextFactory{}
	})
	return webHandlerContextFactoryInstance
}
