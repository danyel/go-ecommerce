package router

import (
	"log"
	"net/http"

	"github.com/danyel/ecommerce/cmd/broker"
	"github.com/danyel/ecommerce/cmd/config"
	"github.com/danyel/ecommerce/internal/category"
	"github.com/danyel/ecommerce/internal/cms"
	"github.com/danyel/ecommerce/internal/management"
	"github.com/danyel/ecommerce/internal/product"
	"github.com/danyel/ecommerce/internal/product-management"
	shoppingbasket "github.com/danyel/ecommerce/internal/shopping-basket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type ApiDefinition struct {
	SC     *config.ServerConfiguration
	DB     *gorm.DB
	Broker *broker.Broker
}

func (a *ApiDefinition) ConfigRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
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
func shoppingBasketV1Routing(r chi.Router, a *ApiDefinition) chi.Router {
	return r.Route("/shopping-basket/v1/shopping-baskets", func(r chi.Router) {
		h := shoppingbasket.NewHandler(a.DB)
		r.Post("/", h.CreateShoppingBasket)
		r.Route("/{shoppingBasketId}", func(r chi.Router) {
			r.Post("/", h.UpdateShoppingBasketItem)
			r.Get("/", h.GetShoppingBasket)
		})
	})
}

// Product Management api V1 /api/product-management/v1
func productManagementV1Routing(r chi.Router, a *ApiDefinition) chi.Router {
	return r.Route("/product-management/v1", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			h := product_management.NewHandler(a.DB)
			r.Get("/", h.GetProducts)
			r.Post("/", h.CreateProduct)
			r.Route("/{productId}", func(r chi.Router) {
				r.Get("/", h.GetProduct)
				r.Delete("/", h.DeleteProduct)
				r.Put("/", h.UpdateProduct)
			})
		})
	})
}

// Order api V1 /api/order/v1/orders
func orderV1Routing(r chi.Router, _ *ApiDefinition) chi.Router {
	return r.Route("/order/v1/orders", func(r chi.Router) {
		//
	})
}

// Category api V1 /api/category/v1/categories
func categoryV1Routing(r chi.Router, a *ApiDefinition) chi.Router {
	return r.Route("/category/v1", func(r chi.Router) {
		h := category.NewHandler(a.DB)
		r.Route("/categories", func(r chi.Router) {
			r.Post("/", h.CreateCategory)
		})
		r.Post("/translations", h.CreateTranslations)
	})
}

// Payment api V1 /api/payment/v1/payments
func paymentV1Routing(r chi.Router, _ *ApiDefinition) chi.Router {
	return r.Route("/payment/v1/payments", func(r chi.Router) {
		//
	})
}

// CMS api V1 /api/cms/v1/translations
func cmsV1Routing(r chi.Router, a *ApiDefinition) chi.Router {
	return r.Route("/cms/v1/translations", func(r chi.Router) {
		h := cms.NewHandler(a.DB)
		r.Get("/", h.GetTranslations)
		r.Get("/{language}/{id}", h.GetTranslation)
	})
}

// Management api V1 /api/management/v1
func managementV1Routing(r chi.Router, a *ApiDefinition) chi.Router {
	return r.Route("/management/v1", func(r chi.Router) {
		h := management.NewHandler(a.DB)
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", h.GetCategories)
		})
		r.Post("/translations", h.CreateTranslations)
	})
}

// Product api V1 /api/product/v1/products
func productV1Routing(r chi.Router, a *ApiDefinition) chi.Router {
	return r.Route("/product/v1/products", func(r chi.Router) {
		h := product.NewApiHandler(product.NewProductService(a.DB))
		r.Get("/", h.GetProducts)
		// /api/product/v1/products/{productId}
		r.Route("/{productId}", func(r chi.Router) {
			r.Get("/", h.GetProduct)
		})
	})
}

func (a *ApiDefinition) Run(r *chi.Mux) {
	log.Printf("Running the server on port %s", a.SC.Addr)
	if err := http.ListenAndServe(a.SC.Addr, r); err != nil {
		log.Fatal(err)
	}

}
