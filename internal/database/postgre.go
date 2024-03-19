package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log"

	"github.com/Mobrick/gophermart/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
)

const accountsTableName = "accounts"

var embedMigrations embed.FS

type PostgreDB struct {
	DatabaseConnection *sql.DB
	DatabaseMap        map[string]string
}

func (dbData PostgreDB) PingDB() error {
	err := dbData.DatabaseConnection.Ping()
	return err
}

// Возвращает true если такой логин уже хранится в базе
func (dbData PostgreDB) AddNewAccount(ctx context.Context, accountData models.SimpleAccountData) (bool, string, error) {
	id := uuid.New().String()

	err := dbData.createAccountsTableIfNotExists(ctx)
	if err != nil {
		return false,"", nil
	}

	insertStmt := "INSERT INTO " + accountsTableName + " (uuid, login, password)" +
		" VALUES ($1, $2, $3)"

	_, err = dbData.DatabaseConnection.ExecContext(ctx, insertStmt, id, accountData.Login, accountData.Password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("login %s already in database", accountData.Login)
			return true,"", nil
		} else {
			log.Printf("Failed to insert a record: " + accountData.Login)
			return false,"", err
		}
	}

	return false, id, nil
}

func (dbData PostgreDB) CheckLogin(ctx context.Context, accountData models.SimpleAccountData) (string, error) {

	checkStmt := "SELECT uuid FROM accounts WHERE login=$1 AND password=$2)"

	var id string

	err := dbData.DatabaseConnection.QueryRowContext(ctx, checkStmt, accountData.Login, accountData.Password).Scan(&id)

	if err != nil {
		log.Printf("Error querying database: " + accountData.Login)
		return "", err
	}

	return id, nil
}

func (dbData PostgreDB) createAccountsTableIfNotExists(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	// добавить embedding
	if err := goose.Up(dbData.DatabaseConnection, "accounts_migrations"); err != nil {
		return err
	}

	return nil
}

func (dbData PostgreDB) Close() {
	dbData.DatabaseConnection.Close()
}

func (dbData PostgreDB) CheckIfOrderExists(ctx context.Context, number string, currentUserUUID string) (bool, error) {
	err := dbData.createOrdersTableIfNotExists(ctx)
	if err != nil {
		return false, err
	}
	var uuid string
	// ищем существует ли, если да то кто владелец заказа
	row := dbData.DatabaseConnection.QueryRowContext(ctx, "SELECT account_uuid FROM orders WHERE number = $1", number)
	err = row.Scan(&uuid)
	if err == nil {
		if uuid == currentUserUUID{
			return true, nil
		}
		return false, nil
	}
	return false, err
}

	// если не существует, добавляем в таблицу горутиной
	// реализация без горутины
func (dbData PostgreDB) PostOrder(ctx context.Context, number string, currentUserUUID string) error {
	// отправка в систему начисления баллов для проверки запроса
	// формирование запроса
	// парсинг ответа
	// в любом случае создаем запись в даблице заказов
	_, err := dbData.DatabaseConnection.ExecContext(ctx, "INSERT INTO url_records (number, account_uuid) VALUES ($1, $2)", number, currentUserUUID)
	if err != nil {
		return err
	}
	return nil
}

func (dbData PostgreDB) createOrdersTableIfNotExists(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	// добавить embedding
	if err := goose.Up(dbData.DatabaseConnection, "orders_migrations"); err != nil {
		return err
	}

	return nil
}