package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"golang.org/x/sync/errgroup"
)

func (env HandlerEnv) RequestAccuralData(ctx context.Context) {
	log.Print("Getting numbers to check")
	numbersToCheck, err := env.Storage.GetNumbersToCheckInAccrual(ctx)
	if err != nil {
		log.Print("что-то пошло не так на этапе получения номеров для отправки" + err.Error())
	}

	g := new(errgroup.Group)
	inputCh := generator(numbersToCheck)

	for number := range inputCh {
		log.Print("sending number to accrual " + number)
		number := number
		g.Go(func() error {
			err := env.SingleAccrualOrderHandle(ctx, number)
			if err != nil {
				return err
			}
			return nil
		})
		time.Sleep(time.Second)
	}

	if err := g.Wait(); err != nil {
		log.Print(err)
	}
}

func (env HandlerEnv) SingleAccrualOrderHandle(ctx context.Context, num string) error {
	requestURL := fmt.Sprintf("http://localhost%s", env.ConfigStruct.FlagAccrualSystemAddress)
	requestPath := "/api/orders/"
	response, err := http.Get(requestURL + requestPath + num)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		// парсинг ответа
		var accrualData models.AccrualData
		var buf bytes.Buffer

		_, err = buf.ReadFrom(response.Body)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(buf.Bytes(), &accrualData); err != nil {
			logger.Log.Debug("could not unmarshal registration data")
			return err
		}

		err = env.Storage.PostOrderWithAccrualData(ctx, num, accrualData)
		if err != nil {
			return err
		}
		return nil
	}
	defer response.Body.Close()
	return nil
}

func generator(input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			inputCh <- data
		}
	}()
	return inputCh
}
