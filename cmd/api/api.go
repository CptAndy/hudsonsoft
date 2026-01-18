package main

import (
	"log"
	"net/http"
	"time"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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
			r.Route("/create" , func(r chi.Router) {
				r.Put("/",app.createEmployeeHandler)
			})
			
			r.Route("/{employeeID}", func(r chi.Router) {
				r.Get("/", app.getEmployeeHandler)
				r.Delete("/", app.deleteEmployeeHandler)
			})

		})

	})
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
