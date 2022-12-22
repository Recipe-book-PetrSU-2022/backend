package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
)

func (server *Server) GetIngredients(c echo.Context) error {
	var ingredients []models.Ingredient
	err := server.DB.Find(&ingredients).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{
			Message: "Не удалось получить список ингредиентов",
		})
	}

	return c.JSON(http.StatusOK, ingredients)
}

type IngredientData struct {
	Name          string `json:"name"`
	Colories      int    `json:"colories"`
	Proteins      int    `json:"proteins"`
	Fats          int    `json:"fats"`
	Carbohydrates int    `json:"carbohydrates"`
}

func (server *Server) NewIngredient(c echo.Context) error {
	var ingredient_data IngredientData
	err := c.Bind(&ingredient_data)
	if err != nil {
		log.Print(err)

		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	ingredient := models.Ingredient{
		StrIngredientName: ingredient_data.Name,
		IntCalories:       ingredient_data.Colories,
		IntProteins:       ingredient_data.Proteins,
		IntFats:           ingredient_data.Fats,
		IntCarbohydrates:  ingredient_data.Carbohydrates,
	}

	err = server.DB.Create(&ingredient).Error
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось создать ингредиент"})
	}

	return c.JSON(http.StatusOK, &IngredientResponse{
		Message: "Ингредиент создан!",
		Id:      ingredient.ID,
	})
}
