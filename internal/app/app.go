package app

import (
	"log"

	"github.com/ZotovSergey/authenticationservice/internal/database/myProfilesDB"
	"github.com/ZotovSergey/authenticationservice/internal/services/rest/httprouter"
	// "github.com/ZotovSergey/authenticationservice/internal/models"
)

func Run() {
	// Поднятие БД
	log.Println("in memory database is raising...")
	err := myProfilesDB.RaiseMyProfilesDB()
	if err != nil {
		log.Printf("fatal error: %s", err.Error())
		return
	}
	log.Println("in memory database raised")

	// Развертывание API
	log.Println("api deployment...")
	err = httprouter.StartAPI()
	if err != nil {
		log.Printf("fatal error: %s", err.Error())
		return
	}

	// Если будет время, сделать shutdown

	// db, err := myProfilesDB.BuildMyProfilesDB("/Volumes/Storage/sergejzotov/Programming/Golang/Authentication service/databaseDumps/dbTest.json")
	// if err == nil {
	// 	fmt.Println(db)
	// } else {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Все логины")
	// logins := db.GetAllLogins()
	// fmt.Println(logins)
	// fmt.Println()
	// fmt.Println("Получение профиля")
	// profile, err := db.GetProfileData("cleric")
	// if err == nil {
	// 	fmt.Println(profile)
	// } else {
	// 	fmt.Println(err)
	// }
	// profile, err = db.GetProfileData("no login")
	// if err == nil {
	// 	fmt.Println(profile)
	// } else {
	// 	fmt.Println(err)
	// }
	// fmt.Println()
	// fmt.Println("Админ")
	// fmt.Println(db.IsAdmin("cleric"))
	// fmt.Println(db.IsAdmin("Little Fox"))
	// fmt.Println(db.IsAdmin("no login"))
	// fmt.Println()
	// fmt.Println("Добавление профиля")
	// err = db.AddProfile("NewLogin", models.ProfileData{Login: "NewLogin", FirstName: "Иван", LastName: "Иванов"}, "password")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(db.GetAllLogins())
	// }
	// err = db.AddProfile("cleric", models.ProfileData{Login: "cleric", FirstName: "Иван", LastName: "Иванов"}, "password")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(db.GetAllLogins())
	// }
	// fmt.Println()
	// fmt.Println("Редактирование профиля")
	// db.AddProfile("NewLogin2", models.ProfileData{Login: "NewLogin2", FirstName: "Иван", LastName: "Иванов"}, "password")
	// db.EditProfile("NewLogin2", models.ProfileData{Login: "NewLogin2", FirstName: "Александрович", LastName: "Александров"})
	// profile, _ = db.GetProfileData("NewLogin2")
	// fmt.Println(profile)
	// // Изменение пароля
	// db.ChangePassword("NewLogin2", "password2")
	// fmt.Println()
	// fmt.Println("Удаление профиля")
	// db.RemoveProfile("NewLogin")
	// fmt.Println(db.GetAllLogins())
	// fmt.Println()
	// fmt.Println("Добавление админа")
	// db.AddAdmin("Little Fox")
	// fmt.Println(db.IsAdmin("Little Fox"))
	// fmt.Println()
	// fmt.Println("Удаление админа")
	// db.DropAdmin("Cleric")
	// fmt.Println(db.IsAdmin("Cleric"))
	// db.Dump()
}
