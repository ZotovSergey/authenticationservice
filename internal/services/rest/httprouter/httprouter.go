package httprouter

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ZotovSergey/authenticationservice/internal/services/rest/requests"
)

/*
Развертывание API; порт для API берется из конфигурационного файла config/apiConfig.json
*/
func StartAPI() error {
	// Чтение номера порта из конфига
	port, err := getPort()
	if err != nil {
		return err
	}

	// Развертывание API
	api := fiber.New()

	// Инициализация запросов
	api.Get("/profile", requests.GetProfileDataRequest)    // запрос на получение данных о профиле
	api.Get("/logins", requests.GetAllLoginsRequest)       // запрос на получение списка логинов всех пользователей
	api.Post("/profile", requests.AddProfileRequest)       // запрос на добавление пользователя
	api.Patch("/profile", requests.EditProfileRequest)     // запрос на изменение данных пользователя
	api.Patch("/password", requests.ChangePasswordRequest) // запрос на изменение пароля профиля
	api.Delete("/profile", requests.RemoveProfileRequest)  // запрос на удаление профиля
	api.Post("/admin", requests.AddAdminRequest)           // запрос добавление администратора
	api.Delete("/admin", requests.DropAdminRequest)        // запрос удаление администратора

	// Выыод API на порт из конфигурационного файла
	api.Listen(port)

	return nil
}
