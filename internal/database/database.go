package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/Mobrick/gophermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage interface {
	PingDB() error
	AddNewAccount(context.Context, models.SimpleAccountData) (bool, string, error)
	CheckLogin(context.Context, models.SimpleAccountData) (string, error)
	CheckIfOrderExists(context.Context, string, string) (bool, error)
	PostOrder(context.Context, string, string) error
	PostOrderWithAccrualData(context.Context, string, string, models.AccrualData) error
	GetOrdersByUserId(context.Context, string) ([]models.OrderData, error)
	GetBalanceByUserId(context.Context, string) (int, int, error)
	WithdrawPoints(context.Context, string, string, int) (error)
	CheckIfEnoughPoints(context.Context, string, int) (bool, error)
	GetWithdrawals(context.Context, string) ([]models.WithdrawData, error)
	GetNumbersToCheckInAccrual(context.Context) ([]string, error)
	Close()
}

func NewDB(connectionString string) Storage {
	dbData := PostgreDB{
		DatabaseMap:        make(map[string]string),
		DatabaseConnection: NewDBConnection(connectionString),
	}

	return dbData
}

func NewDBConnection(connectionString string) *sql.DB {
	// Закрывается в основном потоке
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}
