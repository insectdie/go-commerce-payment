package routes

import (
	"codebase-service/config"
	"codebase-service/util/middleware"
	"log"
	"net/http"
	"strings"
	"time"

	payment "codebase-service/handlers/payments"
	product "codebase-service/handlers/products"

	"github.com/spf13/viper"
)

type Routes struct {
	Router  *http.ServeMux
	Product *product.Handler
	Payment *payment.Handler
}

func URLRewriter(baseURLPath string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, baseURLPath)

		next.ServeHTTP(w, r)
	}
}

func (r *Routes) SetupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL != "" && baseURL != "/" {
		r.Router.HandleFunc(baseURL+"/", URLRewriter(baseURL, r.Router))
	}
}

func (r *Routes) SetupRouter() {
	r.Router = http.NewServeMux()
	r.SetupBaseURL()
	r.productRoutes()
	r.paymentRoutes()
}

func (r *Routes) productRoutes() {
	r.Router.HandleFunc("GET /products/{id}", middleware.ApplyMiddleware(r.Product.GetProduct, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("GET /products", middleware.ApplyMiddleware(r.Product.GetProducts, middleware.EnabledCors, middleware.LoggerMiddleware()))

	r.Router.HandleFunc("POST /products", middleware.ApplyMiddleware(r.Product.CreateProduct, middleware.GetUserId, middleware.EnabledCors, middleware.LoggerMiddleware()))
	r.Router.HandleFunc("DELETE /products/{id}", middleware.ApplyMiddleware(r.Product.DeleteProduct, middleware.GetUserId, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) paymentRoutes() {
	r.Router.HandleFunc("POST /payments", middleware.ApplyMiddleware(r.Payment.CreatePayment, middleware.EnabledCors, middleware.LoggerMiddleware()))
}

func (r *Routes) Run(port string) {
	r.SetupRouter()

	log.Printf("[Running-Success] clients on localhost on port :%s", port)
	srv := &http.Server{
		Handler:      r.Router,
		Addr:         "localhost:" + port,
		WriteTimeout: config.WriteTimeout() * time.Second,
		ReadTimeout:  config.ReadTimeout() * time.Second,
	}

	log.Panic(srv.ListenAndServe())
}
