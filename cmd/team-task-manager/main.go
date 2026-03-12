package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"team-task-manager/internal/api/v1"
	"team-task-manager/internal/clients/mysql"
	"team-task-manager/internal/clients/redis"
	"team-task-manager/internal/logger"
	ratelimiter "team-task-manager/internal/redis_rate_limiter"
	secretprovider "team-task-manager/internal/secret_provider"
	"team-task-manager/internal/service"
	"team-task-manager/internal/supports"
)

// @title           Team Task Managere
// @version         1.0
// @description     Service for managing teams and tasks
// @termsOfService  http://swagger.io/terms/

// @contact.name   Maksim
// @contact.url    https://github.com/hhullen
// @contact.email  hhullen@gmail.com

// @license.name  Creative Commons Attribution-NonCommercial 4.0 International Public License
// @license.url   https://creativecommons.org/licenses/by-nc/4.0/deed.en

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type 'Bearer <token>' to authenticate

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	apiLog := logger.NewLogger(os.Stdout, "API")
	serviceLog := logger.NewLogger(os.Stdout, "SERVICE")

	secrets := secretprovider.NewSecretProvider()

	dbHost, err := secrets.ReadSecret("db_host")
	if err != nil {
		log.Fatal(err)
	}
	dbPort, err := secrets.ReadSecret("db_port")
	if err != nil {
		log.Fatal(err)
	}
	user, err := secrets.ReadSecret("db_app_user")
	if err != nil {
		log.Fatal(err)
	}
	password, err := secrets.ReadSecret("db_app_password")
	if err != nil {
		log.Fatal(err)
	}
	dbName, err := secrets.ReadSecret("db_name")
	if err != nil {
		log.Fatal(err)
	}

	if !supports.IsInContainer() {
		dbHost = "localhost"
	}

	dbConn, err := mysql.NewMySQLConn(ctx, dbHost, dbPort, user, password, dbName)
	if err != nil {
		log.Fatal(err)
	}

	dbClient := mysql.NewClient(ctx, dbConn)

	redisHost, err := secrets.ReadSecret("redis_host")
	if err != nil {
		log.Fatal(err)
	}
	redisPort, err := secrets.ReadSecret("redis_port")
	if err != nil {
		log.Fatal(err)
	}
	redisPassword, err := secrets.ReadSecret("redis_password")
	if err != nil {
		log.Fatal(err)
	}

	if !supports.IsInContainer() {
		redisHost = "localhost"
	}

	redisConn, err := redis.NewRedisConn(ctx, redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	cacheClient := redis.NewClient(redisConn)

	service := service.NewService(ctx, dbClient, serviceLog, cacheClient, secrets)

	limiter := ratelimiter.NewRedisRateLimiter(ctx, redisConn)

	apiService, err := api.NewAPI(ctx, service, service, secrets, limiter, apiLog)
	if err != nil {
		log.Fatal(err)
	}

	err = apiService.StartListening()
	if err != nil {
		log.Fatal(err)
	}
}
