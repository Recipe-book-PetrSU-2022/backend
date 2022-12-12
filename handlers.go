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

type UserData struct {
	Login           string `json:"login"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ChangeUserData struct {
	Login           string `json:"login"`
	Email           string `json:"email"`
	OldPassword     string `json:"old_password"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Photo           string `json:"photo"`
}

type RecipeData struct {
	Name     string `json:"name"`
	Servings int    `json:"servings"`
	Time     int    `json:"time"`
	Country  string `json:"country"`
	Type     string `json:"type"`
	Image    string `json:"image"`
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

func (server *Server) GetUserByClaims(c echo.Context) (*models.User, error) {
	user_token := c.Get("user").(*jwt.Token)
	user_claims := user_token.Claims.(*claims.UserClaims)

	var user models.User

	err := server.DB.First(&user, "id = ?", user_claims.IntUserId).Error
	if err != nil {
		return nil, errors.New("Пользователь не найден")
	}

	return &user, nil
}

func (server *Server) ProfileHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	response := &ProfileResponse{Message: "Удачный вход на страницу профиля", Id: user.ID, Username: user.StrUserName, Email: user.StrUserEmail, ProfilePhoto: user.StrUserImage}

	return c.JSON(http.StatusOK, response)
}

func (server *Server) ChangeProfileHandle(c echo.Context) error {
	// Получаем данные о пользователе
	user, err := server.GetUserByClaims(c)

	// Если не получилось найти пользователя
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	// Если пользователь нашёлся
	var user_data ChangeUserData

	err = c.Bind(&user_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", user_data)

	if len(user_data.Login) != 0 && len(user_data.Email) != 0 && len(user_data.OldPassword) == 0 {
		return c.JSON(http.StatusOK, &DefaultResponse{Message: "Нечего изменять"})
	}

	if len(user_data.Login) != 0 {
		if user_data.Login == user.StrUserName {
			// Пользователь не стал менять никнейм
		} else {
			// Пользователь решил поменять никнейм
			server.DB.Model(&user).Update("StrUserName", user_data.Login)
		}
	}

	if len(user_data.Email) != 0 {
		if user_data.Email == user.StrUserEmail {
			// Пользователь не стал менять почту
		} else {
			// Пользователь решил поменять почту
			server.DB.Model(&user).Update("StrUserEmail", user_data.Email)
		}
	}

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

	// Прошли все проверки

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Пользователь успешно изменил свои данные!"})
}

func (server *Server) DeleteProfileHandle(c echo.Context) error {
	// Получаем данные о пользователе
	user, err := server.GetUserByClaims(c)

	// Если не получилось найти пользователя
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пользователь не найден"})
	}

	// "Удаляем" запись о пользователе
	server.DB.Delete(&user)

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Пользователь успешно удалён!"})
}

func (server *Server) CreateRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) GetRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) GetRecipesHandle(c echo.Context) error {
	return nil
}

func (server *Server) UpdateRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) DeleteRecipeHandle(c echo.Context) error {
	return nil
}

func (server *Server) FindRecipesHandle(c echo.Context) error {
	return nil
}

func (server *Server) AddRecipeToFavoritesHandle(c echo.Context) error {
	return nil
}

func (server *Server) RemoveRecipeFromFavoritesHandle(c echo.Context) error {
	return nil
}
