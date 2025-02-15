package repository

import (
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

// GetUserDataOrRegister с помощью sql функций ищет айди и хэш пароля пользователя
// Если пользователь не существует - создает его и возвращает его айди и переданный в функцию хэш.
func GetUserDataOrRegister(username string, providedPassHash string) (int, []byte, error) {
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
