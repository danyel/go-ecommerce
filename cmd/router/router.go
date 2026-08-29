package router

import (
	Http "net/http"
	OS "os"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Factory "github.com/danyel/ecommerce/cmd/factory/context"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	ApplicationMiddleware "github.com/danyel/ecommerce/cmd/middleware"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	Router "github.com/go-chi/chi/v5"
	Middleware "github.com/go-chi/chi/v5/middleware"
)

const (
	SLASH                        = "/"
	BaseContextPath              = "/api"
	ShoppingBasketRootContext    = "/shopping-basket"
	ShoppingBasketsRootContext   = "/shopping-baskets"
	ProductRootContext           = "/product"
	ProductsRootContext          = "/products"
	ProductManagementRootContext = "/product-management"
	CmsRootContext               = "/cms"
	ManagementRootContext        = "/management"
	VersionOne                   = "/v1"
	CategoryRootContext          = "/category"
	CategoriesRootContext        = "/categories"
	TranslationsRootContext      = "/translations"
	ByID                         = "/{ID}"
	ByCode                       = "/{code}"
	ByLanguage                   = "/{language}"
)

// APIRouter Definition the web layer.
// WebHandlerContextFactory will provide instances of web handlers to be used.
// ServerConfiguration will provide the application port to be used.
type APIRouter interface {
	// Start the http server
	Start()
	// Router configuration of the loggers and api routing
	Router() *Router.Mux
}

// Router implementation of the Router method
func (apiRouter *apiRouter) Router() *Router.Mux {
	apiRouter.rootRouter = Router.NewRouter()
	apiRouter.configureLog()
	apiRouter.configureAPIRouting()
	return apiRouter.rootRouter
}

// Start implementation of the Start method
func (apiRouter *apiRouter) Start() {
	Logger.Log.Info("Running the server on port %s", apiRouter.serverConfiguration.Addr)
	if err := Http.ListenAndServe(apiRouter.serverConfiguration.Addr, apiRouter.Router()); err != nil {
		Logger.Log.Fatal(err)
		OS.Exit(0)
	}
}

// NewAPIRouter Factory method for the ApiRouter interface
func NewAPIRouter(serverConfiguration *Configuration.ServerConfiguration, webHandlerContextFactory Factory.WebHandlerContextFactory) APIRouter {
	apiRouter := &apiRouter{
		serverConfiguration:      serverConfiguration,
		webHandlerContextFactory: webHandlerContextFactory,
	}
	return apiRouter
}

// apiRouter instance of the ApiRouter
type apiRouter struct {
	serverConfiguration      *Configuration.ServerConfiguration
	webHandlerContextFactory Factory.WebHandlerContextFactory
	rootRouter               *Router.Mux
}

// configureLog all configuration for logging is configured here
func (apiRouter *apiRouter) configureLog() {
	apiRouter.rootRouter.Use(Middleware.RequestID)
	apiRouter.rootRouter.Use(Middleware.Logger)
	apiRouter.rootRouter.Use(Middleware.Recoverer)
	apiRouter.rootRouter.Use(ApplicationMiddleware.CorrelationIDMiddleware)
	//apiRouter.Use(ApplicationMiddleware.JwtAuthMiddleware(apiRouter.ServerConfiguration.JwtSecret))
}

// configureAPIRouting All API routing defined here
func (apiRouter *apiRouter) configureAPIRouting() {
	webHandlerContextFactory := apiRouter.webHandlerContextFactory
	apiRouter.rootRouter.Route(BaseContextPath, func(router Router.Router) {
		product(router, webHandlerContextFactory.ProductWebHandler())
		category(router, webHandlerContextFactory.CategoryWebHandler())
		productManagement(router, webHandlerContextFactory.ProductManagementWebHandler())
		management(router, webHandlerContextFactory.ManagementWebHandler())
		cms(router, webHandlerContextFactory.CmsWebHandler())
		shoppingBasket(router, webHandlerContextFactory.ShoppingBasketWebHandler())
	})
}

