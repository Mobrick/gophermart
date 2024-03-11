package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"go.uber.org/zap"
)

func (env HandlerEnv) RegisterHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var registrationData models.SimpleAccountData
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &registrationData); err != nil {
		logger.Log.Debug("could not unmarshal registration data")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	storage := env.Storage
	loginAlreadyInUse, err := storage.AddNewAccount(ctx, registrationData)
	if err != nil {
		logger.Log.Debug("could not copmplete user registration", zap.String("Attempted login", string(registrationData.Login)))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if loginAlreadyInUse {
		logger.Log.Debug("login already in use", zap.String("Attempted login", string(registrationData.Login)))
		http.Error(res, err.Error(), http.StatusConflict)
		return
	}
	res.WriteHeader(http.StatusOK)
}
