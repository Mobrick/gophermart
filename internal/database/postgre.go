package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log"
	"time"

	accountsmigrations "github.com/Mobrick/gophermart/internal/database/accounts_migrations"
	ordersmigrations "github.com/Mobrick/gophermart/internal/database/orders_migrations"
	"github.com/Mobrick/gophermart/internal/logger"
	"github.com/Mobrick/gophermart/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
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

var embedMigrations embed.FS

// Возвращает true если такой логин уже хранится в базе
func (dbData PostgreDB) AddNewAccount(ctx context.Context, accountData models.SimpleAccountData) (bool, string, error) {

	err := dbData.createAccountsTable(ctx)
	if err != nil {
		return false, "", err
	}

	id := uuid.New().String()

	insertStmt := "INSERT INTO " + accountsTableName + " (uuid, username, password)" +
		" VALUES ($1, $2, $3)"

	_, err = dbData.DatabaseConnection.ExecContext(ctx, insertStmt, id, accountData.Login, accountData.Password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("login %s already in database", accountData.Login)
			return true, "", nil
		} else {
			log.Printf("Failed to insert a record: " + accountData.Login)
			return false, "", err
		}
	}

	return false, id, nil
}

func (dbData PostgreDB) createAccountsTable(ctx context.Context) error {
	provider, err := goose.NewProvider(database.DialectPostgres, dbData.DatabaseConnection, accountsmigrations.EmbedAccounts)
	if err != nil {
		return err
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return err
	}

	for _, r := range results {
		log.Printf("%-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}

	logger.Log.Debug("Created table with goose embed")
	return nil
}

func (dbData PostgreDB) CheckLogin(ctx context.Context, accountData models.SimpleAccountData) (string, error) {

	checkStmt := "SELECT uuid FROM accounts WHERE username=$1 AND password=$2"

	var id string

	err := dbData.DatabaseConnection.QueryRowContext(ctx, checkStmt, accountData.Login, accountData.Password).Scan(&id)

	if err != nil {
		log.Printf("Error querying database: " + accountData.Login)
		return "", err
	}

	return id, nil
}

func (dbData PostgreDB) Close() {
	dbData.DatabaseConnection.Close()
}

func (dbData PostgreDB) CheckIfOrderExists(ctx context.Context, number string, currentUserUUID string) (bool, error) {
	var uuid string

	err := dbData.createOrdersTable(ctx)
	if err != nil {
		return false, err
	}
	// ищем существует ли, если да то кто владелец заказа
	row := dbData.DatabaseConnection.QueryRowContext(ctx, "SELECT account_uuid FROM orders WHERE number = $1", number)
	err = row.Scan(&uuid)
	if err == nil {
		if uuid == currentUserUUID {
			return true, nil
		}
		return false, nil
	}
	return false, err
}

// если не существует, добавляем в таблицу горутиной
// реализация без горутины
func (dbData PostgreDB) PostOrder(ctx context.Context, number string, currentUserUUID string) error {
	_, err := dbData.DatabaseConnection.ExecContext(ctx, "INSERT INTO orders (number, account_uuid, status) VALUES ($1, $2, $3)", number, currentUserUUID, "")
	if err != nil {
		return err
	}
	return nil
}

func (dbData PostgreDB) createOrdersTable(ctx context.Context) error {
	provider, err := goose.NewProvider(database.DialectPostgres, dbData.DatabaseConnection, ordersmigrations.EmbedOrders)
	if err != nil {
		return err
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return err
	}

	for _, r := range results {
		log.Printf("%-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}

	logger.Log.Debug("Created table with goose embed")
	return nil
}

func (dbData PostgreDB) PostOrderWithAccrualData(ctx context.Context, number string, accrualData models.AccrualData) error {
	stmt := "UPDATE orders SET proceeded_at = $1, status = $2, accrual = $3 WHERE number = $4"
	_, err := dbData.DatabaseConnection.ExecContext(ctx, stmt, time.Now(), accrualData.Status, accrualData.Accrual, number)
	if err != nil {
		return err
	}
	return nil
}

func (dbData PostgreDB) GetOrdersByUserID(ctx context.Context, id string) ([]models.OrderData, error) {
	var ordersData []models.OrderData
	stmt := "SELECT number, status, accrual, uploaded_at FROM orders WHERE account_uuid = $1"
	rows, err := dbData.DatabaseConnection.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var number, status, uploadedAt string
		var accrual int
		err := rows.Scan(&number, &status, &accrual, &uploadedAt)
		if err != nil {
			return nil, err
		}

		order := models.OrderData{
			Number:     number,
			Status:     status,
			Accrual:    accrual,
			UploadedAt: uploadedAt,
		}
		ordersData = append(ordersData, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return ordersData, nil
}

func (dbData PostgreDB) GetBalanceByUserID(ctx context.Context, id string) (int, int, error) {
	var accural, withdrawn int
	stmt := "SELECT accrual FROM orders WHERE account_uuid = $1"
	rows, err := dbData.DatabaseConnection.QueryContext(ctx, stmt, id)
	if err != nil {
		return 0, 0, err
	}
	for rows.Next() {
		var value int
		err := rows.Scan(&value)
		if err != nil {
			return 0, 0, err
		}
		if value > 0 {
			accural += value
		} else if value < 0 {
			withdrawn += value
		}
	}

	if err = rows.Err(); err != nil {
		return 0, 0, err
	}
	accural += withdrawn
	withdrawn *= -1

	defer rows.Close()
	return accural, withdrawn, nil
}

func (dbData PostgreDB) WithdrawPoints(ctx context.Context, number string, id string, amount int) error {
	// отправка в систему начисления баллов для проверки запроса
	// формирование запроса
	// парсинг ответа
	// в любом случае создаем запись в даблице заказов
	_, err := dbData.DatabaseConnection.ExecContext(ctx, "INSERT INTO orders (number, account_uuid, accrual) VALUES ($1, $2, $3)", number, id, amount)
	if err != nil {
		return err
	}
	return nil
}

func (dbData PostgreDB) CheckIfEnoughPoints(ctx context.Context, id string, amount int) (bool, error) {
	accural, _, err := dbData.GetBalanceByUserID(ctx, id)
	if err != nil {
		return false, err
	}

	if amount > accural {
		return false, nil
	}
	return true, nil
}

func (dbData PostgreDB) GetWithdrawals(ctx context.Context, id string) ([]models.WithdrawData, error) {
	var ordersData []models.WithdrawData
	stmt := "SELECT number, accrual, proceeded_at FROM orders WHERE account_uuid = $1 AND accural < 0"
	rows, err := dbData.DatabaseConnection.QueryContext(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var number string
		var accrual int
		var proceededAt time.Time
		err := rows.Scan(&number, &accrual, &proceededAt)
		if err != nil {
			return nil, err
		}

		order := models.WithdrawData{
			Order:       number,
			Sum:         accrual,
			ProceededAt: proceededAt,
		}
		ordersData = append(ordersData, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return ordersData, nil
}

func (dbData PostgreDB) GetNumbersToCheckInAccrual(ctx context.Context) ([]string, error) {
	var numbers []string
	stmt := "SELECT number FROM orders WHERE status NOT IN ('PROCESSED','INVALID')"
	rows, err := dbData.DatabaseConnection.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var number string
		err := rows.Scan(&number)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, number)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return numbers, nil
}
