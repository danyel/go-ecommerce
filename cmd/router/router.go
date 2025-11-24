package router

import (
	"log"
	"net/http"

	"github.com/dnoulet/ecommerce/cmd/config"
	"github.com/dnoulet/ecommerce/internal/category"
	"github.com/dnoulet/ecommerce/internal/cms"
	"github.com/dnoulet/ecommerce/internal/management"
	"github.com/dnoulet/ecommerce/internal/product"
	"github.com/dnoulet/ecommerce/internal/product-management"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type ApiDefinition struct {
	ServerConfiguration *config.ServerConfiguration
	DB                  *gorm.DB
}

func (a *ApiDefinition) ConfigRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// the routing should start with /api/<capability>/<version>/<resources>/... (GET, POST, DELETE, PUT, ...)
	r.Route("/api", func(r chi.Router) {
		// Product api V1 /api/product/v1/products
		r.Route("/product/v1/products", func(r chi.Router) {
			handler := product.NewHandler(a.DB)
			r.Get("/", handler.GetProducts)
			// /api/product/v1/products/{productId}
			r.Route("/{productId}", func(r chi.Router) {
				r.Get("/", handler.GetProduct)
			})
		})
		// Category api V1 /api/category/v1/categories
		r.Route("/category/v1", func(r chi.Router) {
			handler := category.NewHandler(a.DB)
			r.Route("/categories", func(r chi.Router) {
				r.Post("/", handler.CreateCategory)
			})
			r.Post("/translations", handler.CreateTranslations)
		})

		// Product Management api V1 /api/product-management/v1
		r.Route("/product-management/v1", func(r chi.Router) {
			r.Route("/products", func(r chi.Router) {
				handler := product_management.NewHandler(a.DB)
				r.Get("/", handler.GetProducts)
				r.Post("/", handler.CreateProduct)
				r.Route("/{productId}", func(r chi.Router) {
					r.Get("/", handler.GetProduct)
					r.Delete("/", handler.DeleteProduct)
					r.Put("/", handler.UpdateProduct)
				})
			})
		})

		// Management api V1 /api/management/v1
		r.Route("/management/v1", func(r chi.Router) {
			handler := management.NewHandler(a.DB)
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", handler.GetCategories)
			})
			r.Post("/translations", handler.CreateTranslations)
		})
		// CMS api V1 /api/cms/v1/translations
		r.Route("/cms/v1/translations", func(r chi.Router) {
			handler := cms.NewHandler(a.DB)
			r.Get("/{language}/{id}", handler.GetTranslation)
		})
		// Payment api V1 /api/payment/v1/payments
		r.Route("/payment/v1/payments", func(r chi.Router) {})
		// Payment api V1 /api/payment/v1/payments
		r.Route("/order/v1/orders", func(r chi.Router) {})
		// Shopping Basket api V1 /api/shopping-basket/v1/shopping-baskets
		r.Route("/shopping-basket/v1/shopping-baskets", func(r chi.Router) {})

	})
	return r
}

func (a *ApiDefinition) Run(r *chi.Mux) {
	log.Printf("Running the server on port %s", a.ServerConfiguration.Addr)
	if err := http.ListenAndServe(a.ServerConfiguration.Addr, r); err != nil {
		log.Fatal(err)
	}

}
