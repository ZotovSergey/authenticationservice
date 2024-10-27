package requests

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/ZotovSergey/authenticationservice/internal/database/myProfilesDB"
	"github.com/ZotovSergey/authenticationservice/internal/models"
)

/*
Запрос на вывод данных о профиле по логину

:return: данные о профиле с заданным логином в виде json
*/
func GetProfileDataRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "get profile data")
	// Структура тела запроса
	type getProfileDataRequestBody struct {
		Login string `json:"login"` // логин получаемого профиля
	}

	// Чтение тела запроса
	var body getProfileDataRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Получение профиля из БД
	profileData, err := myProfilesDB.DB.GetProfileData(body.Login)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return ctx.JSON(profileData)
}

/*
Запрос на вывод списка логинов всех профилей

:return: список логинов всех профилей виде json
*/
func GetAllLoginsRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "get all logins")

	// Получение списка логинов из БД
	loginsList := myProfilesDB.DB.GetAllLogins()

	log.Printf("request completed (status %d)", fiber.StatusOK)
	return ctx.JSON(loginsList)
}

/*
Запрос на регистрацию нового пользователя

:return: статус ошибки, если она появляется
*/
func AddProfileRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "add profile")
	// Структура тела запроса
	type addProfileRequestBody struct {
		Login     string `json:"login"`     // логин нового профиля
		FirstName string `json:"firstName"` // имя нового пользователя
		LastName  string `json:"lastName"`  // фамилия нового пользователя
		Password  string `json:"password"`  // пароль для входа нового пользователя (должен быть не длиннее 72 символов)
	}

	// Чтение тела запроса
	var body addProfileRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Добавление пользователя
	err = myProfilesDB.DB.AddProfile(
		body.Login,
		models.ProfileData{
			Login:     body.Login,
			FirstName: body.FirstName,
			LastName:  body.LastName,
		},
		body.Password,
	)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return nil
}

/*
Запрос на регистрацию изменение данных пользователя

:return: статус ошибки, если она появляется
*/
func EditProfileRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "edit profile")
	// Структура тела запроса
	type editProfileRequestBody struct {
		Login     string `json:"login"`     // логин изменяемого профиля
		FirstName string `json:"firstName"` // новое имя пользователя
		LastName  string `json:"lastName"`  // новая фамилия пользователя
	}

	// Чтение тела запроса
	var body editProfileRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Получение текущих данных профиля
	currentProfileData, err := myProfilesDB.DB.GetProfileData(body.Login)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	// Составление новой структуры данных профиля
	var newProfileFirstName string // новое имя пользовате (остается прежним, если не было задано)
	if body.FirstName != "" {
		newProfileFirstName = body.FirstName
	} else {
		newProfileFirstName = currentProfileData.FirstName
	}
	var newProfileLastName string // новое имя пользовате (остается прежним, если не было задано)
	if body.LastName != "" {
		newProfileLastName = body.LastName
	} else {
		newProfileLastName = currentProfileData.LastName
	}

	// Редактирование данных профиля
	err = myProfilesDB.DB.EditProfile(
		body.Login,
		models.ProfileData{
			Login:     body.Login,
			FirstName: newProfileFirstName,
			LastName:  newProfileLastName,
		})
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return nil
}

/*
Запрос на изменение пароля профиля

:return: статус ошибки, если она появляется
*/
func ChangePasswordRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "change password")
	// Структура тела запроса
	type changePasswordRequestBody struct {
		Login       string `json:"login"`       // логин, у которого меняется пароль
		NewPassword string `json:"newPassword"` // новый пароль заданного профиля (должен быть не длиннее 72 символов)
	}

	// Чтение тела запроса
	var body changePasswordRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Изменение пароля пользователя
	err = myProfilesDB.DB.ChangePassword(
		body.Login,
		body.NewPassword,
	)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return nil
}

/*
Запрос на удаление профиля

:return: статус ошибки, если она появляется
*/
func RemoveProfileRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "remove profile")
	// Структура тела запроса
	type removeProfileRequestBody struct {
		Login string `json:"login"` // логин удаляемого профиля
	}

	// Чтение тела запроса
	var body removeProfileRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Удаление профиля
	err = myProfilesDB.DB.RemoveProfile(
		body.Login,
	)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return nil
}

/*
Запрос на добавление администратора

:return: статус ошибки, если она появляется
*/
func AddAdminRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "add admin")
	// Структура тела запроса
	type addAdminRequestBody struct {
		Login string `json:"login"` // логин добавляемого администратора
	}

	// Чтение тела запроса
	var body addAdminRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Добавление администратора
	err = myProfilesDB.DB.AddAdmin(
		body.Login,
	)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return nil
}

/*
Запрос на удаление администратора

:return: статус ошибки, если она появляется
*/
func DropAdminRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "drop admin")
	// Структура тела запроса
	type dropAdminRequest struct {
		Login string `json:"login"` // логин профиля, удаляемого из списка администраторов
	}

	// Чтение тела запроса
	var body dropAdminRequest
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Добавление администратора
	err = myProfilesDB.DB.DropAdmin(
		body.Login,
	)
	if err != nil {
		log.Printf("request completed (status %d) with error: %s", fiber.StatusOK, err.Error())
		return err
	}
	log.Printf("request completed (status %d)", fiber.StatusOK)
	return nil
}
