package router

import (
	"log"
	"net/http"

	"github.com/dnoulet/ecommerce/cmd/config"
	"github.com/dnoulet/ecommerce/internal/cms"
	"github.com/dnoulet/ecommerce/internal/management"
	"github.com/dnoulet/ecommerce/internal/product"
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
			productHandler := product.NewProductHandler(a.DB)
			r.Get("/", productHandler.GetProducts)
			r.Post("/", productHandler.CreateProduct)
			// /api/product/v1/products/{productId}
			r.Route("/{productId}", func(r chi.Router) {
				r.Get("/", productHandler.GetProduct)
				r.Delete("/", productHandler.DeleteProduct)
				r.Put("/", productHandler.UpdateProduct)
			})
		})
		r.Route("/management/v1", func(r chi.Router) {
			managementHandler := management.NewManagementHandler(a.DB)
			r.Route("/categories", func(r chi.Router) {
				r.Post("/", managementHandler.CreateCategory)
				r.Get("/", managementHandler.GetCategories)
			})
			r.Post("/translations", managementHandler.CreateTranslations)
		})
		r.Route("/cms/v1/translations", func(r chi.Router) {
			cmsHandler := cms.NewCmsHandler(a.DB)
			r.Get("/{language}/{id}", cmsHandler.GetTranslation)
		})
		r.Route("/payment/v1/payments", func(r chi.Router) {})
		r.Route("/order/v1/orders", func(r chi.Router) {})
		r.Route("/shopping-basket/v1/shopping-baskets", func(r chi.Router) {})

	})
	return r
}

func (a *ApiDefinition) Run(r *chi.Mux) {
	log.Printf("Running the server on port %s", a.ServerConfiguration.Addr)
	http.ListenAndServe(a.ServerConfiguration.Addr, r)
}
