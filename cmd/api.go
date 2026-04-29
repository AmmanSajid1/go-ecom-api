package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	repo "github.com/AmmanSajid1/go-ecom-api/internal/adapters/postgresql/sqlc"
	"github.com/AmmanSajid1/go-ecom-api/internal/orders"
	"github.com/AmmanSajid1/go-ecom-api/internal/products"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // important for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	productsService := products.NewService(repo.New(app.db))
	productsHandler := products.NewHandler(productsService)
	r.Get("/products", productsHandler.ListProducts)
	r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		productsHandler.FindProductByID(w, r, id)
	})

	ordersService := orders.NewService(app.db)
	ordersHandler := orders.NewHandler(ordersService)
	r.Post("/orders", ordersHandler.PlaceOrder)
	r.Get("/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		ordersHandler.GetOrderByID(w, r, id)

	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
