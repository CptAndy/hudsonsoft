package main

import (
	"log"
	"net/http"
	"time"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/employee", func(r chi.Router) {
			r.Route("/create", func(r chi.Router) {
				r.Put("/", app.createEmployeeHandler)
			})

			r.Route("/{employeeID}", func(r chi.Router) {
				r.Get("/", app.getEmployeeHandler)
				r.Delete("/", app.deleteEmployeeHandler)
			})

		})

		r.Route("/customer", func(r chi.Router) {
			r.Route("/create", func(r chi.Router) {
				r.Put("/", app.createCustomerHandler)
			})

			r.Route("/{customerID}", func(r chi.Router) {
				r.Get("/", app.getCustomerHandler)
				r.Delete("/", app.deleteCustomerHandler)
			})

		})

		r.Route("/inventory", func(r chi.Router) {
			r.Route("/product/addproduct", func(r chi.Router) {
				// adding product to inventory
				r.Put("/", app.createProductHandler)
			})

			r.Route("/product/{productID}", func(r chi.Router) {
				r.Get("/", app.getProductHandler)
				r.Delete("/", app.deleteProductHandler)
			})

			r.Route("/addvariation", func(r chi.Router) {
				r.Put("/", app.createProductTypeHandler)
			})

			r.Route("/variation/{prodTypeID}", func(r chi.Router) {
				r.Get("/", app.getProductTypeHandler)
				r.Delete("/", app.deleteProductTypeHandler)

			})


		})

		r.Route("/adminfunctions", func(r chi.Router) {
				r.Put("/", app.createReturnReasonHandler)
			})

	}) // END OF /v1/
	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("SERVER RUNNING AT %s", app.config.addr)
	log.Printf("ENV: %s", app.config.env)

	return srv.ListenAndServe()
}
