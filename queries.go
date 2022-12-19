package main

import "github.com/Recipe-book-PetrSU-2022/backend/models"

func (server *Server) GetRecipeById(id int) (*models.Recipe, error) {
	var recipe models.Recipe

	err := server.DB.
		Preload("RecipeStages").
		Preload("RecipeStages.StagePhotos").
		Preload("RecipeComments").
		Preload("RecipeIngredients").
		First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}
