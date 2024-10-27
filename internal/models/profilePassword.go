package models

// Структура паролей (их хэшей и соли) профилей, содержащихся в БД
type ProfilePassword struct {
	PasswordHashSalt string `json:"passwordHashSalt"` // хэш+соль паролей профилей - используется для безопасного хранения паролей в зашифрованном виде
}
