package contextfactory

import (
	WebHandler "github.com/danyel/ecommerce/internal/handler"
)

type WebHandlerContextFactory interface {
	ShoppingBasketWebHandler() WebHandler.ShoppingBasketWebHandler
	ProductManagementWebHandler() WebHandler.ProductManagementWebHandler
	CategoryWebHandler() WebHandler.CategoryWebHandler
	CmsWebHandler() WebHandler.CmsWebHandler
	ManagementWebHandler() WebHandler.ManagementWebHandler
	ProductWebHandler() WebHandler.IProductWebHandler
}

type webHandlerContextFactory struct {
	shoppingBasketWebHandler    WebHandler.ShoppingBasketWebHandler
	productManagementWebHandler WebHandler.ProductManagementWebHandler
	productWebHandler           WebHandler.IProductWebHandler
	categoryWebHandler          WebHandler.CategoryWebHandler
	cmsWebHandler               WebHandler.CmsWebHandler
	managementWebHandler        WebHandler.ManagementWebHandler
}

func (webHandlerContextFactory *webHandlerContextFactory) ShoppingBasketWebHandler() WebHandler.ShoppingBasketWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.shoppingBasketWebHandler, func() WebHandler.ShoppingBasketWebHandler {
		return WebHandler.NewShoppingBasketWebHandler(applicationContextFactoryInstance.ShoppingBasketService(), applicationContextFactoryInstance.ShoppingBasketValidator())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) ProductManagementWebHandler() WebHandler.ProductManagementWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.productManagementWebHandler, func() WebHandler.ProductManagementWebHandler {
		return WebHandler.NewProductManagementWebHandler(applicationContextFactoryInstance.CategoryService(), applicationContextFactoryInstance.CmsService(), applicationContextFactoryInstance.ProductManagementService(), applicationContextFactoryInstance.ProductMapper(), applicationContextFactoryInstance.CategoryMapper())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) ProductWebHandler() WebHandler.IProductWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.productWebHandler, func() WebHandler.IProductWebHandler {
		return WebHandler.ProductWebHandler(applicationContextFactoryInstance.ProductService(), applicationContextFactoryInstance.ProductMapper(), applicationContextFactoryInstance.CategoryMapper())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) CategoryWebHandler() WebHandler.CategoryWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.categoryWebHandler, func() WebHandler.CategoryWebHandler {
		return WebHandler.NewCategoryWebHandler(applicationContextFactoryInstance.CategoryService())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) CmsWebHandler() WebHandler.CmsWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.cmsWebHandler, func() WebHandler.CmsWebHandler {
		return WebHandler.NewCmsWebHandler(applicationContextFactoryInstance.CmsService())
	})
}

func (webHandlerContextFactory *webHandlerContextFactory) ManagementWebHandler() WebHandler.ManagementWebHandler {
	return getInstanceOfType(&webHandlerContextFactory.managementWebHandler, func() WebHandler.ManagementWebHandler {
		return WebHandler.NewManagementWebHandler(applicationContextFactoryInstance.CategoryService(), applicationContextFactoryInstance.ManagementService(), applicationContextFactoryInstance.CmsService())
	})
}

func BuildAWebHandlerContextFactory() WebHandlerContextFactory {
	webHandlerContextFactoryOnce.Do(func() {
		webHandlerContextFactoryInstance = &webHandlerContextFactory{}
	})
	return webHandlerContextFactoryInstance
}
