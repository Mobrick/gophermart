package main

import (
	"github.com/Mobrick/gophermart/internal/config"
	"github.com/Mobrick/gophermart/internal/database"
	"github.com/Mobrick/gophermart/internal/handler"
	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {	
	cfg := config.MakeConfig()

	env := &handler.HandlerEnv{
		ConfigStruct: cfg,
		Storage:      database.NewDB(cfg.FlagDBConnectionAddress),
	}

	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)

	r.Get(`/ping`, env.PingDBHandle)

	r.Post("POST /api/user/register", env.RegisterHandle)
}
