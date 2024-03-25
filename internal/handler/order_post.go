package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"github.com/Mobrick/gophermart/internal/userauth"
	"github.com/ShiraazMoollatjie/goluhn"
)

func (env HandlerEnv) OrderPostHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, ok := userauth.CookieIsValid(req)
	if !ok {
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
		shouldReturn := env.postOrder(ctx, number, res, id)
		if shouldReturn {
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

func (env HandlerEnv) postOrder(ctx context.Context, number []byte, res http.ResponseWriter, id string) bool {
	// TODO: в горутину и отправка запроса к системе начисления баллов
	// отправка в систему начисления баллов для проверки запроса
	// формирование запроса
	response, err := env.GetAccrualOrder(string(number))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return true
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		// парсинг ответа
		var accrualData models.AccrualData
		var buf bytes.Buffer

		_, err = buf.ReadFrom(response.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return true
		}

		if err = json.Unmarshal(buf.Bytes(), &accrualData); err != nil {
			logger.Log.Debug("could not unmarshal registration data")
			http.Error(res, err.Error(), http.StatusBadRequest)
			return true
		}

		err = env.Storage.PostOrderWithAccrualData(ctx, string(number), id, accrualData)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return true
		}
		return false
	}

	// в любом случае создаем запись в даблице заказов
	err = env.Storage.PostOrder(ctx, string(number), id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
