package httprouter

import (
	_ "github.com/ZotovSergey/authenticationservice/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/ZotovSergey/authenticationservice/internal/rest/handlers"
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
	app := fiber.New()

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Аутентификация
	app.Use(handlers.BasicAuth())

	// Инициализация запросов
	app.Get("/profile", handlers.GetProfileDataRequest)    // запрос на получение данных о профиле
	app.Get("/logins", handlers.GetAllLoginsRequest)       // запрос на получение списка логинов всех пользователей
	app.Post("/profile", handlers.AddProfileRequest)       // запрос на добавление пользователя
	app.Patch("/profile", handlers.EditProfileRequest)     // запрос на изменение данных пользователя
	app.Patch("/password", handlers.ChangePasswordRequest) // запрос на изменение пароля профиля
	app.Delete("/profile", handlers.RemoveProfileRequest)  // запрос на удаление профиля
	app.Post("/admin", handlers.AddAdminRequest)           // запрос добавление администратора
	app.Delete("/admin", handlers.DropAdminRequest)        // запрос удаление администратора

	// Выыод API на порт из конфигурационного файла
	app.Listen(port)

	return nil
}
