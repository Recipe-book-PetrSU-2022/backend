package main

import "github.com/Recipe-book-PetrSU-2022/backend/models"

// Функция для получения полной информации о рецепте
func (server *Server) GetRecipeById(id int) (*models.Recipe, error) {
	var recipe models.Recipe

	err := server.DB.
		Preload("User").
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		Preload("RecipeIngredients.Ingredient").
		First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

// Функция для получения информации о этапе
func (server *Server) GetStageById(id int) (*models.Stage, error) {
	var stage models.Stage

	err := server.DB.
		Preload("Recipe").
		Preload("StagePhotos").
		First(&stage, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &stage, nil
}
