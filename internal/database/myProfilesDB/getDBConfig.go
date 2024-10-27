package myProfilesDB

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ZotovSergey/authenticationservice/internal/models"
)

/*
Получение пути к дампу БД конфига "../../configs/dbConfig.json"

:return: путь к файлу дампа БД или ошибка, если конфиг не удалось прочитать
*/
func getDBDumpPath() (string, error) {
	// Структура пути к файлу дампа в конфигурации БД
	type dbConfig struct {
		DatabaseDumpPath string `json:"databaseDumpPath"` // путь к дампу БД
	}
	configFilePath := "../../configs/dbConfig.json"
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return "", errors.New("fail to read database config " + configFilePath)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var dbData dbConfig
	json.Unmarshal(byteValue, &dbData)

	return dbData.DatabaseDumpPath, nil
}

/*
Получение данных профиля администратора по умолчанию и его пароля из БД конфига "../../configs/dbConfig.json"

:return: логин, данные и пароль профиля администратора по умолчанию или ошибка, если конфиг не удалось прочитать
*/
func getDefaultAdminProfile() (string, models.ProfileData, string, error) {
	// Структура данных стандартного профиля администратора в конфигурации БД
	type defaultAdminProfileConfig struct {
		DefaultAdminProfile struct {
			Login     string `json:"login"`     // логин стандартного профиля администратора
			FirstName string `json:"firstName"` // имя пользователя стандартного профиля администратора
			LastName  string `json:"lastName"`  // фамилия пользователя стандартного профиля администратора
			Password  string `json:"password"`  // пароль стандартного профиля администратора (должен быть не длиннее 72 символов)
		}
	}

	configFilePath := "../../configs/dbConfig.json"
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return "", models.ProfileData{}, "", errors.New("fail to read database config " + configFilePath)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	var dab defaultAdminProfileConfig
	json.Unmarshal(byteValue, &dab)

	return dab.DefaultAdminProfile.Login,
		models.ProfileData{
			Login:     dab.DefaultAdminProfile.Login,
			FirstName: dab.DefaultAdminProfile.FirstName,
			LastName:  dab.DefaultAdminProfile.LastName},
		dab.DefaultAdminProfile.Password, nil
}
