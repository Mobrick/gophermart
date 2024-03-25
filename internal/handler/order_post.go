package handler

import (
	"context"
	"database/sql"
	"io"
	"net/http"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/userauth"
	"github.com/ShiraazMoollatjie/goluhn"
	"golang.org/x/sync/errgroup"
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
		g := new(errgroup.Group)
		g.Go(func() error {
			err := env.postOrder(ctx, number, res, id)
			if err != nil {
				return err
			}

			return nil
		})

		if err := g.Wait(); err != nil {
			logger.Log.Debug("could not post order")
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

func (env HandlerEnv) postOrder(ctx context.Context, number []byte, res http.ResponseWriter, id string) error {
	// в любом случае создаем запись в даблице заказов
	err := env.Storage.PostOrder(ctx, string(number), id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
