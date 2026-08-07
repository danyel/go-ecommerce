package router

import (
	Log "log"
	Http "net/http"

	Configuration "github.com/danyel/ecommerce/cmd/config"
	Category "github.com/danyel/ecommerce/internal/category"
	CMS "github.com/danyel/ecommerce/internal/cms"
	Port "github.com/danyel/ecommerce/internal/common/port"
	Management "github.com/danyel/ecommerce/internal/management"
	Product "github.com/danyel/ecommerce/internal/product"
	ProductManagement "github.com/danyel/ecommerce/internal/product-management"
	ShoppingBasket "github.com/danyel/ecommerce/internal/shopping-basket"
	Router "github.com/go-chi/chi/v5"
	Middleware "github.com/go-chi/chi/v5/middleware"
	Database "gorm.io/gorm"
)

type ApiDefinition struct {
	SC             *Configuration.ServerConfiguration
	DB             *Database.DB
	EventPublisher Port.EventPublisher
}

func (a *ApiDefinition) ConfigRouter() *Router.Mux {
	r := Router.NewRouter()
	r.Use(Middleware.RequestID)
	r.Use(Middleware.Logger)
	r.Use(Middleware.Recoverer)

	r.Route("/api", func(r Router.Router) {
		productV1Routing(r, a)
		categoryV1Routing(r, a)
		productManagementV1Routing(r, a)
		managementV1Routing(r, a)
		cmsV1Routing(r, a)
		paymentV1Routing(r, a)
		orderV1Routing(r, a)
		shoppingBasketV1Routing(r, a)
	})
	return r
}

// Shopping Basket api V1 /api/shopping-basket/v1/shopping-baskets
func shoppingBasketV1Routing(r Router.Router, a *ApiDefinition) Router.Router {
	return r.Route("/shopping-basket/v1/shopping-baskets", func(r Router.Router) {
		h := ShoppingBasket.NewHandler(a.DB, a.EventPublisher)
		r.Post("/", h.CreateShoppingBasket)
		r.Route("/{shoppingBasketId}", func(r Router.Router) {
			r.Post("/", h.UpdateShoppingBasketItem)
			r.Get("/", h.GetShoppingBasket)
		})
	})
}

// Product Management api V1 /api/product-management/v1
func productManagementV1Routing(r Router.Router, a *ApiDefinition) Router.Router {
	return r.Route("/product-management/v1", func(r Router.Router) {
		r.Route("/products", func(r Router.Router) {
			h := ProductManagement.NewHandler(a.DB)
			r.Get("/", h.GetProducts)
			r.Post("/", h.CreateProduct)
			r.Route("/{productId}", func(r Router.Router) {
				r.Get("/", h.GetProduct)
				r.Delete("/", h.DeleteProduct)
				r.Put("/", h.UpdateProduct)
			})
		})
	})
}

// Order api V1 /api/order/v1/orders
func orderV1Routing(r Router.Router, _ *ApiDefinition) Router.Router {
	return r.Route("/order/v1/orders", func(r Router.Router) {
		//
	})
}

// Category api V1 /api/category/v1/categories
func categoryV1Routing(r Router.Router, a *ApiDefinition) Router.Router {
	return r.Route("/category/v1", func(r Router.Router) {
		h := Category.NewHandler(a.DB)
		r.Route("/categories", func(r Router.Router) {
			r.Post("/", h.CreateCategory)
		})
		r.Post("/translations", h.CreateTranslations)
	})
}

// Payment api V1 /api/payment/v1/payments
func paymentV1Routing(r Router.Router, _ *ApiDefinition) Router.Router {
	return r.Route("/payment/v1/payments", func(r Router.Router) {
		//
	})
}

// CMS api V1 /api/cms/v1/translations
func cmsV1Routing(r Router.Router, a *ApiDefinition) Router.Router {
	return r.Route("/cms/v1/translations", func(r Router.Router) {
		h := CMS.NewHandler(a.DB)
		r.Get("/", h.GetTranslations)
		r.Get("/{language}/{id}", h.GetTranslation)
	})
}

// Management api V1 /api/management/v1
func managementV1Routing(r Router.Router, a *ApiDefinition) Router.Router {
	return r.Route("/management/v1", func(r Router.Router) {
		h := Management.NewHandler(a.DB)
		r.Route("/categories", func(r Router.Router) {
			r.Get("/", h.GetCategories)
		})
		r.Post("/translations", h.CreateTranslations)
	})
}

// Product api V1 /api/product/v1/products
func productV1Routing(r Router.Router, a *ApiDefinition) Router.Router {
	return r.Route("/product/v1/products", func(r Router.Router) {
		h := Product.NewApiHandler(Product.NewProductService(a.DB))
		r.Get("/", h.GetProducts)
		// /api/product/v1/products/{productId}
		r.Route("/{productId}", func(r Router.Router) {
			r.Get("/", h.GetProduct)
		})
	})
}

func (a *ApiDefinition) Run(r *Router.Mux) {
	Log.Printf("Running the server on port %s", a.SC.Addr)
	if err := Http.ListenAndServe(a.SC.Addr, r); err != nil {
		Log.Fatal(err)
	}

}
