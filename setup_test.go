package main

import (
	"os"
	"testing"

	"github.com/Recipe-book-PetrSU-2022/backend/claims"
	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	TestE      = echo.New()
	TestServer = Server{
		E:    TestE,
		Host: "0.0.0.0",
		Port: 11111,
		// DBConnectionInfo: "file::memory:/test?cache=shared", // БД в оперативке
		TokenKey:    []byte("test"),
		UploadsPath: "/tmp/test/recipe_book_uploads",
	}
	UserJWT           = ""
	UserJWT2          = ""
	TestJwtMiddleware echo.MiddlewareFunc
)

func TestMain(m *testing.M) {
	// err := testServer.ConnectDB()
	// if err != nil {
	// 	panic(err)
	// }

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	TestServer.DB = db

	config := middleware.JWTConfig{
		Claims:     &claims.UserClaims{},
		SigningKey: TestServer.TokenKey,
	}
	TestJwtMiddleware = middleware.JWTWithConfig(config)
	TestE.Use(TestJwtMiddleware)

	err = TestServer.DB.AutoMigrate(
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
		panic(err)
	}

	os.Exit(m.Run())
}
