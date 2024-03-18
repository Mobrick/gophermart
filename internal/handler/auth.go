package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"github.com/Mobrick/gophermart/internal/userauth"
)

func (env HandlerEnv) AuthHandle(res http.ResponseWriter, req *http.Request) {
	if userauth.CookieIsValid(req) {
		res.WriteHeader(http.StatusOK)
		return
	}

	ctx := req.Context()

	var loginData models.SimpleAccountData
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &loginData); err != nil {
		logger.Log.Debug("could not unmarshal registration data")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := env.Storage.CheckLogin(ctx, loginData)
	if err != nil {
		logger.Log.Debug("could not check login data")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(id) == 0 {
		res.WriteHeader(http.StatusUnauthorized)
	}

	cookie, err := userauth.CreateNewCookie(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(res, &cookie)

	res.WriteHeader(http.StatusOK)
}
