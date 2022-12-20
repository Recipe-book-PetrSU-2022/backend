package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Recipe-book-PetrSU-2022/backend/models"
	"github.com/labstack/echo/v4"
)

type StageData struct {
	Description string `json:"description"`
}

// Функция для создания этапа
func (server *Server) CreateStageHandle(c echo.Context) error {
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
	var stage_data StageData
	err = c.Bind(&stage_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("stage data = %+v", stage_data)

	// Создаем этап
	stage := models.Stage{
		StrStageDesc: stage_data.Description,
		Recipe:       *recipe,
	}

	// Сохраняем этап в БД
	err = server.DB.Create(&stage).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, &DefaultResponse{
				Message: "Не удалось создать новый этап рецепта",
			},
		)
	}

	return c.JSON(http.StatusOK, &StageResponse{Message: "Создан новый этап", Id: stage.ID})
}

// Функция для обновления данных об этапе
func (server *Server) UpdateStageHandle(c echo.Context) error {
	//
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта, к которому будет добавлен этап
	recipeID, err := strconv.Atoi(c.Param("stage_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию об этапе, который будем изменять
	stage, err := server.GetStageById(recipeID)
	if err != nil {
		log.Printf("Get stage by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти этап"})
	}

	// Проверка на то, что текущий пользователь автор рецепта
	if user.ID != stage.Recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	// Получаем данные с фронтенда
	var stage_data StageData
	err = c.Bind(&stage_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("stage data = %+v", stage_data)

	// Обновляем описание этапа
	stage.StrStageDesc = stage_data.Description

	// Сохраняем изменееные данные в БД
	err = server.DB.Save(&stage).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, &DefaultResponse{
				Message: "Не удалось обновить этап",
			},
		)
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Этап обновлен"})
}

// Функция для добавления фото к этапу
func (server *Server) AddStagePhotoHandle(c echo.Context) error {
	// Получаем информацию о пользователе
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	// Получаем ID рецепта, к которому будет добавлен этап
	recipeID, err := strconv.Atoi(c.Param("stage_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}

	// Получаем информацию об этапе, который будем изменять
	stage, err := server.GetStageById(recipeID)
	if err != nil {
		log.Printf("Get stage by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти этап"})
	}

	// Проверка на то, что текущий пользователь автор рецепта
	if user.ID != stage.Recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить файл из формы: %s", err.Error())})
	}

	// Получаем расширение файла
	fileExt, err := server.GetFileExtByMimetype(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось определить тип файла: %s", err.Error())})
	}

	// Пытаемся прочитать файл
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось прочитать файл: %s", err.Error())})
	}
	defer src.Close()

	// Сохраняем файл
	filename, err := server.SaveFileWithExt(src, fileExt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: fmt.Sprintf("Не удалось сохранить файл: %s", err.Error())})
	}

	// Добавляем имя файла к этапу
	stage.StagePhotos = append(stage.StagePhotos, models.Photo{
		StrImage: filename,
	})

	// Обновляем данные об этапе в БД
	err = server.DB.Save(&stage).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, &DefaultResponse{
				Message: "Не удалось обновить этап",
			},
		)
	}

	return c.JSON(http.StatusOK, &DefaultResponse{Message: "Этап обновлен"})
}
