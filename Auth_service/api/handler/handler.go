package handler

import (
	"Auth_service/service"
	"Auth_service/storage/redis"
	"log/slog"
)

type Handler struct {
	log   *slog.Logger
	auth  service.AuthService
	redis *redis.RedisStorage
}

func NewHandler(auth service.AuthService, log *slog.Logger, redis *redis.RedisStorage) *Handler {
	return &Handler{log, auth, redis}
}
