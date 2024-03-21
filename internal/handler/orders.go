package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/userauth"
)

func (env HandlerEnv) OrdersHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := userauth.CookieIsValid(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}	

	orders, err := env.Storage.GetOrdersByUserId(ctx, id)
	if err != nil {
		logger.Log.Debug("could not get orders by user id")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	// TODO: если accural пустое не включать его в результат
	resp, err := json.Marshal(orders)
	if err != nil {
		logger.Log.Debug("could not marshal response")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(resp))
}