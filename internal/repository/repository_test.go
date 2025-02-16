package repository

import (
	"avito_internship/internal/models"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var mockDB *sql.DB
var mock sqlmock.Sqlmock

// -------------------------------------------
// Инициализация тестов и мокнутой базы данных
// -------------------------------------------
func initMock() {
	var err error
	mockDB, mock, err = sqlmock.New()
	if err != nil {
		panic("Ошибка при создании мока БД: " + err.Error())
	}
	db = mockDB
}

func teardown() {
	mockDB.Close()
}

func TestMain(m *testing.M) {
	initMock()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// ---------------------------------
// Тесты GetUserIDPassHashOrRegister
// ---------------------------------
func TestGetUserIDPassHashOrRegisterValidLogin(t *testing.T) {
	resetMockDB(t)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, password_hash FROM get_user_id_password_hash\\(\\$1\\);").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).AddRow(1, "userPassHash"))
	mock.ExpectCommit()
	userID, passHash, err := GetUserIDPassHashOrRegister("test", "passHash")
	assert.NoError(t, err)
	assert.Equal(t, userID, 1)
	assert.Equal(t, string(passHash), "userPassHash")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserIDPassHashOrRegisterLoginWithError(t *testing.T) {
	resetMockDB(t)
	returningError := fmt.Errorf("error")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, password_hash FROM get_user_id_password_hash\\(\\$1\\);").
		WithArgs("test").
		WillReturnError(returningError)
	mock.ExpectRollback()
	userID, passHash, err := GetUserIDPassHashOrRegister("test", "passHash")
	assert.ErrorIs(t, err, returningError)
	assert.Equal(t, userID, 0)
	assert.Equal(t, string(passHash), "")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserIDPassHashOrRegisterValidRegistration(t *testing.T) {
	resetMockDB(t)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, password_hash FROM get_user_id_password_hash\\(\\$1\\);").
		WithArgs("test").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT register_user\\(\\$1, \\$2\\);").
		WithArgs("test", "passHash").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	userID, passHash, err := GetUserIDPassHashOrRegister("test", "passHash")
	assert.NoError(t, err)
	assert.Equal(t, userID, 1)
	assert.Equal(t, string(passHash), "passHash")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---------------------------------
// Тесты GetUserBalanceInventoryLogs
// ---------------------------------
func TestGetUserBalanceInventoryLogsValid(t *testing.T) {
	resetMockDB(t)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT get_user_balance\\(\\$1\\);").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(250))
	mock.ExpectQuery("SELECT \\* FROM get_user_inventory\\(\\$1\\);").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"type", "quantity"}).
			AddRow("t_shirt", 1).
			AddRow("cup", 2).
			AddRow("book", 1).
			AddRow("powerbank", 1),
		)
	mock.ExpectQuery("SELECT \\* FROM get_user_receive_history\\(\\$1\\);").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user", "amount"}).
			AddRow("user1", 100).
			AddRow("user2", 50),
		)
	mock.ExpectQuery("SELECT \\* FROM get_user_send_history\\(\\$1\\);").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user", "amount"}).
			AddRow("user3", 80).
			AddRow("user4", 30),
		)
	mock.ExpectCommit()

	result, err := GetUserBalanceInventoryLogs(1)
	assert.NoError(t, err)
	assert.Equal(t, 250, result.Coins)
	expectedItems := []models.Item{
		{"t_shirt", 1},
		{"cup", 2},
		{"book", 1},
		{"powerbank", 1},
	}
	assert.Equal(t, expectedItems, result.Inventory)
	expectedReceived := []models.CoinTransaction{
		{"user1", 100},
		{"user2", 50},
	}
	assert.Equal(t, expectedReceived, result.CoinHistory.Received)
	expectedSent := []models.CoinTransaction{
		{"user3", 80},
		{"user4", 30},
	}
	assert.Equal(t, expectedSent, result.CoinHistory.Sent)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserBalanceInventoryLogsInvalid(t *testing.T) {
	resetMockDB(t)
	mock.ExpectBegin()
	returningError := fmt.Errorf("пользователь не найден")
	mock.ExpectQuery("SELECT get_user_balance\\(\\$1\\);").
		WithArgs(999).
		WillReturnError(returningError)
	mock.ExpectRollback()
	_, err := GetUserBalanceInventoryLogs(999)
	assert.ErrorIs(t, err, returningError)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func resetMockDB(t *testing.T) {
	var err error
	mockDB, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мока БД: %v", err)
	}
	db = mockDB
}
