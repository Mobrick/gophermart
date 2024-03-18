package handler

import (
	"io"
	"net/http"

	"github.com/Mobrick/gophermart/internal/userauth"
	"github.com/ShiraazMoollatjie/goluhn"
)

func (env HandlerEnv) OrderPostHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if !userauth.CookieIsValid(req) {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	number, err := io.ReadAll(io.Reader(req.Body))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	if err := goluhn.Validate(string(number)); err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// проверка наличия в базе, если есть, кто занял, если текущий пользователь - 200, если нет - 409, если номера нет - 202
	err = env.Storage.PostOrderOrReturnStatus(ctx, string(number))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusAccepted)
}
