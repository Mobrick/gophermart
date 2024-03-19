package handler

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/Mobrick/gophermart/internal/userauth"
	"github.com/ShiraazMoollatjie/goluhn"
)

func (env HandlerEnv) OrderPostHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := userauth.CookieIsValid(req)
	if  !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	number, err := io.ReadAll(io.Reader(req.Body))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := goluhn.Validate(string(number)); err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// проверка наличия в базе, если есть, кто занял, если текущий пользователь - 200, если нет - 409, если номера нет - 202
	thisUser, err := env.Storage.CheckIfOrderExists(ctx, string(number), id)
	// если номер не найден
	if err == sql.ErrNoRows {
		// TODO: в горутину и отправка запроса к системе начисления баллов
		err := env.Storage.PostOrder(ctx, string(number), id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusAccepted)
		return	
	// если другая, не ожидаемая ошибка то 500
	} else if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// если уже существует номер у этого юзера
	if thisUser {
		res.WriteHeader(http.StatusOK)
		return
	} 
	// если этот номер занят другим юзером
	res.WriteHeader(http.StatusConflict)
}
