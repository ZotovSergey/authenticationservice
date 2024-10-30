package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/ZotovSergey/authenticationservice/internal/database/myProfilesDB"
	"github.com/ZotovSergey/authenticationservice/internal/models"
)

// @Summary Get profile data
// @Security BasicAuth
// @Description Запрос на вывод данных о профиле по логину, доступно всем пользователям (запрос не работает со страницы swagger из браузера, но работает через postman или insomnia)
// @Accept json
// @Produce json
// @Param input body models.LoginData true "логин получаемого профиля"
// @Success      200  {json}	json	model.ProfileData
// @Failure      404  {string}  string	"no such profile"
// @Router /profile [get]
func GetProfileDataRequest(ctx *fiber.Ctx) error {
	// Чтение тела запроса
	var body models.LoginData
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

// @Summary Get all logins
// @Security BasicAuth
// @Description Запрос на вывод списка логинов всех профилей, доступно всем пользователям
// @Produce json
// @Success      200  {json}  json model.FullProfileData "logins"
// @Router /logins [get]
func GetAllLoginsRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "get all logins")

	// Получение списка логинов из БД
	loginsList := myProfilesDB.DB.GetAllLogins()

	log.Printf("request completed (status %d)", fiber.StatusOK)
	return ctx.JSON(loginsList)
}

// @Summary Add profile
// @Security BasicAuth
// @Description Запрос на регистрацию нового пользователя, доступно только администраторам
// @Accept json
// @Param input body models.FullProfileData true "данные нового профиля, включающие логин нового профиля, данные для хранения и пароль"
// @Success      200  {string}  string	"request completed"
// @Failure      404  {string}  string	"no such profile"
// @Router /profile [post]
func AddProfileRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "add profile")

	// Проверка, является ли авторизованный пользователь администратором
	authorizedUserLogin := ctx.Locals("login").(string)
	if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) {
		log.Println(userIsNotAdminErr.Error())
		return userIsNotAdminErr
	}

	// Чтение тела запроса
	var body models.FullProfileData
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
	return ctx.SendString("request completed")
}

// @Summary Edit profile
// @Security BasicAuth
// @Description Запрос на редактирование данных профиля, доступно всем пользователям при редактировании своих профилей и доступно редактирование всех профилей администраторам
// @Accept json
// @Param input body models.ProfileData true "логин редактируемого профиля и новые данные для редактирования (если данные не добавлениы, то они не меняются)"
// @Success      200  {string}  string	"request completed"
// @Failure      404  {string}  string	"no such profile"
// @Router /profile [patch]
func EditProfileRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "edit profile")

	// Чтение тела запроса
	var body models.ProfileData
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Проверка прав достува (является ли авторизованный пользователь администратором или пользователь редактирует свой профиль)
	authorizedUserLogin := ctx.Locals("login").(string)
	if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) && authorizedUserLogin != body.Login {
		if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) {
			log.Println(userIsNotAdminErr.Error())
			return userIsNotAdminErr
		}
		if authorizedUserLogin != body.Login {
			log.Println(canNotEditProfileErr.Error())
			return canNotEditProfileErr
		}
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
	return ctx.SendString("request completed")
}

// @Summary Change password
// @Security BasicAuth
// @Description Запрос на изменение пароля пользователя, доступно всем пользователям при редактировании своих профилей и доступно редактирование всех профилей администраторам
// @Accept json
// @Param input body models.NewPasswordForProfile true "логин профиля, для каторого меняется пароль и новый пароль"
// @Success      200  {string}  string	"request completed"
// @Failure      404  {string}  string	"no such profile"
// @Router /password [patch]
func ChangePasswordRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "change password")

	// Чтение тела запроса
	var body models.NewPasswordForProfile
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Проверка прав достува (является ли авторизованный пользователь администратором или пользователь редактирует свой профиль)
	authorizedUserLogin := ctx.Locals("login").(string)
	if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) && authorizedUserLogin != body.Login {
		if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) {
			log.Println(userIsNotAdminErr.Error())
			return userIsNotAdminErr
		}
		if authorizedUserLogin != body.Login {
			log.Println(canNotEditProfileErr.Error())
			return canNotEditPasswordErr
		}
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
	return ctx.SendString("request completed")
}

// @Summary Remove profile
// @Security BasicAuth
// @Description Запрос удаление профиля, доступно только администраторам, нельзя удалять свой профиль
// @Accept json
// @Param input body models.LoginData true "логин удаляемого профиля"
// @Success      200  {string}  string	"request completed"
// @Failure      404  {string}  string	"no such profile"
// @Router /profile [delete]
func RemoveProfileRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "remove profile")

	// Проверка, является ли авторизованный пользователь администратором
	authorizedUserLogin := ctx.Locals("login").(string)
	if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) {
		log.Println(userIsNotAdminErr.Error())
		return userIsNotAdminErr
	}

	// Чтение тела запроса
	var body models.LoginData
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Проверка профиля: нельзя удалять свой профиль
	if authorizedUserLogin == body.Login {
		log.Println(canNotRemoveOwnProfileErr.Error())
		return canNotRemoveOwnProfileErr
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
	return ctx.SendString("request completed")
}

// @Summary Add admin
// @Security BasicAuth
// @Description Запрос на добавление администратора, доступно только администраторам
// @Accept json
// @Param input body models.LoginData true "логин профиля, добавляемого к списку администраторов"
// @Success      200  {string}  string	"request completed"
// @Failure      404  {string}  string	"no such profile"
// @Router /admin [post]
func AddAdminRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "add admin")

	authorizedUserLogin := ctx.Locals("login").(string)
	if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) {
		log.Println(userIsNotAdminErr.Error())
		return userIsNotAdminErr
	}

	// Чтение тела запроса
	var body models.LoginData
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
	return ctx.SendString("request completed")
}

// @Summary Drop admin
// @Security BasicAuth
// @Description Запрос на удаление профиля из списка администраторов администратора, доступно только администраторам, нельзя удалять из списка администраторов свой профиль
// @Accept json
// @Param input body models.LoginData true "логин профиля, удаляемого из списка администраторов"
// @Success      200  {string}  string	"request completed"
// @Failure      404  {string}  string	"no such profile"
// @Router /admin [delete]
func DropAdminRequest(ctx *fiber.Ctx) error {
	log.Printf("\"%s\" request received", "drop admin")

	authorizedUserLogin := ctx.Locals("login").(string)
	if !myProfilesDB.DB.IsAdmin(authorizedUserLogin) {
		log.Println(userIsNotAdminErr.Error())
		return userIsNotAdminErr
	}

	// Чтение тела запроса
	var body models.LoginData
	err := ctx.BodyParser(&body)
	if err != nil {
		log.Printf("request error (status %d): %s", fiber.StatusBadRequest, err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Проверка профиля: нельзя удалять свой профиль из администраторов
	if authorizedUserLogin == body.Login {
		log.Println(canNotRemoveOwnProfileFromAdminsErr.Error())
		return canNotRemoveOwnProfileFromAdminsErr
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
	return ctx.SendString("request completed")
}
