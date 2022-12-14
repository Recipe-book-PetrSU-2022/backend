package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
)

type RecipeData struct {
	Name     string `json:"name"`
	Servings int    `json:"servings"`
	Time     int    `json:"time"`
	Country  string `json:"country"`
	Type     string `json:"type"`
	// Image    string `json:"image"`
}

type VisibilitySwitch struct {
	Visible bool `json:"visible"`
}

// Функция для создания пустого рецепта
func (server *Server) CreateEmptyRecipeHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Создаем структуру для рецепта
	recipe := models.Recipe{
		User: *user,
	}

	// Сохраняем пустой рецепт в БД
	err = server.DB.Create(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось создать рецепт"})
	}

	// Сохраняем ID рецепта, для передачи на фронтэнд
	recipeID := recipe.ID

	return c.JSON(http.StatusOK, &RecipeResponse{Message: "Создан новый рецепт", Id: recipeID})
}

// Функция для обновления данные о рецепте
func (server *Server) UpdateRecipeHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта с фронтэнда
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Получаем данные о рецепте с фронтэнда
	var recipe_data RecipeData
	err = c.Bind(&recipe_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	// Если не введено название рецепта
	if len(recipe_data.Name) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Название рецепта не может быть пустым"})
	}

	// Проверка на то, что текущий пользователь - автор рецепта
	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	// Заполняем данные о рецепте
	recipe.StrRecipeName = recipe_data.Name
	recipe.IntServings = recipe_data.Servings
	recipe.IntTime = recipe_data.Time
	recipe.StrRecipeCountry = recipe_data.Country
	recipe.StrRecipeType = recipe_data.Type

	// Сохраняем обновленный рецепт в БД
	err = server.DB.Save(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось обновить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Рецепт обновлен"})
}

func (server *Server) GetRecipeHandle(c echo.Context) error {
	// Получаем ID рецепта с фронтэнда
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Если рецепт скрыт из общего доступа, то пишем, что не удалось найти
	if !recipe.BoolRecipeVisibility {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipe)
}

func (server *Server) GetMyRecipeHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта с фронтэнда
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Проверка на то, что текущий пользователь - автор рецепта
	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	return c.JSON(http.StatusOK, recipe)
}

