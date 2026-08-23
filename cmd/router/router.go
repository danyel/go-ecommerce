package router

import (
	Log "log"
	Http "net/http"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Factory "github.com/danyel/ecommerce/internal/common/factory"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/productmanagement"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shoppingbasket"
	Router "github.com/go-chi/chi/v5"
	Middleware "github.com/go-chi/chi/v5/middleware"
)

type APIDefinition struct {
	ServerConfiguration          *Configuration.ServerConfiguration
	ApplicationConnectionFactory Factory.ApplicationConnectionFactory
}

func (apiDefinition *APIDefinition) ConfigRouter() *Router.Mux {
	router := Router.NewRouter()
	router.Use(Middleware.RequestID)
	router.Use(Middleware.Logger)
	router.Use(Middleware.Recoverer)
	// r.Use(ApplicationMiddleware.JwtAuthMiddleware(a.SC.JwtSecret))
	applicationConnectionFactory := apiDefinition.ApplicationConnectionFactory
	router.Route("/api", func(router Router.Router) {
		productV1Routing(router, applicationConnectionFactory)
		categoryV1Routing(router, applicationConnectionFactory)
		productManagementV1Routing(router, applicationConnectionFactory)
		managementV1Routing(router, applicationConnectionFactory)
		cmsV1Routing(router, applicationConnectionFactory)
		paymentV1Routing(router, applicationConnectionFactory)
		orderV1Routing(router, applicationConnectionFactory)
		shoppingBasketV1Routing(router, applicationConnectionFactory)
	})
	return router
}

// Shopping Basket api V1 /api/shopping-basket/v1/shopping-baskets
func shoppingBasketV1Routing(router Router.Router, applicationConnectionFactory Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/shopping-basket/v1/shopping-baskets", func(router Router.Router) {
		shoppingBasketHandler := ShoppingBasket.NewHandler(applicationConnectionFactory.ShoppingBasketService())
		router.Post("/", shoppingBasketHandler.CreateShoppingBasket)
		router.Route("/{shoppingBasketID}", func(router Router.Router) {
			router.Post("/", shoppingBasketHandler.UpdateShoppingBasketItem)
			router.Get("/", shoppingBasketHandler.GetShoppingBasket)
		})
	})
}

// Product Management api V1 /api/product-management/v1
func productManagementV1Routing(router Router.Router, applicationConnectionFactory Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/product-management/v1", func(router Router.Router) {
		router.Route("/products", func(router Router.Router) {
			productManagementHandler := ProductManagement.NewHandler(applicationConnectionFactory.CategoryService(), applicationConnectionFactory.CmsService(), applicationConnectionFactory.ProductManagementService())
			router.Get("/", productManagementHandler.GetProducts)
			router.Post("/", productManagementHandler.CreateProduct)
			router.Route("/{productID}", func(router Router.Router) {
				router.Get("/", productManagementHandler.GetProduct)
				router.Delete("/", productManagementHandler.DeleteProduct)
				router.Put("/", productManagementHandler.UpdateProduct)
			})
		})
	})
}

// Order api V1 /api/order/v1/orders
func orderV1Routing(router Router.Router, _ Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/order/v1/orders", func(router Router.Router) {
		//
	})
}

// Category api V1 /api/category/v1/categories
func categoryV1Routing(router Router.Router, applicationConnectionFactory Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/category/v1", func(router Router.Router) {
		categoryHandler := Category.NewHandler(applicationConnectionFactory.CategoryService())
		router.Route("/categories", func(router Router.Router) {
			router.Post("/", categoryHandler.CreateCategory)
		})
		router.Post("/translations", categoryHandler.CreateTranslations)
	})
}

// Payment api V1 /api/payment/v1/payments
func paymentV1Routing(router Router.Router, _ Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/payment/v1/payments", func(router Router.Router) {
		//
	})
}

// CMS api V1 /api/cms/v1/translations
func cmsV1Routing(router Router.Router, applicationConnectionFactory Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/cms/v1/translations", func(router Router.Router) {
		cmsHandler := CMS.NewHandler(applicationConnectionFactory.CmsService())
		router.Get("/", cmsHandler.GetTranslations)
		router.Get("/{language}/{ID}", cmsHandler.GetTranslation)
	})
}

// Management api V1 /api/management/v1
func managementV1Routing(router Router.Router, applicationConnectionFactory Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/management/v1", func(router Router.Router) {
		managementHandler := Management.NewHandler(applicationConnectionFactory.CategoryService(), applicationConnectionFactory.ManagementService(), applicationConnectionFactory.CmsService())
		router.Route("/categories", func(router Router.Router) {
			router.Get("/", managementHandler.GetCategories)
		})
		router.Post("/translations", managementHandler.CreateTranslations)
	})
}

// Product api V1 /api/product/v1/products
func productV1Routing(router Router.Router, applicationConnectionFactory Factory.ApplicationConnectionFactory) Router.Router {
	return router.Route("/product/v1/products", func(router Router.Router) {
		productHandler := Product.NewHandler(applicationConnectionFactory.ProductService())
		router.Get("/", productHandler.GetProducts)
		// /api/product/v1/products/{productID}
		router.Route("/{productID}", func(router Router.Router) {
			router.Get("/", productHandler.GetProduct)
		})
	})
}

func (apiDefinition *APIDefinition) Run(router *Router.Mux) {
	Log.Printf("Running the server on port %s", apiDefinition.ServerConfiguration.Addr)
	if err := Http.ListenAndServe(apiDefinition.ServerConfiguration.Addr, router); err != nil {
		Log.Fatal(err)
	}
}
