package myProfilesDB

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/ZotovSergey/authenticationservice/internal/models"
)

// Указатель на используемую базу данных
var DB *myProfilesDB

// Структура in memory базы данных профилей
type myProfilesDB struct {
	profilesDataTab      map[string]models.ProfileData // таблица данных профиля
	profilesPasswordsTab map[string]string             // таблица зашифрованных паролей хэш+соль паролей профилей - используется для безопасного хранения паролей в зашифрованном виде)
	adminsTab            map[string]struct{}           // таблица админов сервиса
	dumpFilePath         string                        // путь к файлу с данными из базы на диске (из него данные для заполнения читаются и в него сохраняются)
}

// Структура базы данных профилей в файле (для хранения данных в файле)
type myProfilesDBFileData struct {
	ProfilesDataTab      map[string]models.ProfileData `json:"profilesDataTab"`      // таблица данных профиля
	ProfilesPasswordsTab map[string]string             `json:"profilesPasswordsTab"` // таблица зашифрованных паролей хэш+соль паролей профилей - используется для безопасного хранения паролей в зашифрованном виде)
	AdminsTab            []string                      `json:"adminsTab"`            // таблица админов сервиса
}

/*
"Поднятие базы данных" - построение структуры базы данных MyProfilesDataBase ее и заполнение по заданному json-файлу
или построение пустого экземпляра базы, если по заданному пути файл отсутствует и ее запись по глобальному адресу DB

:return: ошибка, если базу данных не удается поднять
*/
func RaiseMyProfilesDB() error {
	// Чтение пути к файлу с данными для заполнения базы (и для сохранения данных в него)
	dataFilePath, err := getDBDumpPath()
	if err != nil {
		return err
	}
	var db myProfilesDB
	// Проверка, существует ли файл с данными по пути dataFilePath
	_, err = os.Stat(dataFilePath)
	switch {
	// Если файл существует, из файла читаются данные для записи в экземпляр БД
	case err == nil:
		// Чтение данных из файла для записи в экземпляр БД со структурой myProfilesDBFileData
		dataFile, err := os.Open(dataFilePath)
		if err != nil {
			return err
		}
		defer dataFile.Close()
		byteValue, _ := ioutil.ReadAll(dataFile)
		var dbFromFile myProfilesDBFileData
		json.Unmarshal(byteValue, &dbFromFile)

		// Составление in memory БД со структурой myProfilesDB
		//	map для списка админов
		adminsTab := make(map[string]struct{})
		for _, adminLogin := range dbFromFile.AdminsTab {
			adminsTab[adminLogin] = struct{}{}
		}
		db = myProfilesDB{
			profilesDataTab:      dbFromFile.ProfilesDataTab,
			profilesPasswordsTab: dbFromFile.ProfilesPasswordsTab,
			adminsTab:            adminsTab,
			dumpFilePath:         dataFilePath,
		}

	// Если файла не существует, создается пустой экземпляр БД
	case errors.Is(err, os.ErrNotExist):
		db = myProfilesDB{
			profilesDataTab:      make(map[string]models.ProfileData),
			profilesPasswordsTab: make(map[string]string),
			adminsTab:            make(map[string]struct{}),
			dumpFilePath:         dataFilePath,
		}
		// Добавление профиля администратора по умолчанию
		//	Чтение профиля администратора по умолчанию
		defaultAdminLogin, defaultAdminProfilrData, defaultAdminPassword, err := getDefaultAdminProfile()
		if err != nil {
			return err
		}
		//	Добавление профиля
		db.AddProfile(defaultAdminLogin, defaultAdminProfilrData, defaultAdminPassword)
		db.AddAdmin(defaultAdminLogin)
	}

	// Запись базы в глобальный указатель DB
	DB = &db
	return nil
}

/*
Сохранение данных из БД на диск в файл db.dumpFilePath

:return: возвращаеся ошибка, если файл с данными не удается переписать
*/
func (db *myProfilesDB) Dump() error {
	adminsList := make([]string, 0, len(db.adminsTab))
	for adminLogin := range db.adminsTab {
		adminsList = append(adminsList, adminLogin)
	}
	dbToFile := myProfilesDBFileData{
		ProfilesDataTab:      db.profilesDataTab,
		ProfilesPasswordsTab: db.profilesPasswordsTab,
		AdminsTab:            adminsList,
	}

	dataFile, _ := json.MarshalIndent(dbToFile, "", "	")
	err := ioutil.WriteFile(db.dumpFilePath, dataFile, 0644)
	if err != nil {
		return dbDumpFailErr
	}

	return nil
}

/*
Получить данные о профиле по логину

:param login string: логин получаемого профиля

:return: данные о профиле с логином login или ошибка, исли профиля с таким логином нет
*/
func (db *myProfilesDB) GetProfileData(login string) (profileData models.ProfileData, err error) {
	// Поиск данных профиля
	profileData, ok := db.profilesDataTab[login]
	if !ok {
		err = noProfileErr
	}
	return
}

/*
Получить список логинов всех зарегистрированных пользователей

:return: список логинов всех зарегистрированных пользователей
*/
func (db *myProfilesDB) GetAllLogins() (logins []string) {
	for login := range db.profilesDataTab {
		logins = append(logins, login)
	}
	return
}

