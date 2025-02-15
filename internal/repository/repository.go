package repository

import (
	"avito_internship/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)
import _ "github.com/jackc/pgx/v5/stdlib"
import "avito_internship/internal/config"

var db *sql.DB

// Connect устанавливает соединение с базой данных и сохраняет его в переменной db.
func Connect() {
	cfg := config.Get()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
	conn, err := sql.Open("pgx/v5", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	db = conn
	log.Println("Успешное подключение к базе данных")
}

// BuyItemsForUser осуществляет покупку определенного количества вещей
func BuyItemsForUser(userID int, itemName string, amount int) error {
	_, err := db.Query("SELECT buy_item($1, $2, $3);", userID, itemName, amount)
	return err
}

// SendCoins осуществляет перевод коинов от одного пользователя к другому
func SendCoins(userFromID, amount int, userTo string) error {
	_, err := db.Query("SELECT transfer_coins($1, $2, $3);", userFromID, userTo, amount)
	return err
}

// GetUserIDPassHashOrRegister с помощью sql функций ищет айди и хэш пароля пользователя
// Если пользователь не существует - создает его и возвращает его айди и переданный в функцию хэш.
func GetUserIDPassHashOrRegister(username string, providedPassHash string) (int, []byte, error) {
	tx, err := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return 0, nil, err
	}
	var userId int
	var userPassHash string

	err = tx.QueryRow(
		"SELECT id, password_hash FROM get_user_id_password_hash($1);",
		username,
	).Scan(&userId, &userPassHash)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRow("SELECT register_user($1, $2);", username, providedPassHash).Scan(&userId)
		userPassHash = providedPassHash
	}

	if err != nil {
		return 0, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, nil, err
	}
	return userId, []byte(userPassHash), nil
}

// GetUserBalanceInventoryLogs с помощью sql функций получает баланс пользователя, инвентари и последние транзакции
func GetUserBalanceInventoryLogs(userID int) (models.InfoResponse, error) {
	tx, err := db.Begin()
	if err != nil {
		return models.InfoResponse{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var result models.InfoResponse

	err = tx.QueryRow("SELECT get_user_balance($1);", userID).Scan(&result.Coins)
	if err != nil {
		return models.InfoResponse{}, err
	}

	rows, err := tx.Query("SELECT * FROM get_user_inventory($1);", userID)
	if err != nil {
		return models.InfoResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return models.InfoResponse{}, err
		}
		result.Inventory = append(result.Inventory, item)
	}
	if err = rows.Err(); err != nil {
		return models.InfoResponse{}, err
	}

	rows, err = tx.Query("SELECT * FROM get_user_receive_history($1);", userID)
	if err != nil {
		return models.InfoResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var received models.CoinTransaction
		if err := rows.Scan(&received.User, &received.Amount); err != nil {
			return models.InfoResponse{}, err
		}
		result.CoinHistory.Received = append(result.CoinHistory.Received, received)
	}
	if err = rows.Err(); err != nil {
		return models.InfoResponse{}, err
	}

	rows, err = tx.Query("SELECT * FROM get_user_send_history($1);", userID)
	if err != nil {
		return models.InfoResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var sent models.CoinTransaction
		if err := rows.Scan(&sent.User, &sent.Amount); err != nil {
			return models.InfoResponse{}, err
		}
		result.CoinHistory.Sent = append(result.CoinHistory.Sent, sent)
	}
	if err = rows.Err(); err != nil {
		return models.InfoResponse{}, err
	}

	if err = tx.Commit(); err != nil {
		return models.InfoResponse{}, err
	}

	return result, nil
}
