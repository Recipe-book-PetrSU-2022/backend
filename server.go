package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Recipe-book-PetrSU-2022/backend/claims"
	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Структура сервера
//
// Переменные структуры:
//   - Хост для запуска
//   - Порт для запуска
//   - Echo http-сервер
//   - Информация для подключения
//   - Объект ORM
type Server struct {
	Host             string     // Хост для запуска
	Port             int        // Порт для запуска
	E                *echo.Echo // Echo http-сервер
	DBConnectionInfo string     // Информация для подключения
	DB               *gorm.DB   // Объект ORM
	TokenKey         []byte
}

// Функция для поднятия сервера
func (server *Server) Run() error {
	err := server.ConnectDB()

	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(
		&models.User{},
		&models.Filter{},
		&models.Ingredient{},
		&models.Recipe{},
		&models.Stage{},
		&models.Comment{},
		&models.Photo{},
		&models.RecipeIngredient{},
	)

	if err != nil {
		return err
	}

	server.E.Use(middleware.Logger())
	server.E.Use(middleware.Recover())

	config := middleware.JWTConfig{
		Claims:     &claims.UserClaims{},
		SigningKey: server.TokenKey,
	}

	jwtMiddleware := middleware.JWTWithConfig(config)

	server.E.GET("/profile", server.ProfileHandle, jwtMiddleware)
	server.E.POST("/signin", server.SignInHandle)
	server.E.POST("/signup", server.SignUpHandle)

	return server.E.Start(fmt.Sprintf("%s:%d", server.Host, server.Port))
}

// Функция для подключения сервера к БД
func (server *Server) ConnectDB() error {
	db, err := gorm.Open(mysql.Open(server.DBConnectionInfo), &gorm.Config{})

	if err != nil {
		return err
	}

	server.DB = db

	return nil
}

// Основная функция
func main() {
	e := echo.New()

	server := Server{
		E:                e,
		Host:             "0.0.0.0",
		Port:             1337,
		DBConnectionInfo: "root:my-secret-pw@tcp(127.0.0.1:3306)/recipe_book?charset=utf8mb4&parseTime=True",
		TokenKey:         []byte("rakabidasta_test_key"),
	}

	err := server.Run()

	if err != nil {
		log.Fatalf("Can't start server: %s", err.Error())
	}
}
