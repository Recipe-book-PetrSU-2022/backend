package main

import (
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

func (server *Server) CreateEmptyRecipeHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipe := models.Recipe{
		User: user,
	}

	err = server.DB.Create(&recipe).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось создать рецепт"})
	}

	recipeID := recipe.ID

	log.Printf("New recipe id = %d", recipeID)

	return c.JSON(http.StatusOK, &RecipeResponse{Message: "Создан новый рецепт", Id: recipeID})
}

func (server *Server) CreateRecipeHandle(c echo.Context) error {
	// user, err := server.GetUserByClaims(c)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	// }

	// recipeID := c.QueryParam("id")

	// var recipe_data RecipeData

	// err = c.Bind(&recipe_data)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	// }

	// log.Printf("%+v", recipe_data)

	// if len(recipe_data.Name) == 0 {
	// 	return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Название рецепта не может быть пустым"})
	// }

	// server.DB.Model(&user).Update("StrUserName", recipe_data.Name)

	return nil
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

	// stages := StageToResponse()

	// return c.JSON(http.StatusOK, &RecipeResponse{
	// 	Message:   "Ок",
	// 	Id:        recipe.ID,
	// 	UserId:    recipe.User.ID,
	// 	Name:      recipe.StrRecipeName,
	// 	Servings:  recipe.IntServings,
	// 	Time:      recipe.IntTime,
	// 	Country:   recipe.StrRecipeCountry,
	// 	Type:      recipe.StrRecipeType,
	// 	Cover:     recipe.StrRecipeImage,
	// 	IsVisible: recipe.BoolRecipeVisibility,
	// 	Stages:    stages,
	// })

	return c.JSON(http.StatusOK, recipe)
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

func (server *Server) ChangeVisibilityRecipeHandle(c echo.Context) error {
	return nil
}
