package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func (env HandlerEnv) GetAccrualOrder(number string) (http.Response, error) {
	requestURL := fmt.Sprintf("http://localhost%s", env.ConfigStruct.FlagAccrualSystemAddress)
	requestPath := "/api/orders/"

	response, err := http.Get(requestURL + requestPath + number)
	if err != nil {
		return *response, err
	}
	return *response, nil
}

func (env HandlerEnv) RequestAccuralData(ctx context.Context) {

	numbersToCheck, err := env.Storage.GetNumbersToCheckInAccrual(ctx)
	if err != nil {
		log.Print(err)
	}

	g := new(errgroup.Group)
	inputCh := generator(numbersToCheck)

	for number := range inputCh {
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
	resp, err := env.GetAccrualOrder(num)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
