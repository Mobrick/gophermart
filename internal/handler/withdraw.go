package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"github.com/Mobrick/gophermart/internal/userauth"
	"github.com/ShiraazMoollatjie/goluhn"
)

func (env HandlerEnv) WithdrawHandle(res http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	id, ok := userauth.CookieIsValid(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var withdrawData models.WithdrawInputData
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &withdrawData); err != nil {
		logger.Log.Debug("could not unmarshal withdraw data")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	number := withdrawData.Order

	if err := goluhn.Validate(number); err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	sum, err := withdrawData.Sum.Float64()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	enoughPoints, err := env.Storage.CheckIfEnoughPoints(ctx, id, sum)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if !enoughPoints {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}
	err = env.Storage.WithdrawPoints(ctx, number, id, sum)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