// shoppingBasket Shopping Basket api /api/shopping-basket
func shoppingBasket(router Router.Router, shoppingBasketWebHandler ShoppingBasket.ShoppingBasketWebHandler) Router.Router {
	return router.Route(ShoppingBasketRootContext, func(shoppingBasketRouter Router.Router) {
		shoppingBasketRouter.Route(VersionOne, func(versionOneRouter Router.Router) {
			versionOneRouter.Route(ShoppingBasketsRootContext, func(shoppingBasketsRouter Router.Router) {
				shoppingBasketsRouter.Post(SLASH, shoppingBasketWebHandler.HandleCreateShoppingBasketV1)
				shoppingBasketsRouter.Route(ByID, func(byIdRouter Router.Router) {
					Logger.Log.Debug("Shopping Basket By Id")
					byIdRouter.Get(SLASH, shoppingBasketWebHandler.HandleGetShoppingBasketByIDV1)
					byIdRouter.Put(SLASH, shoppingBasketWebHandler.HandleUpdateShoppingBasketItemV1)
				})
			})
		})
	})
}

// productManagement Product Management api /api/product-management
func productManagement(router Router.Router, productManagementWebHandler ProductManagement.ProductManagementWebHandler) Router.Router {
	return router.Route(ProductManagementRootContext, func(productManagementRootRouter Router.Router) {
		productManagementRootRouter.Route(VersionOne, func(versionOneRouter Router.Router) {
			versionOneRouter.Route(ProductsRootContext, func(productsRootRouter Router.Router) {
				productsRootRouter.Get(SLASH, productManagementWebHandler.HandleGetProductsV1)
				productsRootRouter.Post(SLASH, productManagementWebHandler.HandleCreateProductV1)
				productsRootRouter.Route(ByID, func(byIdRouter Router.Router) {
					byIdRouter.Get(SLASH, productManagementWebHandler.HandleGetProductV1)
					byIdRouter.Delete(SLASH, productManagementWebHandler.HandleDeleteProductV1)
					byIdRouter.Put(SLASH, productManagementWebHandler.HandleUpdateProductV1)
				})
			})
		})
	})
}

// category api /api/category
func category(router Router.Router, categoryWebHandler Category.CategoryWebHandler) Router.Router {
	return router.Route(CategoryRootContext, func(categoryRootRouter Router.Router) {
		categoryRootRouter.Route(VersionOne, func(versionOneRouter Router.Router) {
			versionOneRouter.Route(CategoriesRootContext, func(categoriesRootRouter Router.Router) {
				categoriesRootRouter.Post(BaseContextPath, categoryWebHandler.HandleCreateCategoryV1)
			})
			versionOneRouter.Post(TranslationsRootContext, categoryWebHandler.HandleCreateTranslationsV1)
		})
	})
}

// cms api /api/cms
func cms(router Router.Router, cmsWebHandler CMS.CmsWebHandler) Router.Router {
	return router.Route(CmsRootContext, func(cmsRootRouter Router.Router) {
		cmsRootRouter.Route(VersionOne, func(versionOneRouter Router.Router) {
			versionOneRouter.Route(TranslationsRootContext, func(translationsRootRouter Router.Router) {
				Logger.Log.Debug("Hitting Url: %s", BaseContextPath+CmsRootContext+VersionOne+TranslationsRootContext)
				translationsRootRouter.Get(ByLanguage, cmsWebHandler.HandleV1)
				translationsRootRouter.Get(SLASH, cmsWebHandler.HandleV1)
				translationsRootRouter.Get(ByLanguage+ByCode, cmsWebHandler.HandleGetTranslationV1)
			})
		})
	})
}

// product api /api/product
func product(router Router.Router, productWebHandler Product.ProductWebHandler) Router.Router {
	return router.Route(ProductRootContext, func(productRootRouter Router.Router) {
		productRootRouter.Route(VersionOne, func(versionOneRouter Router.Router) {
			versionOneRouter.Get(SLASH, productWebHandler.HandleGetProductsV1)
			versionOneRouter.Route(ByID, func(byIdRouter Router.Router) {
				byIdRouter.Get(SLASH, productWebHandler.HandleGetProductV1)
			})
		})
	})
}

// management api /api/management
func management(router Router.Router, managementWebHandler Management.ManagementWebHandler) Router.Router {
	return router.Route(ManagementRootContext, func(managementRootRouter Router.Router) {
		managementRootRouter.Route(VersionOne, func(versionOneRouter Router.Router) {
			versionOneRouter.Route(CategoriesRootContext, func(categoriesRootRouter Router.Router) {
				categoriesRootRouter.Get(BaseContextPath, managementWebHandler.HandleGetCategoriesV1)
			})
			versionOneRouter.Post(TranslationsRootContext, managementWebHandler.HandleCreateTranslationsV1)
		})
	})
}
