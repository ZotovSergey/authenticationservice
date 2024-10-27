# Authentication service
## Описание
Сервис для хранения информации о пользователях.
Каждый зарегистрированный пользователь может авторизоваться по своим логину и паролю и просматривать информацию о себе и других пользователях.
Пользователи-администраторы могут создавать, изменять и удалять профили.
### Данные пользователей
Данные пользователей, доступные для просмотра:
* Логин
* Имя
* Фамилия
## База данных
В сервисе используется кастомная in memory база данных с возможностью сохранения данных на диск и поднятии при запуске сервиса.
База данных хранит:
* Данные о пользователях
* Пароли
* Список администраторов