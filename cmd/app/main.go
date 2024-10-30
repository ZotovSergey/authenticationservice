package main

import "github.com/ZotovSergey/authenticationservice/internal/app"

// @title Authentication service
// @verion 1.0
// @description Сервис для хранения информации о пользователях. Каждый зарегистрированный пользователь может авторизоваться по своим логину и паролю и просматривать и менять информацию о себе и только просматривать информацию о других пользователях. Пользователи-администраторы могут создавать, изменять и удалять профили.

// @host localhost:3000
// @BasePath /

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func main() {
	app.Run()
}