func (server *Server) GetRecipesHandle(c echo.Context) error {
	// Получаем информацию о рецепте
	var recipes []models.Recipe
	err := server.DB.
		Preload("User").
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		Preload("RecipeIngredients.Ingredient").
		Find(&recipes, "bool_recipe_visibility = true").Error
	if err != nil {
		log.Printf("Get all recipes: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipes)
}

func (server *Server) GetMyRecipesHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем информацию о рецепте
	var recipes []models.Recipe
	err = server.DB.
		Preload("User").
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		Preload("RecipeIngredients.Ingredient").
		Find(&recipes, "int_user_id = ?", user.ID).Error
	if err != nil {
		log.Printf("Get all recipes: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipes)
}

func (server *Server) DeleteRecipeHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта с фронтэнда
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Проверка на то, что текущий пользователь - автор рецепта
	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	// Удаляем данные о рецепте из БД
	err = server.DB.Delete(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить удалить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Рецепт удален"})
}

type FindData struct {
	Text string `json:"text"`
}

func (server *Server) FindRecipesHandle(c echo.Context) error {
	var find_data FindData
	err := c.Bind(&find_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	if find_data.Text == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Пустая строка поиска"})
	}

	// Получаем информацию о рецепте
	var recipes []models.Recipe
	err = server.DB.
		Preload("User").
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		Find(&recipes, "LOWER(str_recipe_name) LIKE ?", fmt.Sprintf("%%%s%%", find_data.Text)).Error
	if err != nil {
		log.Printf("Find recipes: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipes)
}

func (server *Server) AddRecipeToFavoritesHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта с фронтэнда
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Добавляем рецепт в избранные пользователя
	user.UserFavorite = append(user.UserFavorite, *recipe)
	err = server.DB.Save(&user).Error
	if err != nil {
		log.Printf("Favorite: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось добавить рецепт в избранное"})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Ок"})
}

// func (server *Server) RemoveRecipeFromFavoritesHandle(c echo.Context) error {
// 	return nil
// }

func (server *Server) ChangeVisibilityRecipeHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта с фронтэнда
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Проверка на то, что текущий пользователь - автор рецепта
	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	// Получаем данные о видимости с фронтэнда
	var visibility VisibilitySwitch
	err = c.Bind(&visibility)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	// Сохраняем значение видимости
	recipe.BoolRecipeVisibility = visibility.Visible

	// Обновляем запись в БД
	err = server.DB.Save(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить обновить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Рецепт обновлен"})
}

// func (server *Server) UploadRecipeCoverHandle(c echo.Context) error {
// 	// Получаем информацию о пользователе
// 	user, err := server.GetUserByClaims(c)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
// 	}

// 	// Получаем ID рецепта с фронтэнда
// 	recipeID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		log.Printf("Recipe id: %s", err.Error())
// 		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
// 	}

// 	// Получаем информацию о рецепте
// 	recipe, err := server.GetRecipeById(recipeID)
// 	if err != nil {
// 		log.Printf("Get recipe by id: %s", err.Error())
// 		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
// 	}

// 	// Проверка на то, что текущий пользователь - автор рецепта
// 	if user.ID != recipe.IntUserId {
// 		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
// 	}

// 	// Получаем файл из формы
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить файл из формы: %s", err.Error())})
// 	}

// 	// Получаем расширение файла
// 	fileExt, err := server.GetFileExtByMimetype(file)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось определить тип файла: %s", err.Error())})
// 	}

// 	// Пытаемся открыть файл
// 	src, err := file.Open()
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось прочитать файл: %s", err.Error())})
// 	}
// 	defer src.Close()

// 	// Сохраняем файл
// 	filename, err := server.SaveFileWithExt(src, fileExt)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось сохранить файл: %s", err.Error())})
// 	}

// 	// Сохраняем обложку рецепта
// 	recipe.StrRecipeImage = filename

// 	// Обновляем данные о рецепте
// 	err = server.DB.Save(&recipe).Error
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось обновить рецепт: %s", err.Error())})
// 	}

// 	return c.JSON(http.StatusOK, &CoverResponse{Message: "Ок", Cover: filename})

// }

type RecipeIngredientInfo struct {
	IngredientId int `json:"ingredient_id"`
	Grams        int `json:"grams"`
}

func (server *Server) AddIngredientHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта, к которому будет добавлен этап
	recipeID, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию о рецепте, к которому будет добавлен этап
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Проверка на то, что текущий пользователь автор рецепта
	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	// Получаем данные с фронтенда
	var recipe_ingredient_info RecipeIngredientInfo
	err = c.Bind(&recipe_ingredient_info)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	var ingredient models.Ingredient
	err = server.DB.First(&ingredient, "id = ?", recipe_ingredient_info.IngredientId).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, &DefaultResponse{Message: "Ингредиент не найден"})
	}

	recipe.RecipeIngredients = append(recipe.RecipeIngredients, models.RecipeIngredient{
		IntGrams:   recipe_ingredient_info.Grams,
		Ingredient: ingredient,
		Recipe:     *recipe,
	})

	err = server.DB.Save(&recipe).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, &DefaultResponse{
				Message: "Не удалось добавить ингредиент для рецепта",
			},
		)
	}

	return c.JSON(http.StatusOK, &DefaultResponse{
		Message: "Ингредиент добавлен",
	})
}

func (server *Server) RemoveIngredientHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Проверка на то, что текущий пользователь автор рецепта
	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}
	// Получаем данные с фронтенда
	var recipe_ingredient_info RecipeIngredientInfo
	err = c.Bind(&recipe_ingredient_info)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	var ingredient models.Ingredient
	err = server.DB.First(&ingredient, "id = ?", recipe_ingredient_info.IngredientId).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, &DefaultResponse{Message: "Ингредиент не найден"})
	}

	var recipeIngredient models.RecipeIngredient
	err = server.DB.Where(&models.RecipeIngredient{
		Recipe:     *recipe,
		Ingredient: ingredient,
	}).First(&recipeIngredient).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, &DefaultResponse{Message: "Ингредиент рецепта не найден"})
	}

	err = server.DB.Unscoped().Delete(&recipeIngredient).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, &DefaultResponse{
				Message: "Не удалось удалить ингредиент для рецепта",
			},
		)
	}

	return c.JSON(http.StatusOK, &DefaultResponse{
		Message: "Ингредиент удален",
	})
}
