package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Recipe-book-PetrSU-2022/backend/claims"
	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Структура ответа для входа и регистрации
//
// Переменные структуры:
//   - Никнейм
//   - Почта пользователя
//   - Пароль
//   - Пароль для подтверждения
type UserData struct {
	Login           string `json:"login"`            // Никнейм
	Email           string `json:"email"`            // Почта
	Password        string `json:"password"`         // Пароль
	ConfirmPassword string `json:"confirm_password"` // Подтверждение пароля
}

// Структура ответа для входа и регистрации
//
// Переменные структуры:
//   - Никнейм
//   - Почта пользователя
//   - Текущий пароль
//   - Пароль
//   - Пароль для подтверждения
//   - Фото профиля
type ChangeUserData struct {
	Login           string `json:"login"`            // Никнейм
	Email           string `json:"email"`            // Почта
	OldPassword     string `json:"old_password"`     // Текущий пароль
	Password        string `json:"password"`         // Новый пароль
	ConfirmPassword string `json:"confirm_password"` // Подтверждение нового пароля
	Photo           string `json:"photo"`            // Фото профиля
}

// Функция для регистрации
//
// Обрабатывает json с фронтэнда.
// Проверяет на наличие логина, почты, пароля и подтверждение пароля
// После прохождения проверок, хэширует пароль, создаёт пользователя
// и добавляет его в БД
func (server *Server) SignUpHandle(c echo.Context) error {
	//
	// Получение информации о пользователе с фронтэенда
	var user_data UserData
	err := c.Bind(&user_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", user_data)

	// Если введен пустой логин
	if len(user_data.Login) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Имя пользователя не может быть пустым"})
	}

	// Если введена пустая почта
	if len(user_data.Email) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Почта пользователя не может быть пустой"})
	}

	// Если введен пустой пароль
	if len(user_data.Password) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароль пользователя не может быть пустым"})
	}

	// Если введен пустой повторный пароль
	if user_data.Password != user_data.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароли должны совпадать"})
	}

	// Получаем хэш пароля
	passwordHash, err := HashPassword(user_data.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не получилось захешировать пароль"})
	}

	// Сохраняем данные о пользователе
	user := models.User{
		StrUserName:     user_data.Login,
		StrUserPassword: passwordHash,
		StrUserEmail:    user_data.Email,
		IntUserRights:   0,
	}

	// Добавляем пользователя в БД
	err = server.DB.Create(&user).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не получилось создать пользователя"})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Пользователь успешно зарегистрирован!"})
}

// Функция для регистрации
//
// Обрабатывает json с фронтэнда.
// Проверяет на наличие логина, пароля. Вход по почте пока что не сделан
// После прохождения проверок, ищет пользователя в БД
// Если пользователь найден, проверяет введенный пароль с сохранённым хэшем
// Если пароли совпали, то создаётся jwt
// Пример структуры токена в /claims/user_claims.go
func (server *Server) SignInHandle(c echo.Context) error {
	// Получаем данные от пользователя
	var user_data UserData
	err := c.Bind(&user_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", user_data)

	// Если введен пустой логин
	if len(user_data.Login) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Имя пользователя не может быть пустым"})
	}

	// Если введен пустой пароль
	if len(user_data.Password) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароль пользователя не может быть пустым"})
	}

	log.Printf("user = %+v", user_data)

	// Берём информацию о пользователе по логину
	var user models.User
	err = server.DB.First(&user, "str_user_name = ?", user_data.Login).Error
	log.Println(err)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	// Проверяем совпадает ли введенный пароль с сохраненным хэшем
	if !CheckPasswordHash(user.StrUserPassword, user_data.Password) {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Введены неверные данные"})
	}

	// Заполняем структуру для JWT
	user_claims := claims.UserClaims{
		IntUserId:     user.ID,
		StrUserName:   user.StrUserName,
		IntUserRights: user.IntUserRights,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	// Создаем JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user_claims)
	token_string, err := token.SignedString(server.TokenKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не получилось подписать токен"})
	}

	return c.JSON(http.StatusOK, &TokenResponse{Message: "Пользователь успешно вошёл в систему!", Token: token_string})
}

