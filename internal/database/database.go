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
	PostOrderOrReturnStatus(context.Context, string) (error)
	AddNewAccount(context.Context, models.SimpleAccountData) (bool, string, error)
	CheckLogin(context.Context, models.SimpleAccountData) (string, error)
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
