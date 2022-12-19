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

func (server *Server) CreateEmptyRecipeHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipe := models.Recipe{
		User: *user,
	}

	err = server.DB.Create(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось создать рецепт"})
	}

	recipeID := recipe.ID

	log.Printf("New recipe id = %d", recipeID)

	return c.JSON(http.StatusOK, &RecipeResponse{Message: "Создан новый рецепт", Id: recipeID})
}

func (server *Server) UpdateRecipeHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	var recipe_data RecipeData

	err = c.Bind(&recipe_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("%+v", recipe_data)

	if len(recipe_data.Name) == 0 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Название рецепта не может быть пустым"})
	}

	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	recipe.StrRecipeName = recipe_data.Name
	recipe.IntServings = recipe_data.Servings
	recipe.IntTime = recipe_data.Time
	recipe.StrRecipeCountry = recipe_data.Country
	recipe.StrRecipeType = recipe_data.Type

	err = server.DB.Save(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить обновить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Рецепт обновлен"})
}

func (server *Server) GetRecipeHandle(c echo.Context) error {
	recipeIDStr := c.Param("id")
	recipeID, err := strconv.Atoi(recipeIDStr)

	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	if !recipe.BoolRecipeVisibility {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipe)
}

func (server *Server) GetMyRecipeHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeIDStr := c.Param("id")
	recipeID, err := strconv.Atoi(recipeIDStr)

	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipe)
}

func (server *Server) GetRecipesHandle(c echo.Context) error {
	var recipes []models.Recipe

	err := server.DB.
		Preload("User").
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		Find(&recipes, "bool_recipe_visibility = true").Error
	if err != nil {
		log.Printf("Get all recipes: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipes)
}

func (server *Server) GetMyRecipesHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	var recipes []models.Recipe

	err = server.DB.
		Preload("User").
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		Find(&recipes, "int_user_id = ?", user.ID).Error
	if err != nil {
		log.Printf("Get all recipes: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, recipes)
}

func (server *Server) DeleteRecipeHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	err = server.DB.Delete(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить удалить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Рецепт удален"})
}

// хз что тут должно быть
func (server *Server) FindRecipesHandle(c echo.Context) error {
	return nil
}

func (server *Server) AddRecipeToFavoritesHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	user.UserFavorite = append(user.UserFavorite, *recipe)
	err = server.DB.Save(&user).Error
	if err != nil {
		log.Printf("Favorite: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось добавить рецепт в избранное"})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Ок"})
}

func (server *Server) RemoveRecipeFromFavoritesHandle(c echo.Context) error {
	return nil
}

func (server *Server) ChangeVisibilityRecipeHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	var visibility VisibilitySwitch

	err = c.Bind(&visibility)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("visibility = %+v", visibility)

	recipe.BoolRecipeVisibility = visibility.Visible

	err = server.DB.Save(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить обновить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Рецепт обновлен"})
}

func (server *Server) UploadRecipeCoverHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить файл из формы: %s", err.Error())})
	}

	fileExt, err := server.GetFileExtByMimetype(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось определить тип файла: %s", err.Error())})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось прочитать файл: %s", err.Error())})
	}
	defer src.Close()

	filename, err := server.SaveFileWithExt(src, fileExt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить сохранить файл: %s", err.Error())})
	}

	recipe.StrRecipeImage = filename
	err = server.DB.Save(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить обновить рецепт: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, &CoverResponse{Message: "Ок", Cover: filename})

}