// Поиск пользователя в БД
//
// Обрабатывает jwt с фронтэнда.
// Берёт информацию о пользователе из jwt
// Ищет пользователя в БД по id из jwt
// Если пользователь найден, то возращает указатель на пользователя
// Иначе ошибку
func (server *Server) GetUserByClaims(c echo.Context) (*models.User, error) {
	// Получаем данные о JWT с фронтэнда
	user_token := c.Get("user").(*jwt.Token)
	user_claims := user_token.Claims.(*claims.UserClaims)

	// Ищем пользователя по ID из JWT
	var user models.User
	err := server.DB.First(&user, "id = ?", user_claims.IntUserId).Error
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return &user, nil
}

// Получение страницы пользователя
//
// Обрабатывает jwt с фронтэнда.
// Ищет пользователя в БД по jwt (GetUserByClaims)
// Если пользователь найден, то возвращает основную информацию о пользователе
// Иначе ошибку
func (server *Server) ProfileHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	// Создаем ответ с ID, логином, почтой и фотографией профиля
	response := &ProfileResponse{Message: "Удачный вход на страницу профиля", Id: user.ID, Username: user.StrUserName, Email: user.StrUserEmail, ProfilePhoto: user.StrUserImage}

	return c.JSON(http.StatusOK, response)
}

// Изменение данных о пользователе
//
// Обрабатывает jwt и json с фронтэнда.
// Ищет пользователя в БД по jwt (GetUserByClaims)
// Если пользователь найден, то возвращает основную информацию о пользователе
// Иначе ошибку
// Проверяет на наличие измененных данных
// Поочередно проверяет логин, почту, пароль
func (server *Server) ChangeProfileHandle(c echo.Context) error {
	// Получаем данные о пользователе
	user, err := server.GetUserByClaims(c)

	// Если не получилось найти пользователя
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	// Получаем данные с фронтэнда
	var user_data ChangeUserData
	err = c.Bind(&user_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", user_data)

	// Если пользователь ничего не меняет
	if len(user_data.Login) != 0 && len(user_data.Email) != 0 && len(user_data.OldPassword) == 0 {
		return c.JSON(http.StatusOK, &DefaultResponse{Message: "Нечего изменять"})
	}

	// Если пользователь ввёл что-то в поле логина
	if len(user_data.Login) != 0 {
		if user_data.Login == user.StrUserName {
			// Пользователь не стал менять никнейм
		} else {
			// Пользователь решил поменять никнейм
			server.DB.Model(&user).Update("StrUserName", user_data.Login)
		}
	}

	// Если пользователь ввёл что-то в поле почты
	if len(user_data.Email) != 0 {
		if user_data.Email == user.StrUserEmail {
			// Пользователь не стал менять почту
		} else {
			// Пользователь решил поменять почту
			server.DB.Model(&user).Update("StrUserEmail", user_data.Email)
		}
	}

	// Если пользователь ввёл что-то в поле старого пароля
	if len(user_data.OldPassword) == 0 {
		// Пользователь не ввёл текущий пароль
		if len(user_data.Password) != 0 || len(user_data.ConfirmPassword) != 0 {
			return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Вы не ввели текущий пароль"})
		}
	} else {
		if !CheckPasswordHash(user.StrUserPassword, user_data.OldPassword) {
			return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Введен неверный пароль"})
		}

		if user_data.Password != user_data.ConfirmPassword {
			return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Новые пароли должны совпадать"})
		}

		passwordHash, err := HashPassword(user_data.Password)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не получилось захешировать пароль"})
		}

		server.DB.Model(&user).Update("StrUserPassword", passwordHash)
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Пользователь успешно изменил свои данные!"})
}

// Удаление профиля пользователя
//
// Обрабатывает jwt с фронтэнда.
// Ищет пользователя в БД по jwt (GetUserByClaims)
// Если пользователь найден, то возвращает основную информацию о пользователе
// Иначе ошибку
// Удаляет пользователя
func (server *Server) DeleteProfileHandle(c echo.Context) error {
	// Получаем данные о пользователе
	user, err := server.GetUserByClaims(c)

	// Если не получилось найти пользователя
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	// "Удаляем" запись о пользователе
	err = server.DB.Delete(&user).Error

	// Если не получилось удалить пользователя
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось удалить пользователя"})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Пользователь успешно удалён!"})
}
