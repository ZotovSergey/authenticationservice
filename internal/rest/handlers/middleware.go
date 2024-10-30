package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"

	"github.com/ZotovSergey/authenticationservice/internal/authorizers"
	"github.com/ZotovSergey/authenticationservice/internal/database/myProfilesDB"
)

/*
Функция возвращает функцию для middleware basic auth с кастомным config
*/
func BasicAuth() func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Users:      myProfilesDB.DB.GetPasswordsTab(),
		Authorizer: authorizers.CommonUsersAuthorizer,
		Unauthorized: func(ctx *fiber.Ctx) error {
			log.Println(unauthorizedRequestErr.Error())
			return ctx.SendString(unauthorizedRequestErr.Error())
		},
		ContextUsername: "login",
	})
}
