package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"github.com/Mobrick/gophermart/internal/userauth"
)

func (env HandlerEnv) BalanceHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := userauth.CookieIsValid(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	accrual, withdrawn, err := env.Storage.GetBalanceByUserID(ctx, id)
	if err != nil {
		logger.Log.Debug("could not get orders by user id")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	balance := models.BalanceData{
		Current:   json.Number(strconv.FormatFloat(accrual, 'e', -1, 64)),
		Withdrawn: json.Number(strconv.FormatFloat(withdrawn, 'e', -1, 64)),
	}

	resp, err := json.Marshal(balance)
	if err != nil {
		logger.Log.Debug("could not marshal response")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(resp))
}
