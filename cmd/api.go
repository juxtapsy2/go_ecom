package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/juxtapsy2/go_ecom/internal/adapters/postgresql/sqlc"
	"github.com/juxtapsy2/go_ecom/internal/orders"
	"github.com/juxtapsy2/go_ecom/internal/products"
)

type application struct {
	config config
	db     *pgx.Conn
}

// mount
func (app *application) mount() http.Handler {
	// gorilla
	// chi
	// fi
	router := chi.NewRouter()
	// A good base middleware stack
	router.Use(middleware.RequestID) // for rate limiting
	router.Use(middleware.RealIP)    // for rate limit and analytics tracing
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out
	// and further processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	router.Get("/products", productHandler.ListProducts)

	orderService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderService)
	router.Post("/orders", ordersHandler.PlaceOrder)

	return router
}

// run
func (app *application) run(h http.Handler) error {
	service := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at addr %s", app.config.addr)
	return service.ListenAndServe()
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dbConnectionString string
}
