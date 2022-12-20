package main

import (
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
