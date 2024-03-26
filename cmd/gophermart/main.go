package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mobrick/gophermart/internal/config"
	"github.com/Mobrick/gophermart/internal/database"
	"github.com/Mobrick/gophermart/internal/handler"
	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	sugar := *zapLogger.Sugar()
	logger.Sugar = sugar

	cfg := config.MakeConfig()

	env := &handler.HandlerEnv{
		ConfigStruct: cfg,
		Storage:      database.NewDB(cfg.FlagDBConnectionAddress),
	}

	defer env.Storage.Close()
	ctx := context.Background()
	go func() {
		for {
			log.Print("Requests to accrual " + time.Now().GoString())
			time.Sleep(time.Second * 10)
			env.RequestAccuralData(ctx)
		}
	}()

	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)

	r.Get(`/ping`, env.PingDBHandle)
	r.Get(`/api/user/orders`, env.OrdersHandle)
	r.Get(`/api/user/balance`, env.BalanceHandle)
	r.Get(`/api/user/withdrawals`, env.WithdrawalsHandle)

	r.Post("/api/user/register", env.RegisterHandle)
	r.Post("/api/user/login", env.AuthHandle)
	r.Post("/api/user/orders", env.OrderPostHandle)
	r.Post("/api/user/balance/withdraw", env.WithdrawHandle)

	sugar.Infow(
		"Starting server",
		"addr", cfg.FlagRunAddr,
	)

	server := &http.Server{
		Addr:    cfg.FlagRunAddr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	env.Storage.Close()
	sugar.Infow("Server stopped")
}
