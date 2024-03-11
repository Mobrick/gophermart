package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/Mobrick/gophermart/internal/models"
)

type Storage interface {
	PingDB() error
	AddNewAccount(context.Context, models.SimpleAccountData) (bool, error)
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
