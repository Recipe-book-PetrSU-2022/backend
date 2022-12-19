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

func (server *Server) CreateStageHandle(c echo.Context) error {
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

	if user.ID != recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	var stage_data StageData
	err = c.Bind(&stage_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("stage data = %+v", stage_data)

	stage := models.Stage{
		StrStageDesc: stage_data.Description,
		Recipe:       *recipe,
	}

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

func (server *Server) UpdateStageHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("stage_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}
	stage, err := server.GetStageById(recipeID)
	if err != nil {
		log.Printf("Get stage by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти этап"})
	}

	if user.ID != stage.Recipe.IntUserId {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Рецепт принадлежит другому пользователю"})
	}

	var stage_data StageData
	err = c.Bind(&stage_data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: fmt.Sprintf("Не удалось получить данные от пользователя: %s", err.Error())})
	}

	log.Printf("stage data = %+v", stage_data)

	stage.StrStageDesc = stage_data.Description

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

func (server *Server) AddStagePhotoHandle(c echo.Context) error {
	user, err := server.GetUserByClaims(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Не удалось найти пользователя"})
	}

	recipeID, err := strconv.Atoi(c.Param("stage_id"))
	if err != nil {
		log.Printf("Recipe id: %s", err.Error())
		return c.JSON(http.StatusBadRequest, &DefaultResponse{Message: "Неверный id рецепта"})
	}
	stage, err := server.GetStageById(recipeID)
	if err != nil {
		log.Printf("Get stage by id: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &DefaultResponse{Message: "Не удалось найти этап"})
	}

	if user.ID != stage.Recipe.IntUserId {
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

	stage.StagePhotos = append(stage.StagePhotos, models.Photo{
		StrImage: filename,
	})

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
