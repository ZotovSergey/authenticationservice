package authorizers

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/ZotovSergey/authenticationservice/internal/database/myProfilesDB"
)

/*
Авторизация любого пользователя по логину и паролю

:param login string: логин для авторизации пользователя
:param password string: пароль для авторизации пользователя

:return: true - если логин и пароль есть в БД, иначе - false
*/
func CommonUsersAuthorizer(login, password string) bool {
	log.Println("attempt to authorize...")
	// Получение пароля из БД
	passwordHashSalt, err := myProfilesDB.DB.GetPasswordHashSalt(login)
	if err != nil {
		log.Println("access denied: no such profile")
		return false
	}
	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(passwordHashSalt), []byte(password))
	if err != nil {
		log.Println("access denied: wrong password")
		return false
	}
	log.Printf("access granted to user \"%s\"", login)
	return true
}
