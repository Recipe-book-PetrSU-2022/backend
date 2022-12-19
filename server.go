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
	TokenKey         []byte     // ключ подписи токена
	UploadsPath      string     // путь для загрузки файлов
}

// Функция для поднятия сервера
//
// Первоначально устанавливается соединение с БД
// После этого мигрируются модели в БД
// Настраивается middleware, создаются группы для аккаунта и рецептов
// Прописываются эндпоинты
// И в конце запускается сам сервер
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

	err = server.CreateUploadDirs()
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

	recipe_group := server.E.Group("/recipe")                        // от лица кого угодно
	user_recipe_group := server.E.Group("/my-recipe", jwtMiddleware) // от лица владельца
	profile_group := server.E.Group("/profile", jwtMiddleware)
	assets_group := server.E.Group("/assets")

	server.E.POST("/signin", server.SignInHandle)
	server.E.POST("/signup", server.SignUpHandle)

	profile_group.GET("", server.ProfileHandle)
	profile_group.GET("/update", server.ChangeProfileHandle)
	profile_group.GET("/delete", server.DeleteProfileHandle)

	user_recipe_group.POST("/add", server.CreateEmptyRecipeHandle)
	user_recipe_group.POST("/complete/:id", server.UpdateRecipeHandle)
	user_recipe_group.POST("/visible/:id", server.ChangeVisibilityRecipeHandle)
	user_recipe_group.POST("/change/:id", server.UpdateRecipeHandle)
	user_recipe_group.POST("/delete/:id", server.DeleteRecipeHandle)
	user_recipe_group.POST("/upload-cover/:id", server.UploadRecipeCoverHandle)
	user_recipe_group.GET("/:id", server.GetMyRecipeHandle)
	user_recipe_group.GET("/all", server.GetMyRecipesHandle)

	user_recipe_group.POST("/:recipe_id/stage/add", server.CreateStageHandle)
	user_recipe_group.POST("/stage/:stage_id/update", server.UpdateStageHandle)
	user_recipe_group.POST("/stage/:stage_id/upload-photo", server.AddStagePhotoHandle)

	recipe_group.POST("/comment/:id", server.GetCommentHandle)
	recipe_group.POST("/comments", server.GetCommentsHandle)
	recipe_group.POST("/comment/:id/add", server.CreateCommentHandle)
	recipe_group.POST("/comment/:id/delete", server.DeleteCommentHandle)

	recipe_group.GET("/:id", server.GetRecipeHandle)
	recipe_group.GET("/all", server.GetRecipesHandle)
	recipe_group.GET("/find", server.FindRecipesHandle)
	recipe_group.GET("/favorite/:id", server.AddRecipeToFavoritesHandle, jwtMiddleware)

	assets_group.GET("/:filename", server.DownloadFile)

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
		UploadsPath:      "/tmp/recipe_book_uploads/",
	}

	err := server.Run()

	if err != nil {
		log.Fatalf("Can't start server: %s", err.Error())
	}
}