/*
Получение зашифрованного пароля профиля по логину

:param login string: логин профиля, для которого берется зашифрованный пароль

:return: зашифрованный пароль профиля с логином login или ошибка, исли профиля с таким логином нет
*/
func (db *myProfilesDB) GetPasswordHashSalt(login string) (passwordHashSalt string, err error) {
	// Поиск пароля
	passwordHashSalt, ok := db.profilesPasswordsTab[login]
	if !ok {
		err = noProfileErr
	}
	return
}

/*
Получение таблицы логинов и паролей (для авторизаторов)

:return: map db.profilesPasswordsTab
*/
func (db *myProfilesDB) GetPasswordsTab() map[string]string {
	return db.profilesPasswordsTab
}

/*
Проверка профиля, входит ли он в список администраторов

:param login string: логин профиля, проверяемого по списку администраторов

:return: возвращается true, если login в списке администраторов, иначе - false и ошибка
*/
func (db *myProfilesDB) IsAdmin(login string) bool {
	// Проверка логина на наличие в списке администраторов
	_, ok := db.adminsTab[login]
	return ok
}

/*
Запись данных и пароля (в зашифрованном виде) нового пользователя (регистрация)

:param login string: логин нового профиля
:param profileData models.ProfileData: данные для хранения в новом профиле
:param password string: пароль для входа нового пользователя (должен быть не длиннее 72 символов)

:return: возвращается ошибка, если профиль с логином login уже существует, неподходящий пароль или базу данных не удалось сохранить
*/
func (db *myProfilesDB) AddProfile(login string, profileData models.ProfileData, password string) error {
	// Проверка наличия профиля с логином login
	_, ok := db.profilesDataTab[login]
	if ok {
		return profileExistsErr
	}

	// Добавление данных пользователя
	db.profilesDataTab[login] = profileData

	// Добавление зашифрованного пароля
	//	Генерация хэша и соли
	passwordHashSalt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return incorrectPasswordErr
	}
	db.profilesPasswordsTab[login] = string(passwordHashSalt)

	// Сохранение данных в файл
	err = db.Dump()
	if err != nil {
		return err
	}
	return nil
}

/*
Замена данных в профиле на новые

:param login string: логин редактируемого профиля
:param profileData models.ProfileData: новые данные для хранения в профиле

:return: возвращается ошибка, если профиль с логином login не существует или базу данных не удалось сохранить
*/
func (db *myProfilesDB) EditProfile(login string, profileData models.ProfileData) error {
	// Проверка наличия профиля с логином login
	_, ok := db.profilesDataTab[login]
	if !ok {
		return noProfileErr
	}

	// Замена данных в профиле на новые
	db.profilesDataTab[login] = profileData

	// Сохранение данных в файл
	err := db.Dump()
	if err != nil {
		return err
	}
	return nil
}

/*
Изменение пароля в профиле

:param login string: логин, у которого меняется пароль
:param newPassword string: новый пароль заданного профиля (должен быть не длиннее 72 символов)

:return: возвращается ошибка, если профиль с логином login не существует, неподходящий пароль или базу данных не удалось сохранить
*/
func (db *myProfilesDB) ChangePassword(login string, newPassword string) error {
	// Проверка наличия профиля с логином login
	_, ok := db.profilesDataTab[login]
	if !ok {
		return noProfileErr
	}

	// Добавление зашифрованного пароля
	//	Генерация хэша и соли
	newPasswordHashSalt, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return incorrectPasswordErr
	}
	db.profilesPasswordsTab[login] = string(newPasswordHashSalt)

	// Сохранение данных в файл
	err = db.Dump()
	if err != nil {
		return err
	}
	return nil
}

/*
Удаление профиля

:param login string: логин удаляемого профиля

:return: возвращается ошибка, если профиль с логином login не существует или базу данных не удалось сохранить
*/
func (db *myProfilesDB) RemoveProfile(login string) error {
	// Проверка наличия профиля с логином login
	_, ok := db.profilesDataTab[login]
	if !ok {
		return noProfileErr
	}

	// Удаление профиля из всех таблиц
	if db.IsAdmin(login) {
		delete(db.adminsTab, login)
	}
	delete(db.profilesDataTab, login)
	delete(db.profilesPasswordsTab, login)

	// Сохранение данных в файл
	err := db.Dump()
	if err != nil {
		return err
	}
	return nil
}

/*
Добавление логина профиля в список админов

:param login string: логин добавляемого администратора

:return: возвращается ошибка, если профиль с логином login не существует или базу данных не удалось сохранить
*/
func (db *myProfilesDB) AddAdmin(login string) error {
	// Проверка наличия профиля с логином login
	_, ok := db.profilesDataTab[login]
	if !ok {
		return noProfileErr
	}

	// Добавление логина login в список администраторов
	db.adminsTab[login] = struct{}{}

	// Сохранение данных в файл
	err := db.Dump()
	if err != nil {
		return err
	}
	return nil
}

/*
Удаление логина профиля из списка админов

:param login string: логин профиля, удаляемого из списка администраторов

:return: возвращается ошибка, если профиль с логином login не существует или базу данных не удалось сохранить
*/
func (db *myProfilesDB) DropAdmin(login string) error {
	// Проверка наличия профиля с логином login
	_, ok := db.profilesDataTab[login]
	if !ok {
		return noProfileErr
	}

	// Удаление профиля из списка администраторов
	delete(db.adminsTab, login)

	// Сохранение данных в файл
	err := db.Dump()
	if err != nil {
		return err
	}
	return nil
}
