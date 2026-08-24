package router

import (
	Http "net/http"
	OS "os"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Logger "github.com/danyel/ecommerce/cmd/logger"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Factory "github.com/danyel/ecommerce/internal/common/factory/context"
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
	ById                         = "/{ID}"
	ByCode                       = "/{code}"
	ByLanguage                   = "/{language}"
)

// ApiRouter Definition the web layer.
// WebHandlerContextFactory will provide instances of web handlers to be used.
// ServerConfiguration will provide the application port to be used.
type ApiRouter interface {
	// Start the http server
	Start()
	// Router configuration of the loggers and api routing
	Router() *Router.Mux
}

// Router implementation of the Router method
func (apiRouter *apiRouter) Router() *Router.Mux {
	apiRouter.rootRouter = Router.NewRouter()
	apiRouter.configureLog()
	apiRouter.configureApiRouting()
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

// NewApiRouter Factory method for the ApiRouter interface
func NewApiRouter(serverConfiguration *Configuration.ServerConfiguration, webHandlerContextFactory Factory.WebHandlerContextFactory) ApiRouter {
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
	//router.Use(ApplicationMiddleware.JwtAuthMiddleware(apiRouter.ServerConfiguration.JwtSecret))
}

// configureApiRouting All API routing defined here
func (apiRouter *apiRouter) configureApiRouting() {
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
	return router.Route(ShoppingBasketRootContext, func(router Router.Router) {
		router.Route(VersionOne, func(router Router.Router) {
			router.Route(ShoppingBasketsRootContext, func(router Router.Router) {
				router.Post(BaseContextPath, shoppingBasketWebHandler.HandleCreateShoppingBasketV1)
				router.Route(ById, func(router Router.Router) {
					router.Post(BaseContextPath, shoppingBasketWebHandler.HandleUpdateShoppingBasketItemV1)
					router.Get(SLASH, shoppingBasketWebHandler.HandleGetShoppingBasketV1)
				})
			})
		})
	})
}

// productManagement Product Management api /api/product-management
func productManagement(router Router.Router, productManagementWebHandler ProductManagement.ProductManagementWebHandler) Router.Router {
	return router.Route(ProductManagementRootContext, func(router Router.Router) {
		router.Route(VersionOne, func(router Router.Router) {
			router.Route(ProductsRootContext, func(router Router.Router) {
				router.Get(SLASH, productManagementWebHandler.HandleGetProductsV1)
				router.Post(SLASH, productManagementWebHandler.HandleCreateProductV1)
				router.Route(ById, func(router Router.Router) {
					router.Get(SLASH, productManagementWebHandler.HandleGetProductV1)
					router.Delete(SLASH, productManagementWebHandler.HandleDeleteProductV1)
					router.Put(SLASH, productManagementWebHandler.HandleUpdateProductV1)
				})
			})
		})
	})
}

// category api /api/category
func category(router Router.Router, categoryWebHandler Category.CategoryWebHandler) Router.Router {
	return router.Route(CategoryRootContext, func(router Router.Router) {
		router.Route(VersionOne, func(router Router.Router) {
			router.Route(CategoriesRootContext, func(router Router.Router) {
				router.Post(BaseContextPath, categoryWebHandler.HandleCreateCategoryV1)
			})
			router.Post(TranslationsRootContext, categoryWebHandler.HandleCreateTranslationsV1)
		})
	})
}

// cms api /api/cms
func cms(router Router.Router, cmsWebHandler CMS.CmsWebHandler) Router.Router {
	return router.Route(CmsRootContext, func(router Router.Router) {
		router.Route(VersionOne, func(router Router.Router) {
			router.Route(TranslationsRootContext, func(router Router.Router) {
				router.Get(BaseContextPath, cmsWebHandler.HandleV1)
				router.Get(ByLanguage+ByCode, cmsWebHandler.HandleGetTranslationV1)
			})
		})
	})
}

// product api /api/product
func product(router Router.Router, productWebHandler Product.ProductWebHandler) Router.Router {
	return router.Route(ProductRootContext, func(router Router.Router) {
		router.Route(VersionOne, func(router Router.Router) {
			router.Get(SLASH, productWebHandler.HandleGetProductsV1)
			router.Route(ById, func(router Router.Router) {
				router.Get(BaseContextPath, productWebHandler.HandleGetProductV1)
			})
		})
	})
}

// management api /api/management
func management(router Router.Router, managementWebHandler Management.ManagementWebHandler) Router.Router {
	return router.Route(ManagementRootContext, func(router Router.Router) {
		router.Route(VersionOne, func(router Router.Router) {
			router.Route(CategoriesRootContext, func(router Router.Router) {
				router.Get(BaseContextPath, managementWebHandler.HandleGetCategoriesV1)
			})
			router.Post(TranslationsRootContext, managementWebHandler.HandleCreateTranslationsV1)
		})
	})
}
