package transport

import (
	"encoding/json"
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
func Authentication(w http.ResponseWriter, r *http.Request, authenticator Authenticator) {
	if r.Method != http.MethodPost {
		invalidRequestMethodResponse(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		badRequestResponse(w)
		return
	}
	defer r.Body.Close()

	var credentials AuthRequest
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		badRequestResponse(w)
		return
	}

	token, err := authenticator.Authenticate(credentials)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			unauthorizedResponse(w)
		} else {
			internalServerErrorResponse(w)
		}
		return
	}

	tokenResponse(w, token)
}

// Authenticator - интерфейс для аутентификации
type Authenticator interface {
	Authenticate(credentials AuthRequest) (string, error)
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

// unauthorizedResponse генерирует ответ о внутренней ошибке сервера.
// Отправляет статус 500 (Internal Server Error) с общей ошибкой в формате JSON.
func internalServerErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: "Внутренняя ошибка сервера."})
}
