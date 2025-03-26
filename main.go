package main

import (
	"codebase-service/config"
	paymentHandler "codebase-service/handlers/payments"
	productHandler "codebase-service/handlers/products"
	circuitbreaker "codebase-service/infra/circuit_breaker"
	midtranssvc "codebase-service/integration/midtrans"
	"codebase-service/repository/products"
	"codebase-service/routes"
	productSvc "codebase-service/usecases/products"
	"context"
	"database/sql"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	dbConn, err := config.ConnectToDatabase(config.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		return
	}
	defer dbConn.Close()

	redisConn, err := config.ConnectToRedis(config.RedisConnection{
		Host: cfg.RedisHost,
		Port: cfg.RedisPort,
		Pass: cfg.RedisPass,
		DB:   cfg.RedisDB,
	})
	if err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
		return
	}

	// checj if redis is connected
	err = redisConn.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
		return
	} else {
		log.Println("connected to redis")
	}

	validator := validator.New()
	cb := circuitbreaker.NewCircuitBreakerInstance()

	routes := setupRoutes(cfg, dbConn, redisConn, validator, cb)
	routes.Run(cfg.AppPort)
}

func setupRoutes(
	cfg *config.Config,
	db *sql.DB,
	rdb *redis.Client,
	validator *validator.Validate,
	cb *gobreaker.CircuitBreaker,
) *routes.Routes {

	productStore := products.NewStore(db, rdb)
	productSvc := productSvc.NewProductSvc(productStore)
	productHandler := productHandler.NewHandler(productSvc, validator)

	midtranssvc := midtranssvc.NewMidtransContract(cfg, cb)
	paymentHandler := paymentHandler.NewHandler(midtranssvc, validator)

	return &routes.Routes{
		Product: productHandler,
		Payment: paymentHandler,
	}
}
