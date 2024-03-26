package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/userauth"
)

func (env HandlerEnv) WithdrawalsHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := userauth.CookieIsValid(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}	

	orders, err := env.Storage.GetWithdrawals(ctx, id)
	if err != nil {
		logger.Log.Info("could not get withdrawals")
		log.Print("could not get withdrawals " + err.Error())
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
		logger.Log.Info("could not marshal response")
		log.Print("could not marshal response " + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(resp))
}