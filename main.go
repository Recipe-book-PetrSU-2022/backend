package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Recipe-book-PetrSU-2022/backend/claims"
	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/Recipe-book-PetrSU-2022/backend/docs"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
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
	//
	// Подключение к БД
	err := server.ConnectDB()
	if err != nil {
		return err
	}

	// Автомиграция моделей
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

	// Создание директорий для хранения файлов
	err = server.CreateUploadDirs()
	if err != nil {
		return err
	}

	// Использование middleware
	server.E.Use(middleware.Logger())
	server.E.Use(middleware.Recover())
	config := middleware.JWTConfig{
		Claims:     &claims.UserClaims{},
		SigningKey: server.TokenKey,
	}
	jwtMiddleware := middleware.JWTWithConfig(config)

	server.E.Use(middleware.CORS())
	server.E.GET("/swagger/*", echoSwagger.WrapHandler)

	// Создание групп для применения middleware
	recipe_group := server.E.Group("/recipe") // от лица кого угодно
	ingredient_group := server.E.Group("/ingredient")
	user_recipe_group := server.E.Group("/my-recipe", jwtMiddleware) // от лица владельца
	profile_group := server.E.Group("/profile", jwtMiddleware)
	assets_group := server.E.Group("/assets")

	// Эндпоинты для регистрации логина
	server.E.POST("/signin", server.SignInHandle)
	server.E.POST("/signup", server.SignUpHandle)

	// Эндпоинты для администрирования профиля
	profile_group.GET("", server.ProfileHandle)
	profile_group.POST("/update", server.ChangeProfileHandle)
	profile_group.DELETE("/delete", server.DeleteProfileHandle)

	// Эндпоинты для работы с рецептом
	user_recipe_group.POST("/add", server.CreateEmptyRecipeHandle)
	user_recipe_group.POST("/complete/:id", server.UpdateRecipeHandle)
	user_recipe_group.POST("/visible/:id", server.ChangeVisibilityRecipeHandle)
	user_recipe_group.POST("/change/:id", server.UpdateRecipeHandle)
	user_recipe_group.DELETE("/delete/:id", server.DeleteRecipeHandle)
	user_recipe_group.POST("/upload-cover/:id", server.UploadRecipeCoverHandle)
	user_recipe_group.GET("/:id", server.GetMyRecipeHandle)
	user_recipe_group.GET("/all", server.GetMyRecipesHandle)

	// Эндпоинты для работы с этапами
	user_recipe_group.POST("/:recipe_id/stage/add", server.CreateStageHandle)
	user_recipe_group.POST("/:recipe_id/ingredient/add", server.AddIngredientHandle)
	user_recipe_group.DELETE("/:recipe_id/ingredient/delete", server.RemoveIngredientHandle)
	user_recipe_group.DELETE("/stage/:stage_id/delete", server.DeleteStageHandle)
	user_recipe_group.POST("/stage/:stage_id/update", server.UpdateStageHandle)
	user_recipe_group.POST("/stage/:stage_id/upload-photo", server.AddStagePhotoHandle)

	// Эндпоинты для работы с комментариями
	recipe_group.GET("/:recipe_id/comment/:comment_id", server.GetCommentHandle)
	recipe_group.GET("/:recipe_id/comments", server.GetCommentsHandle)
	recipe_group.POST("/:recipe_id/comment/add", server.CreateCommentHandle, jwtMiddleware)
	recipe_group.DELETE("/:recipe_id/comment/:comment_id/delete", server.DeleteCommentHandle, jwtMiddleware)

	// Эндпоинты для работы с группой рецептов
	recipe_group.GET("/:id", server.GetRecipeHandle)
	recipe_group.GET("/all", server.GetRecipesHandle)
	recipe_group.GET("/find", server.FindRecipesHandle)
	recipe_group.POST("/favorite/:id", server.AddRecipeToFavoritesHandle, jwtMiddleware)

	ingredient_group.GET("/all", server.GetIngredients)
	ingredient_group.GET("/create", server.GetIngredients, jwtMiddleware)

	// Эндпоинты для работы с файлами
	assets_group.GET("/:filename", server.DownloadFile)

	return server.E.Start(fmt.Sprintf("%s:%d", server.Host, server.Port))
}

// Функция для подключения сервера к БД
func (server *Server) ConnectDB() error {
	// Соединение с БД
	db, err := gorm.Open(mysql.Open(server.DBConnectionInfo), &gorm.Config{})
	if err != nil {
		return err
	}

	server.DB = db

	return nil
}

//	@title			Recipe Book API
//	@description	Тут будет описание проекта

//	@version	1.0

//	@securityDefinitions.apikey	JWTAuth
//	@in							header
//	@name						Authorization
//	@description				JWT токен пользователя

//	@host	localhost:1337

// Основная функция
func main() {
	e := echo.New()

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")

	server := Server{
		E:                e,
		Host:             "0.0.0.0",
		Port:             1337,
		DBConnectionInfo: fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/recipe_book?charset=utf8mb4&parseTime=True", mysqlUser, mysqlPass),
		TokenKey:         []byte("rakabidasta_test_key"),
		UploadsPath:      "/tmp/recipe_book_uploads/",
	}

	// Запуск сервера
	err := server.Run()

	if err != nil {
		log.Fatalf("Can't start server: %s", err.Error())
	}
}
