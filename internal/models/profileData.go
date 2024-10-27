package models

// Структура данных профилей, содержащихся в БД, для хранения и вывода
type ProfileData struct {
	Login     string `json:"login"`     // логин профиля, используется при авторизации, является первичным ключом для БД, должен быть уникальным
	FirstName string `json:"firstName"` // имя пользователя
	LastName  string `json:"lastName"`  // фамилия пользователя
}
