package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Recipe-book-PetrSU-2022/backend/claims"
	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserData struct {
	Login           string `json:"login"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (server *Server) SignUpHandle(c echo.Context) error {
	var user_data UserData

	err := c.Bind(&user_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", user_data)

	if len(user_data.Login) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Имя пользователя не может быть пустым"})
	}

	if len(user_data.Email) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Почта пользователя не может быть пустой"})
	}

	if len(user_data.Password) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароль пользователя не может быть пустым"})
	}

	if user_data.Password != user_data.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароли должны совпадать"})
	}

	// Прошли все проверки

	passwordHash, err := HashPassword(user_data.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не получилось захешировать пароль"})
	}

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

func (server *Server) SignInHandle(c echo.Context) error {
	var user_data UserData

	err := c.Bind(&user_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", user_data)

	if len(user_data.Login) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Имя пользователя не может быть пустым"})
	}

	// if len(user_data.Email) == 0 {
	// 	return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Почта пользователя не может быть пустой"})
	// }

	if len(user_data.Password) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароль пользователя не может быть пустым"})
	}

	// if user_data.Password != user_data.ConfirmPassword {
	// 	return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пароли должны совпадать"})
	// }

	// Прошли все проверки

	var user models.User

	err = server.DB.First(&user, "str_user_name = ?", user_data.Login).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	if !CheckPasswordHash(user.StrUserPassword, user_data.Password) {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Введены неверные данные"})
	}

	user_claims := claims.UserClaims{
		IntUserId:     user.ID,
		StrUserName:   user.StrUserName,
		IntUserRights: user.IntUserRights,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user_claims)

	token_string, err := token.SignedString(server.TokenKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не получилось подписать токен"})
	}

	return c.JSON(http.StatusOK, &TokenResponse{Message: "Пользователь успешно вошёл в систему!", Token: token_string})
}

func (server *Server) ProfileHandle(c echo.Context) error {

	return nil
}
