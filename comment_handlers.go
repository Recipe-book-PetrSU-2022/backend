package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
)

type CommentData struct {
	Text string `json:"text"`
	Rate int    `json:"rate"`
}

func (server *Server) GetCommentHandle(c echo.Context) error {

	recipeID, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		log.Printf("Comment id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id комментария"})
	}

	var comment models.Comment
	err = server.DB.First(&comment, "int_recipe_id = ? AND id = ?", recipeID, commentID).Error
	if err != nil {
		log.Printf("Get comment by id: %s", err.Error())
		return c.JSON(http.StatusNotFound, &DefaultResponse{
			Message: "Комментарий не найден",
		})
	}
	return c.JSON(http.StatusOK, &comment)
}

func (server *Server) GetCommentsHandle(c echo.Context) error {

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
	if !recipe.BoolRecipeVisibility {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	return c.JSON(http.StatusOK, &recipe.RecipeComments)
}

func (server *Server) CreateCommentHandle(c echo.Context) error {
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

	// Получаем информацию о рецепте, к которому будет добавлен этап
	recipe, err := server.GetRecipeById(recipeID)
	if err != nil {
		log.Printf("Get recipe by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}
	if !recipe.BoolRecipeVisibility {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти рецепт"})
	}

	// Проверка на то, что текущий пользователь НЕ автор рецепта
	if user.ID == recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Нельзя оставить отзыв о собственном рецепте"})
	}

	var comment_data CommentData
	err = c.Bind(&comment_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	if comment_data.Rate < 0 || comment_data.Rate > 5 {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Оценка пользователя должна быть от 0 до 5"})
	}

	comment := models.Comment{
		StrCommentDesc:      comment_data.Text,
		IntRate:             comment_data.Rate,
		Recipe:              *recipe,
		User:                *user,
		IntCommentTimestamp: int(time.Now().Unix()),
	}

	err = server.DB.Create(&comment).Error
	if err != nil {
		log.Printf("Create comment error: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалость создать комментарий"})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Комментарий создан!"})
}

func (server *Server) DeleteCommentHandle(c echo.Context) error {
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

	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id комментария"})
	}

	var comment models.Comment
	err = server.DB.First(&comment, "int_recipe_id = ? AND int_user_id = ? AND id = ?", recipeID, user.ID, commentID).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, &DefaultResponse{
			Message: "Комментарий не найден",
		})
	}

	err = server.DB.Delete(&comment).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{
			Message: "Не удалось удалить комментарий",
		})
	}

	return c.JSON(http.StatusOK, &DefaultResponse{
		Message: "Комментарий удален",
	})
}
