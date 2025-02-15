package transport

import (
	"avito_internship/internal/auth"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Authentication обрабатывает запрос на аутентификацию пользователей.
// Ожидает POST-запрос с JSON-данными, содержащими учетные данные пользователя (имя и пароль).
// Если метод запроса не POST, возвращает ошибку 405 (Method Not Allowed).
// Если данные запроса некорректны или не могут быть разобраны, возвращает ошибку 400 (Bad Request).
// Если аутентификация не удалась (неверные учетные данные), возвращает ошибку 401 (Unauthorized).
// В случае успешной аутентификации возвращает токен в формате JSON и статус 200 (OK).
func Authentication(w http.ResponseWriter, r *http.Request, authFunc func(string, string) (string, error)) {
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

	var credentials AuthRequest
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

// invalidRequestMethodResponse генерирует сообщение об ошибке неверного типа запроса.
// Отправляет статус 405 (Method Not Allowed) с описанием ошибки в формате JSON.
func invalidRequestMethodResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: fmt.Sprintf("Метод %s не разрешен.", r.Method)})
}

// badRequestResponse генерирует сообщение об ошибке неверного запроса.
// Отправляет статус 400 (Bad Request) с общей ошибкой в формате JSON.
func badRequestResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: "Неверный запрос."})
}

// tokenResponse отправляет ответ с токеном в формате JSON.
// Отправляет статус 200 (OK) и переданный токен в теле ответа.
func tokenResponse(w http.ResponseWriter, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

// unauthorizedResponse генерирует ответ об ошибке авторизации.
// Отправляет статус 401 (Unauthorized) с общей ошибкой в формате JSON.
func unauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: "Неавторизован."})
}

// internalServerErrorResponse генерирует ответ о внутренней ошибке сервера.
// Отправляет статус 500 (Internal Server Error) с общей ошибкой в формате JSON.
func internalServerErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: "Внутренняя ошибка сервера."})
}
