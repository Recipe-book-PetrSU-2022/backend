package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
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
	e                *echo.Echo // Echo http-сервер
	dbConnectionInfo string     // Информация для подключения
	db               *gorm.DB   // Объект ORM
}

// Функция для поднятия сервера
func (server *Server) Run() error {
	err := server.ConnectDB()

	if err != nil {
		return err
	}

	err = server.db.AutoMigrate(
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

	return server.e.Start(fmt.Sprintf("%s:%d", server.Host, server.Port))
}

// Функция для подключения сервера к БД
func (server *Server) ConnectDB() error {
	db, err := gorm.Open(mysql.Open(server.dbConnectionInfo), &gorm.Config{})

	if err != nil {
		return err
	}

	server.db = db

	return nil
}

// Основная функция
func main() {
	e := echo.New()

	server := Server{
		e:                e,
		Host:             "0.0.0.0",
		Port:             1337,
		dbConnectionInfo: "root:my-secret-pw@tcp(127.0.0.1:3306)/recipe_book?charset=utf8mb4&parseTime=True",
	}

	err := server.Run()

	if err != nil {
		log.Fatalf("Can't start server: %s", err.Error())
	}
}
