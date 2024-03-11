package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Mobrick/gophermart/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const accountsTableName = "accounts"

type PostgreDB struct {
	DatabaseConnection *sql.DB
	DatabaseMap        map[string]string
}

func (dbData PostgreDB) PingDB() error {
	err := dbData.DatabaseConnection.Ping()
	return err
}

// Возвращает true если такой логин уже хранится в базе
func (dbData PostgreDB) AddNewAccount(ctx context.Context, accountData models.SimpleAccountData) (bool, error) {
	id := uuid.New().String()

	err := dbData.createAccountsTableIfNotExists(ctx)
	if err != nil {
		return false, nil
	}

	insertStmt := "INSERT INTO " + accountsTableName + " (uuid, login, password)" +
		" VALUES ($1, $2, $3)"

	_, err = dbData.DatabaseConnection.ExecContext(ctx, insertStmt, id, accountData.Login, accountData.Password)

	
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("login %s already in database", accountData.Login)			
			return true, nil
		} else {
			log.Printf("Failed to insert a record: " + accountData.Login)
			return false, err
		}
	}

	return false, nil
}

func (dbData PostgreDB) createAccountsTableIfNotExists(ctx context.Context) error {
	_, err := dbData.DatabaseConnection.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS "+accountsTableName+
			` (uuid TEXT PRIMARY KEY, 
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL)`)

	if err != nil {
		return err
	}
	return nil
}
