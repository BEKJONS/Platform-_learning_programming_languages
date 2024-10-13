package main

import (
	"Auth_service/api"
	"Auth_service/api/handler"
	"Auth_service/pkg/config"
	"Auth_service/pkg/logs"
	"Auth_service/service"
	"Auth_service/storage/postgres"
	"Auth_service/storage/redis"
	"log"
)

func main() {
	cfg := config.Load()
	logger := logs.InitLogger()

	redisClient := redis.ConnectDB()

	db, err := postgres.ConnectPostgres(cfg)
	if err != nil {
		logger.Error("Error connecting to database", "error", err)
		log.Fatal(err)
	}

	defer db.Close()

	storage1 := postgres.NewAuthStorage(db)

	service1 := service.NewAuthService(storage1, logger)

	redis1 := redis.NewRedisStorage(redisClient, logger)

	handler1 := handler.NewHandler(service1, logger, redis1)

	router := api.NewRouter(handler1)
	err = router.Run(cfg.AUTH_PORT)

	if err != nil {
		logger.Error("Error starting server", "error", err)
		log.Fatal(err)
	}
}
