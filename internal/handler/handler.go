package handler

import (
	"github.com/Mobrick/gophermart/internal/config"
	"github.com/Mobrick/gophermart/internal/database"
)

type HandlerEnv struct {
	ConfigStruct *config.Config
	Storage      database.Storage
}

