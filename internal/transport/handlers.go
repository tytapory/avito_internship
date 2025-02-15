package transport

import (
	"avito_internship/internal/auth"
	"avito_internship/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Authenticate это middleware который отвечает за проверку предоставленного jwt токена.
// Он парсит токен и если он валидный то передает найденный в нем айди пользователя в handler
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				badRequestResponse(w)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := auth.VerifyJWT(tokenString)
			if err != nil {
				unauthorizedResponse(w)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userID", userID)))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// BuyItems обрабатывает покупку предметов пользователем.
// Ожидает GET-запрос по пути "/api/buy/{item}", где {item} — название предмета.
// Извлекает идентификатор пользователя из контекста, переданного через middleware Authenticate.
// Вызывает переданную функцию buyFunc с параметрами: userID, название предмета и количество (1).
// Если метод запроса не GET или URL не соответствует формату, возвращает ошибку 400 (Bad Request).
// Если во время покупки произошла ошибка, возвращает ошибку 500 (Internal Server Error).
func BuyItems(w http.ResponseWriter, r *http.Request, buyFunc func(int, string, int) error) {
	if r.Method != "GET" {
		badRequestResponse(w)
		return
	}
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[1] != "api" || parts[2] != "buy" {
		badRequestResponse(w)
		return
	}
	item := parts[3]
	err := buyFunc(r.Context().Value("userID").(int), item, 1)
	if err != nil {
		badRequestResponse(w)
		return
	}
}

// TransferCoins осуществляет перевод от одного пользователя к другому.
// Ожидает POST-запрос с JSON-данными, содержащими сумму перевода и ID получателя.
// Если метод запроса не POST, возвращает ошибку 400 (Bad Request).
// Если тело запроса не удалось прочитать или распарсить, возвращает ошибку 400 (Bad Request).
// Извлекает ID отправителя из контекста, переданного middleware Authenticate.
// Если перевод успешен, возвращает статус 200 (OK), иначе 400 (Bad Request).
func TransferCoins(w http.ResponseWriter, r *http.Request, transferFunc func(int, int, string) error) {
	if r.Method != http.MethodPost {
		invalidRequestMethodResponse(w, r)
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		badRequestResponse(w)
		return
	}
	var transferData models.SendCoinRequest
	err = json.Unmarshal(body, &transferData)
	if err != nil {
		badRequestResponse(w)
		return
	}
	err = transferFunc(r.Context().Value("userID").(int), transferData.Amount, transferData.ToUser)
	if err != nil {
		badRequestResponse(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetJWT обрабатывает запрос на аутентификацию пользователей.
// Ожидает POST-запрос с JSON-данными, содержащими учетные данные пользователя (имя и пароль).
// Если метод запроса не POST, возвращает ошибку 405 (Method Not Allowed).
// Если данные запроса некорректны или не могут быть разобраны, возвращает ошибку 400 (Bad Request).
// Если аутентификация не удалась (неверные учетные данные), возвращает ошибку 401 (Unauthorized).
// В случае успешной аутентификации возвращает токен в формате JSON и статус 200 (OK).
func GetJWT(w http.ResponseWriter, r *http.Request, authFunc func(string, string) (string, error)) {
	if r.Method != http.MethodPost {
		invalidRequestMethodResponse(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		badRequestResponse(w)
		return
	}

	var credentials models.AuthRequest
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		badRequestResponse(w)
		return
	}

	token, err := authFunc(credentials.Username, credentials.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) || errors.Is(err, auth.ErrExpiredToken) {
			unauthorizedResponse(w)
		} else {
			internalServerErrorResponse(w)
		}
		return
	}

	tokenResponse(w, token)
}

// GetUserInfo обрабатывает GET-запрос для получения информации о пользователе.
// Ожидает заголовок Authorization с валидным JWT-токеном, который уже был обработан middleware.
// Если метод запроса не GET, возвращает ошибку 405 (Method Not Allowed).
// Если возникает ошибка при получении данных, возвращает 500 (Internal Server Error).
// В случае успеха возвращает информацию о пользователе в формате JSON со статусом 200 (OK).
func GetUserInfo(w http.ResponseWriter, r *http.Request, userInfoFunc func(int) (models.InfoResponse, error)) {
	if r.Method != http.MethodGet {
		invalidRequestMethodResponse(w, r)
		return
	}
	info, err := userInfoFunc(r.Context().Value("userID").(int))
	if err != nil {
		internalServerErrorResponse(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}

// invalidRequestMethodResponse генерирует сообщение об ошибке неверного типа запроса.
// Отправляет статус 405 (Method Not Allowed) с описанием ошибки в формате JSON.
func invalidRequestMethodResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(models.ErrorResponse{Errors: fmt.Sprintf("Метод %s не разрешен.", r.Method)})
}

// badRequestResponse генерирует сообщение об ошибке неверного запроса.
// Отправляет статус 400 (Bad Request) с общей ошибкой в формате JSON.
func badRequestResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(models.ErrorResponse{Errors: "Неверный запрос."})
}

// tokenResponse отправляет ответ с токеном в формате JSON.
// Отправляет статус 200 (OK) и переданный токен в теле ответа.
func tokenResponse(w http.ResponseWriter, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{Token: token})
}

// unauthorizedResponse генерирует ответ об ошибке авторизации.
// Отправляет статус 401 (Unauthorized) с общей ошибкой в формате JSON.
func unauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(models.ErrorResponse{Errors: "Неавторизован."})
}

// internalServerErrorResponse генерирует ответ о внутренней ошибке сервера.
// Отправляет статус 500 (Internal Server Error) с общей ошибкой в формате JSON.
func internalServerErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(models.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
}
